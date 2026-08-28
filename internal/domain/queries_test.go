package domain

import (
	"testing"
	"time"
)

func TestQueryMatching(t *testing.T) {
	r := NewRecord("1", "Summer Garden", "Pearl Hall", "East", 88, time.Unix(1, 0))
	if !Match(r, Query{Text: "garden", MinScore: 80}) {
		t.Fatal("match")
	}
	if Match(r, Query{Region: "West"}) {
		t.Fatal("region")
	}
}
