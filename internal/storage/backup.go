package storage

import (
	"encoding/json"
	"os"
	"time"
	"weddingtemplates/internal/domain"
)

func (s *Store) WriteSnapshot(path string, now time.Time) error {
	snap, e := s.Snapshot(now)
	if e != nil {
		return e
	}
	data, e := json.MarshalIndent(snap, "", "  ")
	if e != nil {
		return e
	}
	return os.WriteFile(path, data, 0600)
}
func ReadSnapshot(path string) (Snapshot, error) {
	data, e := os.ReadFile(path)
	if e != nil {
		return Snapshot{}, e
	}
	var snap Snapshot
	e = json.Unmarshal(data, &snap)
	return snap, e
}
func (s *Store) RestoreSnapshot(path string) error {
	snap, e := ReadSnapshot(path)
	if e != nil {
		return e
	}
	return s.Restore(snap)
}
func (s *Store) ValidateEntity(id string) (bool, error) {
	r, e := s.Record(id)
	if e != nil {
		return false, e
	}
	return r.Validate() == nil, nil
}

var _ domain.Record
