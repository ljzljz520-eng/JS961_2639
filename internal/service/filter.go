package service

import (
	"strings"
	"weddingtemplates/internal/domain"
)

type FilterPreset struct {
	Name  string
	Query domain.Query
}

func NormalizeQuery(q domain.Query) domain.Query {
	q.Text = strings.TrimSpace(q.Text)
	q.Region = strings.TrimSpace(q.Region)
	q.Venue = strings.TrimSpace(q.Venue)
	if q.Limit < 1 {
		q.Limit = 50
	}
	if q.Limit > 200 {
		q.Limit = 200
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
	return q
}
func (s *Service) SearchPreset(p FilterPreset) (domain.Page, error) {
	return s.Search(NormalizeQuery(p.Query))
}
func ScoreBand(min, max int) domain.Query {
	if min < 0 {
		min = 0
	}
	if max > 100 {
		max = 100
	}
	if max < min {
		max = min
	}
	return domain.Query{MinScore: min, Limit: 200}
}
func InRegions(regions []string) domain.Query {
	q := domain.Query{Limit: 200}
	if len(regions) > 0 {
		q.Region = regions[0]
	}
	return q
}
