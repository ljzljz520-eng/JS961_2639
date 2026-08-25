package repository

import (
	"strings"
	"weddingtemplates/internal/domain"
)

func FilterTitle(records []domain.Record, title string) []domain.Record {
	out := []domain.Record{}
	title = strings.ToLower(title)
	for _, r := range records {
		if strings.Contains(strings.ToLower(r.Title), title) {
			out = append(out, r)
		}
	}
	return out
}
func FilterScore(records []domain.Record, min, max int) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if r.Score >= min && r.Score <= max {
			out = append(out, r)
		}
	}
	return out
}
func GroupStatus(records []domain.Record) map[string][]domain.Record {
	out := map[string][]domain.Record{}
	for _, r := range records {
		out[r.Status] = append(out[r.Status], r)
	}
	return out
}
func GroupRegion(records []domain.Record) map[string][]domain.Record {
	out := map[string][]domain.Record{}
	for _, r := range records {
		out[r.Region] = append(out[r.Region], r)
	}
	return out
}
func ContainsID(records []domain.Record, id string) bool {
	for _, r := range records {
		if r.ID == id {
			return true
		}
	}
	return false
}
func MergeRecords(primary, secondary []domain.Record) []domain.Record {
	out := append([]domain.Record{}, primary...)
	for _, x := range secondary {
		if !ContainsID(out, x.ID) {
			out = append(out, x)
		}
	}
	return out
}
func RemoveArchived(records []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if r.Status != domain.StatusArchived {
			out = append(out, r)
		}
	}
	return out
}
