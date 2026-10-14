package client

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"homecloud/internal/database"
)

var (
	crypt32DLL             = windows.NewLazySystemDLL("crypt32.dll")
	procCryptProtectData   = crypt32DLL.NewProc("CryptProtectData")
	procCryptUnprotectData = crypt32DLL.NewProc("CryptUnprotectData")
	kernel32DLL            = windows.NewLazySystemDLL("kernel32.dll")
	procLocalFree          = kernel32DLL.NewProc("LocalFree")
)

// cryptProtectUIForbidden suppresses any interactive DPAPI prompt.
const cryptProtectUIForbidden = 0x01

const devicesKey = "saved_devices.v1"

// dataBlob mirrors the Win32 DATA_BLOB layout used by Crypt(Un)ProtectData.
type dataBlob struct {
	cbData uint32
	pbData *byte
}

// Device is a remembered mini-PC profile used by the SSH and RDP tabs.
// Password is only populated by GetDevice; ListDevices never returns secrets.
type Device struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	ServerPort  int    `json:"serverPort"`
	User        string `json:"user"`
	Fingerprint string `json:"fingerprint"`
	RDPPort     int    `json:"rdpPort"`
	HasPassword bool   `json:"hasPassword"`
	Password    string `json:"password,omitempty"`
	// ClearPassword erases the stored secret on the next SaveDevice call.
	ClearPassword bool `json:"clearPassword"`
}

// deviceRecord is the on-disk form: Device plus the DPAPI-encrypted secret.
type deviceRecord struct {
	Device
	Secret string `json:"secret"`
	Used   int64  `json:"used"`
}

var (
	devStoreMu sync.Mutex
	devStore   *database.Store
)

func clientStorePath() (string, error) {
	if override := strings.TrimSpace(os.Getenv("HOMECLOUD_CLIENT_DB")); override != "" {
		return override, nil
	}
	base := strings.TrimSpace(os.Getenv("LOCALAPPDATA"))
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, "AppData", "Local")
	}
	dir := filepath.Join(base, "HomeCloud")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "client.db"), nil
}

func deviceStore() (*database.Store, error) {
	devStoreMu.Lock()
	defer devStoreMu.Unlock()
	if devStore != nil {
		return devStore, nil
	}
	p, err := clientStorePath()
	if err != nil {
		return nil, err
	}
	s, err := database.Open(p)
	if err != nil {
		return nil, fmt.Errorf("cannot open device store: %w", err)
	}
	devStore = s
	return s, nil
}

func utf16Bytes(s string) ([]byte, error) {
	u, err := windows.UTF16FromString(s)
	if err != nil {
		return nil, err
	}
	if n := len(u); n > 0 && u[n-1] == 0 {
		u = u[:n-1]
	}
	b := make([]byte, len(u)*2)
	for i, c := range u {
		b[2*i] = byte(c)
		b[2*i+1] = byte(c >> 8)
	}
	return b, nil
}

func dpapiCall(proc *windows.LazyProc, in []byte) ([]byte, error) {
	if len(in) == 0 {
		return nil, nil
	}
	inBlob := dataBlob{cbData: uint32(len(in)), pbData: &in[0]}
	var outBlob dataBlob
	desc, err := windows.UTF16PtrFromString("HomeCloud saved device credential")
	if err != nil {
		return nil, err
	}
	r, _, callErr := proc.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		uintptr(unsafe.Pointer(desc)),
		0, 0, 0,
		cryptProtectUIForbidden,
		uintptr(unsafe.Pointer(&outBlob)),
	)
	if r == 0 {
		return nil, fmt.Errorf("DPAPI call failed: %w", callErr)
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outBlob.pbData)))
	out := make([]byte, outBlob.cbData)
	copy(out, unsafe.Slice(outBlob.pbData, outBlob.cbData))
	return out, nil
}

func protectSecret(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	raw, err := utf16Bytes(plain)
	if err != nil {
		return "", err
	}
	enc, err := dpapiCall(procCryptProtectData, raw)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}

func unprotectSecret(encoded string) (string, error) {
	if strings.TrimSpace(encoded) == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	dec, err := dpapiCall(procCryptUnprotectData, raw)
	if err != nil {
		return "", err
	}
	if len(dec) == 0 {
		return "", nil
	}
	u := make([]uint16, len(dec)/2)
	for i := range u {
		u[i] = uint16(dec[2*i]) | uint16(dec[2*i+1])<<8
	}
	return windows.UTF16ToString(u), nil
}

func loadDevices() ([]deviceRecord, error) {
	s, err := deviceStore()
	if err != nil {
		return nil, err
	}
	raw, err := s.Get(devicesKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var recs []deviceRecord
	if err = json.Unmarshal([]byte(raw), &recs); err != nil {
		return nil, fmt.Errorf("saved device list is corrupted: %w", err)
	}
	return recs, nil
}

func persistDevices(recs []deviceRecord) error {
	s, err := deviceStore()
	if err != nil {
		return err
	}
	b, err := json.Marshal(recs)
	if err != nil {
		return err
	}
	return s.Set(devicesKey, string(b))
}

func newDeviceID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}

