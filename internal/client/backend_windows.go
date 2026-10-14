package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/pkg/sftp"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/windows"
	"homecloud/internal/database"
)

type ConnectionStatus struct {
	Connected    bool   `json:"connected"`
	Host         string `json:"host"`
	ServerPort   int    `json:"serverPort"`
	TokenReady   bool   `json:"tokenReady"`
	TerminalOpen bool   `json:"terminalOpen"`
	Message      string `json:"message"`
}
type FileEntry struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"isDir"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}
type App struct {
	mu        sync.Mutex
	connectMu sync.Mutex
	ctx       context.Context
	conn      *ssh.Client
	http      *http.Client
	token     string
	host      string
	port      int
	terminal  *ssh.Session
	input     io.WriteCloser
	sftp      *sftp.Client
	rdpHost   string
}

func (a *App) startup(ctx context.Context) {
	a.mu.Lock()
	a.ctx = ctx
	a.mu.Unlock()
	wruntime.WindowCenter(ctx)
}
func (a *App) shutdown(context.Context) { _ = a.Disconnect() }
func (a *App) emit(event, text string) {
	a.mu.Lock()
	ctx := a.ctx
	a.mu.Unlock()
	if ctx != nil {
		wruntime.EventsEmit(ctx, event, text)
	}
}
func (a *App) Status() ConnectionStatus {
	a.mu.Lock()
	defer a.mu.Unlock()
	return ConnectionStatus{Connected: a.conn != nil, Host: a.host, ServerPort: a.port, TokenReady: a.token != "", TerminalOpen: a.terminal != nil}
}
func (a *App) Connect(c SSHConfig) (ConnectionStatus, error) {
	a.connectMu.Lock()
	defer a.connectMu.Unlock()
	if _, err := deploymentScript(c.ServerPort); err != nil {
		return a.Status(), err
	}
	if strings.TrimSpace(c.User) == "" || c.Password == "" || !strings.HasPrefix(c.Fingerprint, "SHA256:") {
		return a.Status(), fmt.Errorf("SSH user, password and independently verified fingerprint required")
	}
	conn, err := connectSSH(c, func(_ string, _ net.Addr, key ssh.PublicKey) error {
		if ssh.FingerprintSHA256(key) != c.Fingerprint {
			return fmt.Errorf("SSH host fingerprint mismatch; connection blocked")
		}
		return nil
	}, true)
	if err != nil {
		return a.Status(), err
	}
	timer := time.AfterFunc(30*time.Second, func() { conn.Close() })
	token, err := retrieveToken(conn, c.ServerPort)
	if err != nil {
		timer.Stop()
		conn.Close()
		return a.Status(), err
	}
	tr := &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		return conn.DialContext(ctx, network, addr)
	}}
	hc := &http.Client{Transport: tr, Timeout: 30 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:%d/api/v1/stats", c.ServerPort), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := hc.Do(req)
	if err == nil {
		resp.Body.Close()
		if resp.StatusCode != 200 {
			err = fmt.Errorf("server authentication check: HTTP %d", resp.StatusCode)
		}
	}
	timer.Stop()
	if err != nil {
		tr.CloseIdleConnections()
		conn.Close()
		return a.Status(), err
	}
	a.disconnect()
	a.mu.Lock()
	a.conn = conn
	a.http = hc
	a.token = token
	a.host = c.Host
	a.port = c.ServerPort
	a.mu.Unlock()
	go func() {
		_ = conn.Wait()
		a.mu.Lock()
		current := a.conn == conn
		a.mu.Unlock()
		if current {
			a.disconnectIf(conn)
			a.emit("connection:closed", "SSH connection closed")
		}
	}()
	return a.Status(), nil
}
func (a *App) disconnectIf(expected *ssh.Client) {
	a.mu.Lock()
	if a.conn != expected {
		a.mu.Unlock()
		return
	}
	conn := a.conn
	hc := a.http
	t := a.terminal
	sf := a.sftp
	a.conn = nil
	a.http = nil
	a.token = ""
	a.terminal = nil
	a.input = nil
	a.sftp = nil
	a.mu.Unlock()
	if t != nil {
		t.Close()
	}
	if sf != nil {
		sf.Close()
	}
	if hc != nil {
		hc.CloseIdleConnections()
	}
	if conn != nil {
		conn.Close()
	}
}
func (a *App) disconnect() { a.mu.Lock(); conn := a.conn; a.mu.Unlock(); a.disconnectIf(conn) }
func (a *App) Disconnect() error {
	a.connectMu.Lock()
	defer a.connectMu.Unlock()
	a.disconnect()
	return nil
}

