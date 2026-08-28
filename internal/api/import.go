package api

import (
	"net/http"
	"weddingtemplates/internal/domain"
)

func (s *Server) importTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	in, e := parseImport(r)
	if e != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": e.Error()})
		return
	}
	result, e := s.Service.Import(in.Rows, in.Actor)
	if e != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": e.Error()})
		return
	}
	writeJSON(w, http.StatusAccepted, result)
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	sum, e := s.Service.Summary()
	if e != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": e.Error()})
		return
	}
	writeJSON(w, http.StatusOK, sum)
}
func queryFrom(r *http.Request) domain.Query {
	return domain.Query{Text: r.URL.Query().Get("q"), Region: r.URL.Query().Get("region"), Venue: r.URL.Query().Get("venue"), Status: r.URL.Query().Get("status"), Limit: 50}
}
