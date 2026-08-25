package flow080

import (
	"path/filepath"
	"testing"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/repository"
	"weddingtemplates/internal/service"
	"weddingtemplates/internal/storage"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "f.db"))
	defer s.Close()
	svc := service.New(repository.New(s), DeterministicClock())
	svc.Register("1", "A", "H", "E", 80)
	r, e := RunCreateReviewArchive(svc, "1")
	if e != nil || r.Status != domain.StatusArchived {
		t.Fatal(e, r)
	}
}
func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "f2.db"))
	defer s.Close()
	svc := service.New(repository.New(s), DeterministicClock())
	svc.Register("1", "A", "H", "E", 80)
	p, r, e := RunSearchUpdatePublish(svc, domain.Query{Text: "A"}, "1")
	if e != nil || p.Total != 1 || r.Status != domain.StatusApproved {
		t.Fatal(e, p, r)
	}
}
func TestWorkflowImportReport(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "f3.db"))
	defer s.Close()
	svc := service.New(repository.New(s), DeterministicClock())
	r, sum, e := RunImportReport(svc, []domain.ImportRow{{ExternalID: "1", Title: "A", Score: 80}})
	if e != nil || len(r.Accepted) != 1 || sum.Total != 1 {
		t.Fatal(e, r, sum)
	}
}
