package storage

import (
	"fmt"
	"go.etcd.io/bbolt"
	"weddingtemplates/internal/domain"
)

func (s *Store) Healthy() bool { return s != nil && s.db != nil }
func (s *Store) VerifyBuckets() error {
	if !s.Healthy() {
		return fmt.Errorf("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		for name, b := range buckets {
			if tx.Bucket(b) == nil {
				return fmt.Errorf("missing bucket %s", name)
			}
		}
		return nil
	})
}
func (s *Store) EntityCounts() (map[string]int, error) {
	out := map[string]int{}
	for name := range buckets {
		n, e := s.Count(name)
		if e != nil {
			return nil, e
		}
		out[name] = n
	}
	return out, nil
}
func (s *Store) HasRecord(id string) bool { _, e := s.Record(id); return e == nil }
func (s *Store) ValidateAll() ([]string, error) {
	items, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	bad := []string{}
	for _, r := range items {
		if r.Validate() != nil {
			bad = append(bad, r.ID)
		}
	}
	return bad, nil
}

var _ = domain.StatusDraft
