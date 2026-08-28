package service

import (
	"fmt"
	"weddingtemplates/internal/domain"
)

type ImportReport struct {
	Accepted, Rejected int
	RejectedIDs        []string
	AverageScore       float64
}

func (s *Service) ImportReport(rows []domain.ImportRow, actor string) (ImportReport, error) {
	result, e := s.Import(rows, actor)
	if e != nil {
		return ImportReport{}, e
	}
	out := ImportReport{Accepted: len(result.Accepted), Rejected: len(result.Rejected), RejectedIDs: result.Rejected}
	for _, r := range result.Accepted {
		out.AverageScore += float64(r.Score)
	}
	if out.Accepted > 0 {
		out.AverageScore /= float64(out.Accepted)
	}
	if out.Accepted == 0 {
		return out, fmt.Errorf("no valid rows")
	}
	return out, nil
}
func (s *Service) ValidateImport(rows []domain.ImportRow) []string {
	var bad []string
	for _, r := range rows {
		if r.ExternalID == "" || r.Title == "" || r.Score < 0 || r.Score > 100 {
			bad = append(bad, r.ExternalID)
		}
	}
	return bad
}
