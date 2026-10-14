package client

import (
	"fmt"
	"strings"
	"testing"
)

func TestSSHAddress(t *testing.T) {
	for _, c := range []SSHConfig{{Host: "", Port: 22}, {Host: "bad host", Port: 22}, {Host: "host", Port: 0}, {Host: "host", Port: 65536}} {
		if _, e := sshAddress(c); e == nil {
			t.Fatalf("accepted invalid address: %+v", c)
		}
	}
	got, err := sshAddress(SSHConfig{Host: "::1", Port: 22})
	if err != nil || got != "[::1]:22" {
		t.Fatalf("%q %v", got, err)
	}
}
func TestDeployRequiresConfirmation(t *testing.T) {
	if _, err := (&App{}).Deploy(SSHConfig{Host: "localhost", Port: 22, ServerPort: 8081, User: "test", Password: "test"}); err == nil {
		t.Fatal("accepted unconfirmed host key")
	}
}

func TestDeploymentPort(t *testing.T) {
	for _, port := range []int{-1, 0, 22, 1023, 65536} {
		if _, err := deploymentScript(port); err == nil {
			t.Fatalf("accepted %d", port)
		}
	}
	for _, port := range []int{1024, 8081, 9000, 65535} {
		script, err := deploymentScript(port)
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{fmt.Sprintf("--addr=127.0.0.1:%d", port), fmt.Sprintf("/instances/%d", port)} {
			if !strings.Contains(script, want) {
				t.Fatalf("missing %s", want)
			}
		}
		if strings.Contains(script, "__SERVER_PORT__") {
			t.Fatal("unexpanded port")
		}
	}
}