// ListDevices returns saved devices, most recently used first, without passwords.
func (a *App) ListDevices() ([]Device, error) {
	recs, err := loadDevices()
	if err != nil {
		return nil, err
	}
	out := make([]Device, 0, len(recs))
	sort.SliceStable(recs, func(i, j int) bool { return recs[i].Used > recs[j].Used })
	for _, r := range recs {
		d := r.Device
		d.Password = ""
		d.HasPassword = r.Secret != ""
		out = append(out, d)
	}
	return out, nil
}

// GetDevice returns one saved device including its decrypted password.
func (a *App) GetDevice(id string) (Device, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Device{}, fmt.Errorf("не указано устройство")
	}
	recs, err := loadDevices()
	if err != nil {
		return Device{}, err
	}
	for _, r := range recs {
		if r.ID != id {
			continue
		}
		d := r.Device
		d.HasPassword = r.Secret != ""
		if d.HasPassword {
			pw, err := unprotectSecret(r.Secret)
			if err != nil {
				return Device{}, fmt.Errorf("не удалось расшифровать пароль (другой пользователь Windows?): %w", err)
			}
			d.Password = pw
		}
		return d, nil
	}
	return Device{}, fmt.Errorf("устройство не найдено")
}

// SaveDevice stores a device profile. An empty Password keeps the previously
// saved secret; pass ClearPassword to erase it.
func (a *App) SaveDevice(d Device) (Device, error) {
	host := strings.TrimSpace(d.Host)
	if host == "" || strings.ContainsAny(host, " /\\\r\n\t") {
		return Device{}, fmt.Errorf("укажите IP/имя хоста без пробелов и схемы")
	}
	if len(host) > 253 {
		return Device{}, fmt.Errorf("слишком длинное имя хоста")
	}
	user := strings.TrimSpace(d.User)
	if user == "" || strings.ContainsAny(user, "\x00\r\n") {
		return Device{}, fmt.Errorf("укажите пользователя SSH")
	}
	name := strings.TrimSpace(d.Name)
	if name == "" {
		name = host
	}
	if len(name) > 120 || strings.ContainsAny(name, "\x00\r\n") {
		return Device{}, fmt.Errorf("недопустимое имя устройства")
	}
	if d.Port == 0 {
		d.Port = 22
	}
	if d.Port < 1 || d.Port > 65535 {
		return Device{}, fmt.Errorf("SSH-порт должен быть 1–65535")
	}
	if d.ServerPort == 0 {
		d.ServerPort = 8081
	}
	if d.ServerPort < 1 || d.ServerPort > 65535 {
		return Device{}, fmt.Errorf("порт панели должен быть 1–65535")
	}
	if d.RDPPort == 0 {
		d.RDPPort = 3389
	}
	if d.RDPPort < 1 || d.RDPPort > 65535 {
		return Device{}, fmt.Errorf("RDP-порт должен быть 1–65535")
	}
	if d.Fingerprint != "" && !strings.HasPrefix(d.Fingerprint, "SHA256:") {
		return Device{}, fmt.Errorf("недопустимый отпечаток хоста")
	}

	recs, err := loadDevices()
	if err != nil {
		return Device{}, err
	}
	id := strings.TrimSpace(d.ID)
	idx := -1
	for i := range recs {
		if (id != "" && recs[i].ID == id) ||
			(id == "" && recs[i].Host == host && recs[i].Port == d.Port && recs[i].User == user) {
			idx = i
			break
		}
	}
	rec := deviceRecord{}
	if idx >= 0 {
		rec = recs[idx]
	} else if id != "" {
		for i := range recs {
			if recs[i].ID == id {
				rec = recs[i]
				idx = i
				break
			}
		}
	}
	if rec.ID == "" {
		rec.ID = newDeviceID()
	}
	secret := rec.Secret
	if d.ClearPassword {
		secret = ""
	} else if d.Password != "" {
		secret, err = protectSecret(d.Password)
		if err != nil {
			return Device{}, fmt.Errorf("не удалось защитить пароль: %w", err)
		}
	}
	rec.Secret = secret
	rec.Name, rec.Host, rec.User = name, host, user
	rec.Port, rec.ServerPort, rec.RDPPort = d.Port, d.ServerPort, d.RDPPort
	rec.Fingerprint = d.Fingerprint
	rec.Password = ""
	rec.Used = time.Now().UnixNano()

	if idx >= 0 {
		recs[idx] = rec
	} else {
		recs = append(recs, rec)
	}
	if err = persistDevices(recs); err != nil {
		return Device{}, err
	}
	out := rec.Device
	out.HasPassword = rec.Secret != ""
	return out, nil
}

// DeleteDevice removes a saved device profile and its encrypted secret.
func (a *App) DeleteDevice(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("не указано устройство")
	}
	recs, err := loadDevices()
	if err != nil {
		return err
	}
	out := recs[:0]
	found := false
	for _, r := range recs {
		if r.ID == id {
			found = true
			continue
		}
		out = append(out, r)
	}
	if !found {
		return fmt.Errorf("устройство не найдено")
	}
	return persistDevices(out)
}