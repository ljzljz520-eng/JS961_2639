package service

import (
	"fmt"
	"sort"
	"weddingtemplates/internal/domain"
)

type BatchOptions struct {
	Actor       string
	Strict      bool
	Deduplicate bool
}
type BatchOutcome struct {
	Accepted   []domain.Record
	Rejected   []string
	Duplicates []string
	Warnings   []string
}

func (s *Service) ProcessBatch(rows []domain.ImportRow, opt BatchOptions) (BatchOutcome, error) {
	if opt.Actor == "" {
		opt.Actor = "batch"
	}
	out := BatchOutcome{}
	seen := map[string]bool{}
	for _, row := range rows {
		if opt.Deduplicate && seen[row.ExternalID] {
			out.Duplicates = append(out.Duplicates, row.ExternalID)
			continue
		}
		seen[row.ExternalID] = true
		r := domain.NewRecord(row.ExternalID, row.Title, row.Venue, row.Region, row.Score, s.Clock.Now())
		if e := r.Validate(); e != nil {
			out.Rejected = append(out.Rejected, row.ExternalID)
			if opt.Strict {
				return out, fmt.Errorf("row %s rejected: %w", row.ExternalID, e)
			}
			continue
		}
		out.Accepted = append(out.Accepted, r)
	}
	report, e := s.Repo.SaveValidated(out.Accepted)
	if e != nil {
		return out, e
	}
	if report.Saved != len(out.Accepted) {
		out.Warnings = append(out.Warnings, "some records were skipped")
	}
	return out, nil
}
func SortBatch(out *BatchOutcome) {
	sort.SliceStable(out.Accepted, func(i, j int) bool {
		if out.Accepted[i].Score == out.Accepted[j].Score {
			return out.Accepted[i].ID < out.Accepted[j].ID
		}
		return out.Accepted[i].Score > out.Accepted[j].Score
	})
}
func (out BatchOutcome) AcceptedCount() int { return len(out.Accepted) }
func (out BatchOutcome) RejectedCount() int { return len(out.Rejected) }
func (out BatchOutcome) HasErrors() bool    { return len(out.Rejected) > 0 || len(out.Duplicates) > 0 }
func (s *Service) ReplaceFromBatch(ids []string, score int) error {
	for _, id := range ids {
		if _, e := s.Repo.ReplaceScore(id, score, s.Clock.Now()); e != nil {
			return e
		}
	}
	return nil
}