// Restrict the temporary directory's Windows ACL BEFORE writing any secret data.
func privateTemp() (string, error) {
	dir, err := os.MkdirTemp("", "homecloud-token-")
	if err != nil {
		return "", err
	}
	fail := func(e error) (string, error) { os.RemoveAll(dir); return "", e }
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return fail(err)
	}
	sd, err := windows.SecurityDescriptorFromString("D:P(A;OICI;FA;;;" + user.User.Sid.String() + ")")
	if err != nil {
		return fail(err)
	}
	acl, _, err := sd.DACL()
	if err != nil {
		return fail(err)
	}
	if err = windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil); err != nil {
		return fail(err)
	}
	return dir, nil
}
func retrieveToken(conn *ssh.Client, port int) (string, error) {
	dir, err := privateTemp()
	if err != nil {
		return "", fmt.Errorf("secure token temporary directory: %w", err)
	}
	defer os.RemoveAll(dir)
	name := filepath.Join(dir, "homecloud.db")
	f, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	session, err := conn.NewSession()
	if err != nil {
		f.Close()
		return "", err
	}
	defer session.Close()
	out, err := session.StdoutPipe()
	if err != nil {
		f.Close()
		return "", err
	}
	err = session.Start(fmt.Sprintf("cat -- \"$HOME/.local/share/homecloud/instances/%d/state/homecloud.db\"", port))
	if err != nil {
		f.Close()
		return "", fmt.Errorf("cannot read server token database")
	}
	n, copyErr := io.Copy(f, io.LimitReader(out, 16*1024*1024+1))
	closeErr := f.Close()
	if n > 16*1024*1024 {
		session.Close()
		return "", fmt.Errorf("token database exceeds size limit")
	}
	waitErr := session.Wait()
	if copyErr != nil || closeErr != nil || waitErr != nil {
		return "", fmt.Errorf("cannot copy server token database; deploy this instance first")
	}
	db, err := database.Open(name)
	if err != nil {
		return "", fmt.Errorf("cannot open token database snapshot")
	}
	defer db.Close()
	token, err := db.Get("auth_token")
	decoded, e := hex.DecodeString(token)
	if err != nil || e != nil || len(decoded) != 32 {
		return "", fmt.Errorf("server token missing or invalid; deploy/restart instance")
	}
	return token, nil
}
func (a *App) remote(ctx context.Context, method, p string, body io.Reader) (*http.Response, error) {
	a.mu.Lock()
	hc, token, port := a.http, a.token, a.port
	a.mu.Unlock()
	if hc == nil {
		return nil, fmt.Errorf("connect SSH first")
	}
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://127.0.0.1:%d", port)+p, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if method == "PROPFIND" {
		req.Header.Set("Depth", "1")
		req.Header.Set("Content-Type", "application/xml")
	} else if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SSH server request failed: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("server returned HTTP %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return resp, nil
}
func (a *App) Request(p, body string) (string, error) {
	method := "GET"
	switch p {
	case "/api/v1/health", "/api/v1/stats", "/api/v1/processes", "/api/v1/containers":
		if body != "" {
			return "", fmt.Errorf("read endpoint does not accept a body")
		}
	case "/api/v1/containers/control":
		var v struct {
			ContainerID string `json:"container_id"`
			Action      string `json:"action"`
		}
		if len(body) > 4096 || json.Unmarshal([]byte(body), &v) != nil || v.ContainerID == "" {
			return "", fmt.Errorf("invalid container request")
		}
		switch v.Action {
		case "start", "stop", "restart":
		default:
			return "", fmt.Errorf("unsupported container action")
		}
		method = "POST"
	default:
		return "", fmt.Errorf("API endpoint is not supported")
	}
	resp, err := a.remote(context.Background(), method, p, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024*1024+1))
	if err != nil {
		return "", err
	}
	if len(b) > 8*1024*1024 || !json.Valid(b) {
		return "", fmt.Errorf("invalid or oversized server JSON response")
	}
	return string(b), nil
}
func (a *App) MountDrive(letter string) (string, error) {
	return "", fmt.Errorf("automatic drive mapping unavailable: Windows WebClient cannot use HomeCloud Bearer authentication. Use Files over SSH; no unauthenticated proxy is opened")
}

