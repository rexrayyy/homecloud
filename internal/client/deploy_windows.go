package client

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

//go:embed payload/homecloud-linux-amd64
var linuxServer []byte

// SSHConfig is supplied by the local desktop form. Passwords are never persisted.
type SSHConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	ServerPort  int    `json:"serverPort"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

func sshAddress(c SSHConfig) (string, error) {
	if strings.TrimSpace(c.Host) == "" || strings.ContainsAny(c.Host, " /\\\r\n\t") || c.Port < 1 || c.Port > 65535 {
		return "", fmt.Errorf("укажите IP/имя хоста и SSH-порт 1–65535")
	}
	return net.JoinHostPort(c.Host, fmt.Sprint(c.Port)), nil
}

// ProbeHostKey reads the untrusted fingerprint without sending a password.
// The user must verify it via a separate trusted channel before deployment.
func (a *App) ProbeHostKey(c SSHConfig) (string, error) {
	var fingerprint string
	conn, err := connectSSH(c, func(_ string, _ net.Addr, key ssh.PublicKey) error {
		fingerprint = ssh.FingerprintSHA256(key)
		return fmt.Errorf("host key probe complete")
	}, false)
	if conn != nil {
		conn.Close()
	}
	if fingerprint != "" {
		return fingerprint, nil
	}
	return "", err
}

func connectSSH(c SSHConfig, callback ssh.HostKeyCallback, authenticate bool) (*ssh.Client, error) {
	addr, err := sshAddress(c)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{User: c.User, HostKeyCallback: callback, Timeout: 15 * time.Second}
	if authenticate {
		config.Auth = []ssh.AuthMethod{ssh.Password(c.Password)}
	}
	raw, err := net.DialTimeout("tcp", addr, 15*time.Second)
	if err != nil {
		return nil, err
	}
	raw.SetDeadline(time.Now().Add(30 * time.Second))
	cc, chans, reqs, err := ssh.NewClientConn(raw, addr, config)
	if err != nil {
		raw.Close()
		return nil, err
	}
	raw.SetDeadline(time.Time{})
	return ssh.NewClient(cc, chans, reqs), nil
}

// Deploy uploads the embedded Linux binary and starts an unprivileged background
// server. All commands are fixed strings; credentials never enter shell commands.
func (a *App) Deploy(c SSHConfig) (string, error) {
	script, err := deploymentScript(c.ServerPort)
	if err != nil {
		return "", err
	}
	if c.User == "" || c.Password == "" || !strings.HasPrefix(c.Fingerprint, "SHA256:") {
		return "", fmt.Errorf("нужны логин, пароль и подтверждённый отпечаток")
	}
	conn, err := connectSSH(c, func(_ string, _ net.Addr, key ssh.PublicKey) error {
		if ssh.FingerprintSHA256(key) != c.Fingerprint {
			return fmt.Errorf("SSH-ключ изменился: подключение заблокировано")
		}
		return nil
	}, true)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	deadline := time.AfterFunc(3*time.Minute, func() { _ = conn.Close() })
	defer deadline.Stop()
	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	session.Stdin = bytes.NewReader(linuxServer)
	output, err := session.CombinedOutput(script)
	session.Close()
	if err != nil {
		return "", fmt.Errorf("установка: %w\n%s", err, output)
	}
	transport := &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		return conn.DialContext(ctx, network, addr)
	}}
	defer transport.CloseIdleConnections()
	hc := &http.Client{Transport: transport, Timeout: 3 * time.Second}
	for i := 0; i < 15; i++ {
		resp, e := hc.Get(fmt.Sprintf("http://127.0.0.1:%d/api/v1/health", c.ServerPort))
		if e == nil {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
			if resp.StatusCode == 200 && strings.Contains(string(body), `"status":"ok"`) {
				return fmt.Sprintf("Сервер запущен и проверен через SSH. Адрес на мини-ПК: 127.0.0.1:%d. Файлы: ~/.local/share/homecloud/instances/%d; журнал: server.log. После перезагрузки нажмите запуск снова.", c.ServerPort, c.ServerPort), nil
			}
		}
		time.Sleep(time.Second)
	}
	return "", fmt.Errorf("сервер на порту %d не ответил; возможно, порт занят. Выберите другой порт или проверьте ~/.local/share/homecloud/instances/%d/server.log на мини-ПК", c.ServerPort, c.ServerPort)
}

// Only a validated integer is interpolated into the remote shell script.
func deploymentScript(port int) (string, error) {
	if port < 1024 || port > 65535 {
		return "", fmt.Errorf("порт сервера должен быть от 1024 до 65535")
	}
	return strings.ReplaceAll(deployScript, "__SERVER_PORT__", fmt.Sprint(port)), nil
}

const deployScript = `set -eu
test "$(uname -s)" = Linux || { echo 'Requires Linux'; exit 1; }
test "$(uname -m)" = x86_64 || { echo 'Requires Intel/AMD x86_64'; exit 1; }
command -v flock >/dev/null
command -v nohup >/dev/null
umask 077
mkdir -p "$HOME/.local/share/homecloud/instances/__SERVER_PORT__"
cd "$HOME/.local/share/homecloud/instances/__SERVER_PORT__"
exec 9>deploy.lock
flock -n 9 || { echo 'Another deployment is running'; exit 1; }
if ! flock -n server.lock true; then
 echo 'Server already running; checking health'
 exit 0
fi
cat > homecloud.upload
chmod 700 homecloud.upload
mv homecloud.upload homecloud
mkdir -p files state
HOMECLOUD_STORAGE="$PWD/files" HOMECLOUD_STATE="$PWD/state" nohup flock -n server.lock ./homecloud --mode=server --addr=127.0.0.1:__SERVER_PORT__ >server.log 2>&1 < /dev/null 9>&- &
echo 'Server launched'
`
