package service

import (
	"fmt"
	"weddingtemplates/internal/domain"
)

type ValidationReport struct {
	Valid      bool
	Violations []domain.Violation
}

func (s *Service) ValidateRecord(r domain.Record) ValidationReport {
	rules := domain.DefaultRuleSet()
	v := rules.Check(r)
	return ValidationReport{Valid: len(v) == 0, Violations: v}
}
func (s *Service) ValidateIDs(ids []string) error {
	seen := map[string]bool{}
	for _, id := range ids {
		if id == "" {
			return fmt.Errorf("empty id")
		}
		if seen[id] {
			return fmt.Errorf("duplicate id %s", id)
		}
		seen[id] = true
	}
	return nil
}
func (s *Service) EnsureEditable(id string) (domain.Record, error) {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return r, e
	}
	if domain.IsTerminal(r.Status) {
		return r, fmt.Errorf("record %s is archived", id)
	}
	return r, nil
}
func (s *Service) EnsurePublishable(id string) (domain.Record, error) {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return r, e
	}
	if !domain.IsPublishable(r.Status) {
		return r, fmt.Errorf("record %s is not approved", id)
	}
	return r, nil
}
func (s *Service) UpdateScore(id string, score int) (domain.Record, error) {
	r, e := s.EnsureEditable(id)
	if e != nil {
		return r, e
	}
	r.Score = score
	r.Version++
	r.UpdatedAt = s.Clock.Now()
	if e = r.Validate(); e != nil {
		return r, e
	}
	return r, s.Repo.SaveRecord(r)
}
func (s *Service) UpdateVenue(id, venue string) (domain.Record, error) {
	r, e := s.EnsureEditable(id)
	if e != nil {
		return r, e
	}
	r.Venue = venue
	r.Version++
	r.UpdatedAt = s.Clock.Now()
	if e = r.Validate(); e != nil {
		return r, e
	}
	return r, s.Repo.SaveRecord(r)
}