type eventWriter struct{ app *App }

func (w eventWriter) Write(p []byte) (int, error) {
	w.app.emit("terminal:output", string(p))
	return len(p), nil
}
func (a *App) OpenTerminal() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.conn == nil {
		return fmt.Errorf("connect SSH first")
	}
	if a.terminal != nil {
		return nil
	}
	s, err := a.conn.NewSession()
	if err != nil {
		return err
	}
	if err = s.RequestPty("xterm-256color", 24, 80, ssh.TerminalModes{ssh.ECHO: 1, ssh.TTY_OP_ISPEED: 14400, ssh.TTY_OP_OSPEED: 14400}); err != nil {
		s.Close()
		return err
	}
	input, err := s.StdinPipe()
	if err != nil {
		s.Close()
		return err
	}
	s.Stdout = eventWriter{a}
	s.Stderr = eventWriter{a}
	if err = s.Shell(); err != nil {
		s.Close()
		return err
	}
	a.terminal = s
	a.input = input
	go func() {
		err := s.Wait()
		s.Close()
		a.mu.Lock()
		if a.terminal == s {
			a.terminal = nil
			a.input = nil
		}
		a.mu.Unlock()
		msg := "Terminal closed"
		if err != nil {
			msg = "Terminal session ended"
		}
		a.emit("terminal:closed", msg)
	}()
	return nil
}
func (a *App) TerminalInput(data string) error {
	if len(data) > 65536 {
		return fmt.Errorf("terminal input too large")
	}
	a.mu.Lock()
	in := a.input
	a.mu.Unlock()
	if in == nil {
		return fmt.Errorf("open terminal first")
	}
	_, err := io.WriteString(in, data)
	return err
}
func (a *App) ResizeTerminal(cols, rows int) error {
	if cols < 2 || cols > 500 || rows < 2 || rows > 200 {
		return fmt.Errorf("invalid terminal size")
	}
	a.mu.Lock()
	s := a.terminal
	a.mu.Unlock()
	if s == nil {
		return fmt.Errorf("open terminal first")
	}
	return s.WindowChange(rows, cols)
}
func (a *App) CloseTerminal() error {
	a.mu.Lock()
	s := a.terminal
	a.terminal = nil
	a.input = nil
	a.mu.Unlock()
	if s != nil {
		return s.Close()
	}
	return nil
}

type DiscoveredServer struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (a *App) Discover() ([]DiscoveredServer, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	entries := make(chan *zeroconf.ServiceEntry, 32)
	if err = resolver.Browse(ctx, "_homecloud._tcp", "local.", entries); err != nil {
		return nil, err
	}
	out := []DiscoveredServer{}
	seen := map[string]bool{}
	for {
		select {
		case <-ctx.Done():
			return out, nil
		case e, ok := <-entries:
			if !ok {
				return out, nil
			}
			host := strings.TrimSuffix(e.HostName, ".")
			if len(e.AddrIPv4) > 0 {
				host = e.AddrIPv4[0].String()
			}
			key := fmt.Sprintf("%s:%d", host, e.Port)
			if !seen[key] {
				seen[key] = true
				out = append(out, DiscoveredServer{e.Instance, host, e.Port})
			}
		}
	}
}

