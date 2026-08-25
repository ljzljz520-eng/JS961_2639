package service

import (
	"path/filepath"
	"testing"
	"time"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/repository"
	"weddingtemplates/internal/storage"
)

func TestSummary(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "sum.db"))
	defer s.Close()
	svc := New(repository.New(s), FixedClock{T: time.Unix(1, 0)})
	svc.Import([]domain.ImportRow{{ExternalID: "1", Title: "A", Score: 80}, {ExternalID: "2", Title: "B", Score: 90}}, "x")
	sum, e := svc.Summary()
	if e != nil || sum.Total != 2 || sum.AverageScore != 85 {
		t.Fatalf("%v %#v", e, sum)
	}
}
