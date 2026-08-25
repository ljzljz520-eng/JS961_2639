package domain

import (
	"fmt"
	"strings"
	"time"
)

func FormatRecord(r Record) string {
	return fmt.Sprintf("%s | %s | %s | %s | %d | %s", r.ID, r.Title, r.Venue, r.Region, r.Score, r.Status)
}
func FormatAudit(a AuditEvent) string {
	return fmt.Sprintf("%s %s %s %s", a.At.UTC().Format(time.RFC3339), a.Action, a.RecordID, a.Actor)
}
func ParseStatus(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case StatusDraft, StatusReview, StatusApproved, StatusArchived:
		return s, nil
	}
	return "", fmt.Errorf("unknown status %q", s)
}
func NormalizeRecord(r Record) Record {
	r.ID = strings.TrimSpace(r.ID)
	r.Title = strings.TrimSpace(r.Title)
	r.Venue = strings.TrimSpace(r.Venue)
	r.Region = strings.TrimSpace(r.Region)
	r.Status = strings.ToLower(strings.TrimSpace(r.Status))
	return r
}
func CloneRecord(r Record) Record {
	return Record{ID: r.ID, Title: r.Title, Venue: r.Venue, Region: r.Region, Score: r.Score, Status: r.Status, Version: r.Version, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}
}
func CloneRecords(in []Record) []Record {
	out := make([]Record, len(in))
	for i, r := range in {
		out[i] = CloneRecord(r)
	}
	return out
}
func SameContent(a, b Record) bool {
	return a.ID == b.ID && a.Title == b.Title && a.Venue == b.Venue && a.Region == b.Region && a.Score == b.Score && a.Status == b.Status
}
func IsTerminal(status string) bool    { return status == StatusArchived }
func IsPublishable(status string) bool { return status == StatusApproved }
func NextStatus(status string) string {
	switch status {
	case StatusDraft:
		return StatusReview
	case StatusReview:
		return StatusApproved
	case StatusApproved:
		return StatusArchived
	}
	return ""
}
