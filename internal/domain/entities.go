package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidRecord     = errors.New("invalid wedding template record")
	ErrInvalidTransition = errors.New("invalid workflow transition")
)

type Record struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Venue     string    `json:"venue"`
	Region    string    `json:"region"`
	Score     int       `json:"score"`
	Status    string    `json:"status"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuditEvent struct {
	ID, RecordID, Action, Actor, Detail string
	At                                  time.Time
}
type Workflow struct {
	ID, RecordID, State, Owner string
	Revision                   int
	UpdatedAt                  time.Time
}
type Attachment struct {
	ID, RecordID, Name, MediaType string
	Size                          int64
	Checksum                      string
}

const (
	StatusDraft    = "draft"
	StatusReview   = "review"
	StatusApproved = "approved"
	StatusArchived = "archived"
)

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.Title) == "" {
		return ErrInvalidRecord
	}
	if r.Score < 0 || r.Score > 100 {
		return fmt.Errorf("score must be between 0 and 100: %w", ErrInvalidRecord)
	}
	if r.Status == "" {
		return fmt.Errorf("status is required: %w", ErrInvalidRecord)
	}
	return nil
}

func CanTransition(from, to string) bool {
	switch from {
	case StatusDraft:
		return to == StatusReview
	case StatusReview:
		return to == StatusApproved || to == StatusDraft
	case StatusApproved:
		return to == StatusArchived
	case StatusArchived:
		return false
	}
	return false
}

func (r *Record) Transition(to string, now time.Time) error {
	if !CanTransition(r.Status, to) {
		return ErrInvalidTransition
	}
	r.Status = to
	r.Version++
	r.UpdatedAt = now
	return nil
}

func (r Record) Key() string { return r.ID }
func NewRecord(id, title, venue, region string, score int, now time.Time) Record {
	return Record{ID: id, Title: title, Venue: venue, Region: region, Score: score, Status: StatusDraft, Version: 1, CreatedAt: now, UpdatedAt: now}
}
func (a AuditEvent) Key() string { return a.ID }
func (w Workflow) Key() string   { return w.ID }
func (a Attachment) Key() string { return a.ID }
