package storage

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"path/filepath"
	"sync"
	"weddingtemplates/internal/domain"
)

var buckets = map[string][]byte{"records": []byte("records"), "audits": []byte("audits"), "workflows": []byte("workflows"), "attachments": []byte("attachments")}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, e := tx.CreateBucketIfNotExists(b); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) Path() string      { return s.path }
func encode(v any) ([]byte, error) { return json.Marshal(v) }
func put(tx *bbolt.Tx, b []byte, key string, v any) error {
	raw, e := encode(v)
	if e != nil {
		return e
	}
	return tx.Bucket(b).Put([]byte(key), raw)
}
func get[T any](tx *bbolt.Tx, b []byte, key string) (T, error) {
	var out T
	raw := tx.Bucket(b).Get([]byte(key))
	if raw == nil {
		return out, fmt.Errorf("not found: %s", key)
	}
	e := json.Unmarshal(raw, &out)
	return out, e
}
func (s *Store) PutRecord(r domain.Record) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, buckets["records"], r.Key(), r) })
}
func (s *Store) Record(id string) (domain.Record, error) {
	var r domain.Record
	e := s.db.View(func(tx *bbolt.Tx) error { var x error; r, x = get[domain.Record](tx, buckets["records"], id); return x })
	return r, e
}
func (s *Store) ListRecords() (out []domain.Record, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(buckets["records"]).ForEach(func(_, v []byte) error {
			var r domain.Record
			if e := json.Unmarshal(v, &r); e != nil {
				return e
			}
			out = append(out, r)
			return nil
		})
	})
	return
}
func (s *Store) PutAudit(a domain.AuditEvent) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, buckets["audits"], a.Key(), a) })
}
func (s *Store) Audits(recordID string) (out []domain.AuditEvent, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(buckets["audits"]).ForEach(func(_, v []byte) error {
			var a domain.AuditEvent
			if e := json.Unmarshal(v, &a); e != nil {
				return e
			}
			if a.RecordID == recordID {
				out = append(out, a)
			}
			return nil
		})
	})
	return
}
func (s *Store) PutWorkflow(w domain.Workflow) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, buckets["workflows"], w.Key(), w) })
}
func (s *Store) Workflow(id string) (domain.Workflow, error) {
	var w domain.Workflow
	e := s.db.View(func(tx *bbolt.Tx) error {
		var x error
		w, x = get[domain.Workflow](tx, buckets["workflows"], id)
		return x
	})
	return w, e
}
func (s *Store) PutAttachment(a domain.Attachment) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, buckets["attachments"], a.Key(), a) })
}
