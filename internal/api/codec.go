package api

import (
	"encoding/json"
	"net/http"
	"weddingtemplates/internal/domain"
)

type importRequest struct {
	Rows  []domain.ImportRow `json:"rows"`
	Actor string             `json:"actor"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func parseImport(r *http.Request) (importRequest, error) {
	var in importRequest
	e := json.NewDecoder(r.Body).Decode(&in)
	if in.Actor == "" {
		in.Actor = "api"
	}
	return in, e
}
