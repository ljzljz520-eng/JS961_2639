package service

import (
	"path/filepath"
	"testing"
	"time"
	"weddingtemplates/internal/repository"
	"weddingtemplates/internal/storage"
)

func TestRegisterReview(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "s.db"))
	defer s.Close()
	svc := New(repository.New(s), FixedClock{T: time.Unix(1, 0)})
	if _, e := svc.Register("1", "Template", "Hall", "East", 80); e != nil {
		t.Fatal(e)
	}
	if _, e := svc.Review("1", "reviewer"); e != nil {
		t.Fatal(e)
	}
}
