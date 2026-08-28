package service

import (
	"sort"
	"weddingtemplates/internal/domain"
)

type RegionStat struct {
	Region  string
	Count   int
	Average float64
}

func (s *Service) RegionStats() ([]RegionStat, error) {
	p, e := s.Search(domain.Query{Limit: 100000})
	if e != nil {
		return nil, e
	}
	groups := map[string][]domain.Record{}
	for _, r := range p.Items {
		groups[r.Region] = append(groups[r.Region], r)
	}
	out := make([]RegionStat, 0, len(groups))
	for region, list := range groups {
		total := 0
		for _, r := range list {
			total += r.Score
		}
		out = append(out, RegionStat{Region: region, Count: len(list), Average: float64(total) / float64(len(list))})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Region < out[j].Region })
	return out, nil
}
func (s *Service) TopTemplates(limit int) ([]domain.Record, error) {
	if limit < 1 {
		limit = 10
	}
	p, e := s.Search(domain.Query{Limit: 100000})
	if e != nil {
		return nil, e
	}
	domain.SortRecords(p.Items, "score", true)
	if limit > len(p.Items) {
		limit = len(p.Items)
	}
	return p.Items[:limit], nil
}
func (s *Service) ArchiveEligible(policy domain.ReviewPolicy) (int, error) {
	p, e := s.Search(domain.Query{Status: domain.StatusApproved, Limit: 100000})
	if e != nil {
		return 0, e
	}
	n := 0
	for _, r := range p.Items {
		if policy.Eligible(r) {
			if _, e = s.Archive(r.ID, "archiver"); e != nil {
				return n, e
			}
			n++
		}
	}
	return n, nil
}
