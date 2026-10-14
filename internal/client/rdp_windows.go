package client

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// RDPRequest describes one Remote Desktop connection attempt.
type RDPRequest struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Fullscreen bool   `json:"fullscreen"`
	Remember   bool   `json:"remember"`
}

// RDPStatus is the result of a TCP reachability probe of the RDP port.
type RDPStatus struct {
	Reachable bool   `json:"reachable"`
	Target    string `json:"target"`
	Message   string `json:"message"`
}

const defaultRDPPort = 3389

func validateRDPHost(host string) (string, error) {
	host = strings.TrimSpace(host)
	if host == "" || strings.ContainsAny(host, " /\\\r\n\t") {
		return "", fmt.Errorf("укажите IP/имя хоста RDP без пробелов и схемы")
	}
	if len(host) > 253 {
		return "", fmt.Errorf("слишком длинное имя хоста")
	}
	return host, nil
}

func validateRDPPort(port int) (int, error) {
	if port == 0 {
		port = defaultRDPPort
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("RDP-порт должен быть 1–65535")
	}
	return port, nil
}

func rdpTarget(host string, port int) string {
	if port == defaultRDPPort {
		return host
	}
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func system32(name string) string {
	root := strings.TrimSpace(os.Getenv("SystemRoot"))
	if root == "" {
		root = `C:\Windows`
	}
	return filepath.Join(root, "System32", name)
}

// runHidden runs a console helper without flashing a window.
func runHidden(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// RDPCheck reports whether the RDP port answers from this machine.
func (a *App) RDPCheck(r RDPRequest) (RDPStatus, error) {
	host, err := validateRDPHost(r.Host)
	if err != nil {
		return RDPStatus{}, err
	}
	port, err := validateRDPPort(r.Port)
	if err != nil {
		return RDPStatus{}, err
	}
	target := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", target, 6*time.Second)
	if err != nil {
		return RDPStatus{
			Reachable: false,
			Target:    target,
			Message: fmt.Sprintf("RDP-порт %s не отвечает: %v. Проверьте, что на мини-ПК "+
				"установлен и запущен xrdp, и что порт открыт.", target, err),
		}, nil
	}
	conn.Close()
	return RDPStatus{Reachable: true, Target: target, Message: "RDP-порт " + target + " отвечает"}, nil
}

func rdpFileContents(host string, port int, user string, width, height int, fullscreen bool) string {
	mode := 1
	if fullscreen {
		mode = 2
	}
	lines := []string{
		"screen mode id:i:" + strconv.Itoa(mode),
		"use multimon:i:0",
		"desktopwidth:i:" + strconv.Itoa(width),
		"desktopheight:i:" + strconv.Itoa(height),
		"session bpp:i:32",
		"compression:i:1",
		"keyboardhook:i:2",
		"audiocapturemode:i:0",
		"videoplaybackmode:i:1",
		"displayconnectionbar:i:1",
		"enableworkspacereconnect:i:0",
		"disable full window drag:i:1",
		"disable menu anims:i:1",
		"disable themes:i:0",
		"bitmapcachepersistenable:i:1",
		"full address:s:" + rdpTarget(host, port),
		"audiomode:i:0",
		"redirectprinters:i:0",
		"redirectcomports:i:0",
		"redirectsmartcards:i:0",
		"redirectclipboard:i:1",
		"redirectposdevices:i:0",
		"redirectdrives:i:0",
		"autoreconnection enabled:i:1",
		"authentication level:i:2",
		"prompt for credentials on client:i:0",
		"negotiate security layer:i:1",
		"remoteapplicationmode:i:0",
		"alternate shell:s:",
		"shell working directory:s:",
		"gatewayhostname:s:",
		"gatewayusagemethod:i:4",
		"gatewaycredentialssource:i:4",
		"gatewayprofileusagemethod:i:0",
		"promptcredentialonce:i:0",
		"gatewaybrokeringtype:i:0",
		"use redirection server name:i:0",
		"rdgiskdcproxy:i:0",
		"kdcproxyname:s:",
	}
	if user != "" {
		lines = append(lines, "username:s:"+user)
	}
	return strings.Join(lines, "\r\n") + "\r\n"
}

// writeRDPFile stores the generated .rdp in the per-user private temp folder so
// it inherits the same restricted ACL as the SSH key material.
func writeRDPFile(contents string) (string, error) {
	dir, err := privateTemp()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, "homecloud-rdp.rdp")
	if err = os.WriteFile(p, []byte(contents), 0o600); err != nil {
		return "", fmt.Errorf("не удалось записать .rdp: %w", err)
	}
	return p, nil
}

// seedRDPCredential stores the RDP login in Windows Credential Manager so mstsc
// signs in automatically on every later launch.
func seedRDPCredential(host, user, password string) error {
	out, err := runHidden(system32("cmdkey.exe"),
		"/generic:TERMSRV/"+host, "/user:"+user, "/pass:"+password)
	if err != nil {
		if out != "" {
			return fmt.Errorf("cmdkey: %v (%s)", err, out)
		}
		return fmt.Errorf("cmdkey: %w", err)
	}
	return nil
}

// RDPConnect launches a real Remote Desktop session with mstsc.exe.
func (a *App) RDPConnect(r RDPRequest) (string, error) {
	host, err := validateRDPHost(r.Host)
	if err != nil {
		return "", err
	}
	port, err := validateRDPPort(r.Port)
	if err != nil {
		return "", err
	}
	user := strings.TrimSpace(r.User)
	if strings.ContainsAny(user, "\x00\r\n") {
		return "", fmt.Errorf("недопустимое имя пользователя")
	}
	width, height := r.Width, r.Height
	if width < 640 || width > 7680 {
		width = 1920
	}
	if height < 480 || height > 4320 {
		height = 1080
	}

	a.mu.Lock()
	a.rdpHost = host
	a.mu.Unlock()

	seeded := false
	if r.Remember && user != "" && r.Password != "" {
		if err = seedRDPCredential(host, user, r.Password); err != nil {
			return "", fmt.Errorf("не удалось сохранить учётные данные RDP: %w", err)
		}
		seeded = true
	}

	file, err := writeRDPFile(rdpFileContents(host, port, user, width, height, r.Fullscreen))
	if err != nil {
		return "", err
	}
	mstsc := system32("mstsc.exe")
	if _, statErr := os.Stat(mstsc); statErr != nil {
		mstsc = "mstsc.exe"
	}
	cmd := exec.Command(mstsc, file)
	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("не удалось запустить mstsc.exe: %w", err)
	}
	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}
	msg := "RDP-сессия запущена: " + rdpTarget(host, port)
	if seeded {
		msg += " (учётные данные сохранены в Windows)"
	}
	return msg, nil
}

