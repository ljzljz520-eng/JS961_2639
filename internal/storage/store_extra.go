package storage

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"weddingtemplates/internal/domain"
)

func (s *Store) Attachments(recordID string) (out []domain.Attachment, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(buckets["attachments"]).ForEach(func(_, v []byte) error {
			var a domain.Attachment
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
func (s *Store) DeleteRecord(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(buckets["records"]).Delete([]byte(id)) })
}
func (s *Store) Count(bucket string) (n int, err error) {
	err = s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(buckets[bucket]).ForEach(func(_, _ []byte) error { n++; return nil })
	})
	return
}
