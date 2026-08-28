package flow080

import (
	"path/filepath"
	"testing"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/repository"
	"weddingtemplates/internal/service"
	"weddingtemplates/internal/storage"
)

func Test961BusinessRegression(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "bug.db"))
	defer s.Close()
	svc := service.New(repository.New(s), DeterministicClock())
	h := NewBatchHandler(svc, DeterministicClock().T)
	r, e := h.Sync([]domain.ImportRow{{ExternalID: "first", Title: "First", Score: 61}, {ExternalID: "second", Title: "Second", Score: 94}})
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Accepted) != 2 || r.Accepted[1].Score != 94 {
		t.Fatalf("expected second score 94, got %#v", r.Accepted)
	}
}
