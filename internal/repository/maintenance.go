package repository

import (
	"strings"
	"time"
	"weddingtemplates/internal/domain"
)

type LifecycleReport struct {
	RecordID   string
	Current    string
	Events     int
	CanPublish bool
}

func (r *Repository) Lifecycle(id string) (LifecycleReport, error) {
	x, e := r.FindRecord(id)
	if e != nil {
		return LifecycleReport{}, e
	}
	events, e := r.Audits(id)
	if e != nil {
		return LifecycleReport{}, e
	}
	return LifecycleReport{RecordID: id, Current: x.Status, Events: len(events), CanPublish: x.Status == domain.StatusApproved}, nil
}
func (r *Repository) SearchByIDs(ids []string) ([]domain.Record, error) {
	out := []domain.Record{}
	for _, id := range ids {
		x, e := r.FindRecord(id)
		if e != nil {
			continue
		}
		out = append(out, x)
	}
	return out, nil
}
func (r *Repository) FindByPrefix(prefix string) ([]domain.Record, error) {
	all, e := r.Store.ListRecords()
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, x := range all {
		if strings.HasPrefix(x.ID, prefix) {
			out = append(out, x)
		}
	}
	return out, nil
}
func (r *Repository) MarkReviewed(id, actor string, now time.Time) (domain.Record, error) {
	return r.Transition(id, domain.StatusReview, actor, "reviewed", now)
}
func (r *Repository) MarkApproved(id, actor string, now time.Time) (domain.Record, error) {
	return r.Transition(id, domain.StatusApproved, actor, "approved", now)
}
func (r *Repository) MarkArchived(id, actor string, now time.Time) (domain.Record, error) {
	return r.Transition(id, domain.StatusArchived, actor, "archived", now)
}
func (r *Repository) AuditCount(id string) (int, error) { a, e := r.Audits(id); return len(a), e }
