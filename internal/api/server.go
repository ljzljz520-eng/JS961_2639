package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/service"
)

type Server struct{ Service *service.Service }

func New(s *service.Service) *Server { return &Server{Service: s} }
func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/health", s.health)
	m.HandleFunc("/templates", s.templates)
	m.HandleFunc("/templates/", s.template)
	m.HandleFunc("/import", s.importTemplates)
	m.HandleFunc("/metrics", s.metrics)
	return m
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
func (s *Server) templates(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		q := domain.Query{Text: r.URL.Query().Get("q"), Region: r.URL.Query().Get("region"), Status: r.URL.Query().Get("status"), Limit: 50}
		p, e := s.Service.Search(q)
		if e != nil {
			http.Error(w, e.Error(), 500)
			return
		}
		_ = json.NewEncoder(w).Encode(p)
		return
	}
	if r.Method == http.MethodPost {
		var in domain.Record
		if e := json.NewDecoder(r.Body).Decode(&in); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		out, e := s.Service.Register(in.ID, in.Title, in.Venue, in.Region, in.Score)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(out)
		return
	}
	http.Error(w, "method not allowed", 405)
}
func (s *Server) template(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/templates/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		r0, e := s.Service.Repo.FindRecord(id)
		if e != nil {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(r0)
	case http.MethodPost:
		action := r.URL.Query().Get("action")
		var e error
		switch action {
		case "review":
			_, e = s.Service.Review(id, "api")
		case "approve":
			_, e = s.Service.Approve(id, "api")
		case "archive":
			_, e = s.Service.Archive(id, "api")
		default:
			http.Error(w, "unknown action", 400)
			return
		}
		if e != nil {
			http.Error(w, e.Error(), 409)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", 405)
	}
}
