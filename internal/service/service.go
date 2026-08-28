package service

import (
	"fmt"
	"time"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/repository"
)

type Clock interface{ Now() time.Time }
type FixedClock struct{ T time.Time }

func (c FixedClock) Now() time.Time { return c.T }

type Service struct {
	Repo  *repository.Repository
	Clock Clock
}

func New(r *repository.Repository, c Clock) *Service { return &Service{Repo: r, Clock: c} }
func (s *Service) Register(id, title, venue, region string, score int) (domain.Record, error) {
	r := domain.NewRecord(id, title, venue, region, score, s.Clock.Now())
	if e := r.Validate(); e != nil {
		return r, e
	}
	if e := s.Repo.SaveRecord(r); e != nil {
		return r, e
	}
	_ = s.Repo.BeginWorkflow(id, "registrar", s.Clock.Now())
	return r, nil
}
func (s *Service) Review(id, actor string) (domain.Record, error) {
	return s.Repo.Transition(id, domain.StatusReview, actor, "submitted for review", s.Clock.Now())
}
func (s *Service) Approve(id, actor string) (domain.Record, error) {
	return s.Repo.Transition(id, domain.StatusApproved, actor, "approved", s.Clock.Now())
}
func (s *Service) Archive(id, actor string) (domain.Record, error) {
	return s.Repo.Transition(id, domain.StatusArchived, actor, "archived", s.Clock.Now())
}
func (s *Service) Revise(id, title, venue, region string, score int) (domain.Record, error) {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return r, e
	}
	if r.Status == domain.StatusArchived {
		return r, fmt.Errorf("archived records cannot change")
	}
	r.Title = title
	r.Venue = venue
	r.Region = region
	r.Score = score
	r.Version++
	r.UpdatedAt = s.Clock.Now()
	if e = r.Validate(); e != nil {
		return r, e
	}
	return r, s.Repo.SaveRecord(r)
}
func (s *Service) Search(q domain.Query) (domain.Page, error) { return s.Repo.Search(q) }
func (s *Service) Import(rows []domain.ImportRow, actor string) (domain.ImportResult, error) {
	result := domain.BuildImport(rows, actor, s.Clock.Now())
	if e := s.Repo.SaveImport(result); e != nil {
		return result, e
	}
	return result, nil
}
func (s *Service) Attach(id, name, media, checksum string, size int64) error {
	return s.Repo.AddAttachment(domain.Attachment{ID: id + "/" + name, RecordID: id, Name: name, MediaType: media, Checksum: checksum, Size: size})
}
func (s *Service) History(id string) ([]domain.AuditEvent, error) { return s.Repo.Audits(id) }
