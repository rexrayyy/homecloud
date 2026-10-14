package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/grandcat/zeroconf"
	"golang.org/x/net/webdav"
	"homecloud/internal/database"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Server struct {
	token   string
	storage string
	docker  *DockerController
	db      *database.Store
}

func Run(addr string) error {
	storage := os.Getenv("HOMECLOUD_STORAGE")
	if storage == "" {
		storage = "./storage"
	}
	if err := os.MkdirAll(storage, 0750); err != nil {
		return err
	}
	state := os.Getenv("HOMECLOUD_STATE")
	if state == "" {
		state = "./state"
	}
	if err := os.MkdirAll(state, 0700); err != nil {
		return err
	}
	db, err := database.Open(filepath.Join(state, "homecloud.db"))
	if err != nil {
		return err
	}
	defer db.Close()
	token, _ := db.Get("auth_token")
	if token == "" {
		b := make([]byte, 32)
		if _, err = rand.Read(b); err != nil {
			return err
		}
		token = hex.EncodeToString(b)
		if err = db.Set("auth_token", token); err != nil {
			return err
		}
	}
	d, _ := NewDockerController()
	s := &Server{token: token, storage: storage, docker: d, db: db}
	mux := http.NewServeMux()
	s.routes(mux)
	dav := http.StripPrefix("/dav", &webdav.Handler{FileSystem: webdav.Dir(storage), LockSystem: webdav.NewMemLS()})
	mux.HandleFunc("/dav/", func(w http.ResponseWriter, r *http.Request) {
		if s.auth(w, r) {
			dav.ServeHTTP(w, r)
		}
	})
	host, portText, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		return err
	}
	if host != "127.0.0.1" && host != "localhost" {
		beacon, e := zeroconf.Register("homecloud", "_homecloud._tcp", "local.", port, []string{"path=/api/v1/health"}, nil)
		if e == nil {
			defer beacon.Shutdown()
		} else {
			log.Printf("mDNS registration failed: %v", e)
		}
	}
	log.Printf("HomeCloud server listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

var _ = context.Background
