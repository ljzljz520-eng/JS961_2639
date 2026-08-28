package repository

import (
	"time"
	"weddingtemplates/internal/domain"
)

type BulkReport struct {
	Saved, Skipped int
	IDs            []string
}

func (r *Repository) SaveValidated(records []domain.Record) (BulkReport, error) {
	out := BulkReport{}
	for _, x := range records {
		if e := x.Validate(); e != nil {
			out.Skipped++
			continue
		}
		if e := r.SaveRecord(x); e != nil {
			return out, e
		}
		out.Saved++
		out.IDs = append(out.IDs, x.ID)
	}
	return out, nil
}
func (r *Repository) Touch(id string, now time.Time) error {
	x, e := r.FindRecord(id)
	if e != nil {
		return e
	}
	x.UpdatedAt = now
	return r.SaveRecord(x)
}
func (r *Repository) ReplaceScore(id string, score int, now time.Time) (domain.Record, error) {
	x, e := r.FindRecord(id)
	if e != nil {
		return x, e
	}
	x.Score = score
	x.Version++
	x.UpdatedAt = now
	if e = x.Validate(); e != nil {
		return x, e
	}
	return x, r.SaveRecord(x)
}
