package api

import (
	"net/http"
	"strconv"
	"weddingtemplates/internal/domain"
)

func parseQuery(r *http.Request) domain.Query {
	q := domain.Query{Text: r.URL.Query().Get("q"), Region: r.URL.Query().Get("region"), Venue: r.URL.Query().Get("venue"), Status: r.URL.Query().Get("status")}
	if v, e := strconv.Atoi(r.URL.Query().Get("min_score")); e == nil {
		q.MinScore = v
	}
	if v, e := strconv.Atoi(r.URL.Query().Get("limit")); e == nil {
		q.Limit = v
	}
	if v, e := strconv.Atoi(r.URL.Query().Get("offset")); e == nil {
		q.Offset = v
	}
	return q
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
func methodAllowed(r *http.Request, methods ...string) bool {
	for _, m := range methods {
		if r.Method == m {
			return true
		}
	}
	return false
}
func statusForError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return http.StatusBadRequest
}
func normalizeActor(r *http.Request) string {
	actor := r.Header.Get("X-Actor")
	if actor == "" {
		actor = "api"
	}
	return actor
}
func decodeID(r *http.Request) string {
	if r.URL.Query().Get("id") != "" {
		return r.URL.Query().Get("id")
	}
	return ""
}