// CaptureScreen takes a screenshot of the remote desktop via SSH.
// It tries multiple methods: scrot, import (ImageMagick), or xwd + convert.
// Returns base64-encoded PNG image data.
func (a *App) CaptureScreen() (string, error) {
	a.mu.Lock()
	conn := a.conn
	a.mu.Unlock()
	if conn == nil {
		return "", fmt.Errorf("connect SSH first")
	}

	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// Build a script that sets DISPLAY/XAUTHORITY and tries multiple screenshot tools.
	// This is needed because SSH sessions don't inherit the desktop session environment.
	script := `
export DISPLAY="${DISPLAY:-:0}"
# Try common XAUTHORITY locations for desktop sessions
if [ -z "$XAUTHORITY" ]; then
  for auth in "$HOME/.Xauthority" "/run/user/$(id -u)/gdm/Xauthority" "/run/user/$(id -u)/.Xauthority" "/var/run/lightdm/$(whoami)/:0"; do
    if [ -f "$auth" ]; then
      export XAUTHORITY="$auth"
      break
    fi
  done
fi

tmp="/tmp/homecloud_screen_$$.png"

# Try scrot first (most common on Linux desktops)
if command -v scrot >/dev/null 2>&1; then
  if scrot -o "$tmp" 2>/dev/null; then
    base64 "$tmp"
    rm -f "$tmp"
    exit 0
  fi
fi

# Try ImageMagick import
if command -v import >/dev/null 2>&1; then
  if import -window root "$tmp" 2>/dev/null; then
    base64 "$tmp"
    rm -f "$tmp"
    exit 0
  fi
fi

# Try xwd + convert
if command -v xwd >/dev/null 2>&1 && command -v convert >/dev/null 2>&1; then
  if xwd -root -silent 2>/dev/null | convert xwd:- png:"$tmp" 2>/dev/null; then
    base64 "$tmp"
    rm -f "$tmp"
    exit 0
  fi
fi

rm -f "$tmp"
echo "no screenshot tool available or display not accessible" >&2
exit 1
`

	output, err := session.Output(script)
	if err != nil {
		return "", fmt.Errorf("screen capture failed: %v (ensure an active desktop session with DISPLAY=:0 and scrot, imagemagick, or xwd+convert installed)", err)
	}

	// Strip newlines from base64 output (some systems wrap lines)
	b64 := strings.NewReplacer("\n", "", "\r", "", " ", "").Replace(string(output))
	if len(b64) < 100 {
		return "", fmt.Errorf("screen capture returned no image data (is a desktop session active?)")
	}

	return b64, nil
}

// GetScreenInfo returns basic screen information from the remote system
func (a *App) GetScreenInfo() (map[string]interface{}, error) {
	a.mu.Lock()
	conn := a.conn
	a.mu.Unlock()
	if conn == nil {
		return nil, fmt.Errorf("connect SSH first")
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// Try to get screen resolution with proper DISPLAY setup
	script := `
export DISPLAY="${DISPLAY:-:0}"
if [ -z "$XAUTHORITY" ]; then
  for auth in "$HOME/.Xauthority" "/run/user/$(id -u)/gdm/Xauthority" "/run/user/$(id -u)/.Xauthority"; do
    if [ -f "$auth" ]; then
      export XAUTHORITY="$auth"
      break
    fi
  done
fi
if command -v xdpyinfo >/dev/null 2>&1; then
  xdpyinfo 2>/dev/null | grep -E "dimensions|depth of root" || echo "display info unavailable"
else
  echo "xdpyinfo not available"
fi
`

	output, err := session.Output(script)
	if err != nil {
		return map[string]interface{}{"status": "unavailable", "detail": err.Error()}, nil
	}

	return map[string]interface{}{
		"status":     "ok",
		"dimensions": strings.TrimSpace(string(output)),
	}, nil
}
