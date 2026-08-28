package api

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"weddingtemplates/internal/repository"
	"weddingtemplates/internal/service"
	"weddingtemplates/internal/storage"
)

func TestHealthEndpoint(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "a.db"))
	defer s.Close()
	h := New(service.New(repository.New(s), service.FixedClock{})).Handler()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/health", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "ok") {
		t.Fatal(w.Code, w.Body.String())
	}
}
