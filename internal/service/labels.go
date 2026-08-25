package service

import (
	"strings"
	"weddingtemplates/internal/domain"
)

func LabelForStatus(status string) string {
	switch status {
	case domain.StatusDraft:
		return "Draft"
	case domain.StatusReview:
		return "In review"
	case domain.StatusApproved:
		return "Approved"
	case domain.StatusArchived:
		return "Archived"
	}
	return "Unknown"
}
func Labels(records []domain.Record) map[string]string {
	out := map[string]string{}
	for _, r := range records {
		out[r.ID] = LabelForStatus(r.Status)
	}
	return out
}
func NormalizeLabel(s string) string { return strings.TrimSpace(strings.ToLower(s)) }
