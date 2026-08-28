package repository

import (
	"time"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/storage"
)

type Repository struct{ Store *storage.Store }

func New(s *storage.Store) *Repository                            { return &Repository{Store: s} }
func (r *Repository) SaveRecord(x domain.Record) error            { return r.Store.PutRecord(x) }
func (r *Repository) FindRecord(id string) (domain.Record, error) { return r.Store.Record(id) }
func (r *Repository) Search(q domain.Query) (domain.Page, error) {
	all, e := r.Store.ListRecords()
	if e != nil {
		return domain.Page{}, e
	}
	filtered := make([]domain.Record, 0, len(all))
	for _, x := range all {
		if domain.Match(x, q) {
			filtered = append(filtered, x)
		}
	}
	domain.SortRecords(filtered, "updated", true)
	total := len(filtered)
	off := q.Offset
	if off < 0 {
		off = 0
	}
	if off > total {
		off = total
	}
	lim := q.Limit
	if lim <= 0 {
		lim = 50
	}
	end := off + lim
	if end > total {
		end = total
	}
	return domain.Page{Items: filtered[off:end], Total: total, Offset: off, Limit: lim}, nil
}
func (r *Repository) Transition(id, to, actor, detail string, now time.Time) (domain.Record, error) {
	x, e := r.FindRecord(id)
	if e != nil {
		return x, e
	}
	if e = x.Transition(to, now); e != nil {
		return x, e
	}
	if e = r.SaveRecord(x); e != nil {
		return x, e
	}
	_ = r.Store.PutAudit(domain.AuditEvent{ID: id + "-" + to + "-" + string(rune(x.Version)), RecordID: id, Action: to, Actor: actor, Detail: detail, At: now})
	return x, nil
}
func (r *Repository) SaveImport(result domain.ImportResult) error {
	for _, x := range result.Accepted {
		if e := r.SaveRecord(x); e != nil {
			return e
		}
	}
	for _, a := range result.Events {
		if e := r.Store.PutAudit(a); e != nil {
			return e
		}
	}
	return nil
}
