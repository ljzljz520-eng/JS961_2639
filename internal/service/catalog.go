package service

import (
	"fmt"
	"weddingtemplates/internal/domain"
)

func (s *Service) RegisterWithCatalog(c domain.Catalog, id, title, venue, region string, score int) (domain.Record, error) {
	r := domain.NewRecord(id, title, venue, region, score, s.Clock.Now())
	if !c.ValidRecord(r) {
		return r, fmt.Errorf("record does not satisfy catalog")
	}
	if e := s.Repo.SaveRecord(r); e != nil {
		return r, e
	}
	return r, nil
}
func (s *Service) ReviewWithPolicy(p domain.ReviewPolicy, id, actor string) (domain.Record, error) {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return r, e
	}
	if reasons := p.Evaluate(r); len(reasons) > 0 {
		return r, fmt.Errorf("review rejected: %v", reasons)
	}
	return s.Review(id, actor)
}
func (s *Service) Catalog() domain.Catalog { return domain.DefaultCatalog() }
