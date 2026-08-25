package storage

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"time"
	"weddingtemplates/internal/domain"
)

type Snapshot struct {
	Records   []domain.Record
	Audits    []domain.AuditEvent
	CreatedAt time.Time
}

func (s *Store) Snapshot(now time.Time) (Snapshot, error) {
	rs, e := s.ListRecords()
	if e != nil {
		return Snapshot{}, e
	}
	var as []domain.AuditEvent
	s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(buckets["audits"]).ForEach(func(_, v []byte) error {
			var a domain.AuditEvent
			if json.Unmarshal(v, &a) == nil {
				as = append(as, a)
			}
			return nil
		})
	})
	return Snapshot{Records: rs, Audits: as, CreatedAt: now}, nil
}
func (s *Store) Restore(snapshot Snapshot) error {
	for _, r := range snapshot.Records {
		if e := s.PutRecord(r); e != nil {
			return e
		}
	}
	for _, a := range snapshot.Audits {
		if e := s.PutAudit(a); e != nil {
			return e
		}
	}
	return nil
}
