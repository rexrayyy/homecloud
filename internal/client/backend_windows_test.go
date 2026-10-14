package client

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBackendRejectsUnsupportedRequests(t *testing.T) {
	a := &App{}
	for _, tt := range []struct{ path, body string }{{"https://example.com", ""}, {"/api/v1/stats?x=1", ""}, {"/api/v1/stats", "{}"}, {"/api/v1/containers/control", `{"container_id":"x","action":"delete"}`}} {
		if _, err := a.Request(tt.path, tt.body); err == nil {
			t.Fatalf("accepted %s", tt.path)
		}
	}
	if _, err := a.Request("/api/v1/stats", ""); err == nil || !strings.Contains(err.Error(), "connect SSH") {
		t.Fatalf("disconnected: %v", err)
	}
	if _, err := a.Connect(SSHConfig{ServerPort: 8081, User: "x", Password: "x"}); err == nil {
		t.Fatal("accepted missing pin")
	}
}
func TestBackendPrivateTemp(t *testing.T) {
	dir, err := privateTemp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	if err = os.WriteFile(filepath.Join(dir, "secret"), []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}
}
