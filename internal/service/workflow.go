package service

import (
	"fmt"
	"weddingtemplates/internal/domain"
)

type WorkflowView struct {
	Record      domain.Record
	History     []domain.AuditEvent
	Publishable bool
}

func (s *Service) ViewWorkflow(id string) (WorkflowView, error) {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return WorkflowView{}, e
	}
	h, e := s.History(id)
	if e != nil {
		return WorkflowView{}, e
	}
	return WorkflowView{Record: r, History: h, Publishable: r.Status == domain.StatusApproved}, nil
}
func (s *Service) RequireState(id, state string) error {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return e
	}
	if r.Status != state {
		return fmt.Errorf("expected %s, got %s", state, r.Status)
	}
	return nil
}
func (s *Service) Advance(id, actor string) (domain.Record, error) {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return r, e
	}
	next := domain.NextStatus(r.Status)
	if next == "" {
		return r, fmt.Errorf("no next state")
	}
	return s.Repo.Transition(id, next, actor, "advanced", s.Clock.Now())
}
func (s *Service) Replay(id string) ([]string, error) {
	events, e := s.History(id)
	if e != nil {
		return nil, e
	}
	out := make([]string, 0, len(events))
	for _, a := range events {
		out = append(out, a.Action)
	}
	return out, nil
}
func (s *Service) CanDelete(id string) bool {
	r, e := s.Repo.FindRecord(id)
	return e == nil && r.Status == domain.StatusDraft
}
func (s *Service) DeleteDraft(id string) error {
	if !s.CanDelete(id) {
		return fmt.Errorf("only drafts can be deleted")
	}
	return s.Repo.Remove(id)
}
