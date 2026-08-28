package repository

import (
	"path/filepath"
	"testing"
	"time"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/storage"
)

func TestSearchRepository(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "r.db"))
	defer s.Close()
	r := New(s)
	r.SaveRecord(domain.NewRecord("1", "One", "Hall", "East", 90, time.Unix(1, 0)))
	r.SaveRecord(domain.NewRecord("2", "Two", "Hall", "West", 70, time.Unix(2, 0)))
	p, e := r.Search(domain.Query{Region: "East"})
	if e != nil || p.Total != 1 {
		t.Fatal(e, p)
	}
}
