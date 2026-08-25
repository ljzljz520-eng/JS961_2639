package api

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func respondError(w http.ResponseWriter, status int, code, message, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorBody{Code: code, Message: message, RequestID: requestID})
}
func requestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return "local-request"
}
func acceptJSON(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type") == "application/json"
}
func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Actor")
}
