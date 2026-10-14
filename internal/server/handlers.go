package server

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"sort"
)

func (s *Server) routes(m *http.ServeMux) {
	m.HandleFunc("/api/v1/health", s.health)
	m.HandleFunc("/api/v1/stats", s.stats)
	m.HandleFunc("/api/v1/processes", s.processes)
	m.HandleFunc("/api/v1/containers", s.containers)
	m.HandleFunc("/api/v1/containers/control", s.control)
}
func (s *Server) auth(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == "/api/v1/health" {
		return true
	}
	if r.Header.Get("Authorization") != "Bearer "+s.token {
		http.Error(w, "unauthorized", 401)
		return false
	}
	return true
}
func write(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]string{"status": "ok"})
}
func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	if !s.auth(w, r) {
		return
	}
	m, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(0, false)
	write(w, map[string]any{"cpu_percent": c, "memory": m})
}

type ProcessInfo struct {
	PID    int32  `json:"pid"`
	Name   string `json:"name"`
	Memory uint64 `json:"memory_bytes"`
}

func (s *Server) processes(w http.ResponseWriter, r *http.Request) {
	if !s.auth(w, r) {
		return
	}
	ps, _ := process.Processes()
	out := make([]ProcessInfo, 0, len(ps))
	for _, p := range ps {
		mi, e := p.MemoryInfo()
		n, _ := p.Name()
		if e == nil {
			out = append(out, ProcessInfo{p.Pid, n, mi.RSS})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Memory > out[j].Memory })
	if len(out) > 20 {
		out = out[:20]
	}
	write(w, out)
}
func (s *Server) containers(w http.ResponseWriter, r *http.Request) {
	if !s.auth(w, r) {
		return
	}
	if s.docker == nil {
		http.Error(w, "docker unavailable", 503)
		return
	}
	v, e := s.docker.List(r.Context())
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	write(w, v)
}
func (s *Server) control(w http.ResponseWriter, r *http.Request) {
	if !s.auth(w, r) {
		return
	}
	var x struct {
		ContainerID string `json:"container_id"`
		Action      string `json:"action"`
	}
	if json.NewDecoder(r.Body).Decode(&x) != nil || s.docker == nil {
		http.Error(w, "bad request", 400)
		return
	}
	if e := s.docker.Control(r.Context(), x.ContainerID, x.Action); e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	write(w, map[string]string{"status": "ok"})
}
