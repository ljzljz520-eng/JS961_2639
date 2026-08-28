package service

import (
	"sort"
	"weddingtemplates/internal/domain"
)

type Summary struct {
	Total, Draft, Review, Approved, Archived int
	AverageScore                             float64
	Regions                                  []string
}

func (s *Service) Summary() (Summary, error) {
	p, e := s.Search(domain.Query{Limit: 100000})
	if e != nil {
		return Summary{}, e
	}
	var out Summary
	seen := map[string]bool{}
	for _, r := range p.Items {
		out.Total++
		switch r.Status {
		case domain.StatusDraft:
			out.Draft++
		case domain.StatusReview:
			out.Review++
		case domain.StatusApproved:
			out.Approved++
		case domain.StatusArchived:
			out.Archived++
		}
		out.AverageScore += float64(r.Score)
		if !seen[r.Region] {
			seen[r.Region] = true
			out.Regions = append(out.Regions, r.Region)
		}
	}
	if out.Total > 0 {
		out.AverageScore /= float64(out.Total)
	}
	sort.Strings(out.Regions)
	return out, nil
}
func (s *Service) Publishable(id string) (bool, error) {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return false, e
	}
	return r.Status == domain.StatusApproved, nil
}