// RDPForget removes the stored Windows credential for the last RDP host.
func (a *App) RDPForget() (string, error) {
	a.mu.Lock()
	host := a.rdpHost
	a.rdpHost = ""
	a.mu.Unlock()
	if host == "" {
		return "", fmt.Errorf("нет сохранённой RDP-цели")
	}
	if _, err := validateRDPHost(host); err != nil {
		return "", err
	}
	if _, err := runHidden(system32("cmdkey.exe"), "/delete:TERMSRV/"+host); err != nil {
		return "", fmt.Errorf("не удалось удалить учётные данные: %w", err)
	}
	return "Учётные данные RDP для " + host + " удалены", nil
}

// rdpProbeScript builds a fixed shell probe; only a validated port is embedded.
func rdpProbeScript(port int) string {
	var b strings.Builder
	b.WriteString("set -u\n")
	b.WriteString("p=" + strconv.Itoa(port) + "\n")
	b.WriteString(`if command -v xrdp >/dev/null 2>&1; then echo "xrdp=installed"; else echo "xrdp=missing"; fi` + "\n")
	b.WriteString(`if command -v xrdp-sesman >/dev/null 2>&1; then echo "sesman=installed"; else echo "sesman=missing"; fi` + "\n")
	b.WriteString("l=\"\"\n")
	b.WriteString(`if command -v ss >/dev/null 2>&1; then l=$(ss -ltn 2>/dev/null || true)` + "\n")
	b.WriteString(`elif command -v netstat >/dev/null 2>&1; then l=$(netstat -ltn 2>/dev/null || true); fi` + "\n")
	b.WriteString(`if printf '%s\n' "$l" | grep -q ":$p "; then echo "listening=yes"; else echo "listening=no"; fi` + "\n")
	b.WriteString(`if [ "$(id -u)" = "0" ]; then echo "sudo=root"; else echo "sudo=user"; fi` + "\n")
	return b.String()
}

// RDPProbe inspects the mini-PC over the open SSH session and reports whether
// an RDP server is installed and listening.
func (a *App) RDPProbe(port int) (map[string]interface{}, error) {
	port, err := validateRDPPort(port)
	if err != nil {
		return nil, err
	}
	a.mu.Lock()
	conn := a.conn
	a.mu.Unlock()
	if conn == nil {
		return map[string]interface{}{
			"status":  "no-ssh",
			"message": "Подключитесь по SSH, чтобы проверить наличие xrdp на мини-ПК.",
		}, nil
	}
	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	output, err := session.Output(rdpProbeScript(port))
	if err != nil {
		return map[string]interface{}{
			"status":  "probe-failed",
			"message": "Не удалось выполнить проверку по SSH: " + err.Error(),
		}, nil
	}
	values := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if k, v, ok := strings.Cut(strings.TrimSpace(line), "="); ok {
			values[k] = v
		}
	}
	installed := values["xrdp"] == "installed"
	listening := values["listening"] == "yes"
	portText := strconv.Itoa(port)
	message := "xrdp установлен и слушает порт " + portText
	if !installed {
		message = "xrdp не установлен. Выполните на мини-ПК: sudo apt install xrdp && sudo systemctl enable --now xrdp"
	} else if !listening {
		message = "xrdp установлен, но порт " + portText + " не слушается. Проверьте: sudo systemctl status xrdp"
	}
	status := "unavailable"
	if installed && listening {
		status = "ok"
	}
	return map[string]interface{}{
		"status":        status,
		"xrdpInstalled": installed,
		"sesman":        values["sesman"],
		"listening":     listening,
		"port":          port,
		"message":       message,
	}, nil
}