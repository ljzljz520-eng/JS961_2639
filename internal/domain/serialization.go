package domain

import (
	"encoding/json"
	"fmt"
)

func EncodeRecord(r Record) ([]byte, error) {
	if e := r.Validate(); e != nil {
		return nil, e
	}
	return json.Marshal(r)
}
func DecodeRecord(data []byte) (Record, error) {
	var r Record
	if e := json.Unmarshal(data, &r); e != nil {
		return r, e
	}
	if e := r.Validate(); e != nil {
		return r, e
	}
	return r, nil
}
func EncodeAudit(a AuditEvent) ([]byte, error) {
	if a.ID == "" || a.RecordID == "" {
		return nil, fmt.Errorf("audit identity required")
	}
	return json.Marshal(a)
}
func DecodeAudit(data []byte) (AuditEvent, error) {
	var a AuditEvent
	e := json.Unmarshal(data, &a)
	return a, e
}
