package domain

import (
	"testing"
	"time"
)

func TestRecordValidation(t *testing.T) {
	r := NewRecord("1", "Garden", "Hall", "East", 80, time.Unix(1, 0))
	if e := r.Validate(); e != nil {
		t.Fatal(e)
	}
	if !CanTransition(StatusDraft, StatusReview) {
		t.Fatal("transition")
	}
}
