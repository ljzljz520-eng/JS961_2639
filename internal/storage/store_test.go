package storage

import (
	"path/filepath"
	"testing"
	"time"
	"weddingtemplates/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "data.db")
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	r := domain.NewRecord("persist", "Saved", "Hall", "East", 77, time.Unix(10, 0))
	if e = s.PutRecord(r); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.Record("persist")
	if e != nil || got.Score != 77 {
		t.Fatalf("%v %#v", e, got)
	}
}
