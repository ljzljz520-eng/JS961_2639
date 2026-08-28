package domain

import (
	"fmt"
	"strings"
)

type RuleCode string

const (
	RuleIdentity RuleCode = "identity"
	RuleScore    RuleCode = "score"
	RuleRegion   RuleCode = "region"
	RuleStatus   RuleCode = "status"
	RuleVersion  RuleCode = "version"
	RuleTitle    RuleCode = "title"
	RuleVenue    RuleCode = "venue"
)

type Violation struct {
	Code    RuleCode
	Field   string
	Message string
}
type RuleSet struct {
	RequireRegion, RequireVenue, RequireTitle bool
	MinScore, MaxScore                        int
	AllowedStatuses                           []string
}

func DefaultRuleSet() RuleSet {
	return RuleSet{RequireRegion: true, RequireVenue: true, RequireTitle: true, MinScore: 0, MaxScore: 100, AllowedStatuses: []string{StatusDraft, StatusReview, StatusApproved, StatusArchived}}
}
func (r RuleSet) Check(x Record) []Violation {
	out := []Violation{}
	if r.RequireTitle && strings.TrimSpace(x.Title) == "" {
		out = append(out, Violation{RuleTitle, "title", "title required"})
	}
	if r.RequireRegion && strings.TrimSpace(x.Region) == "" {
		out = append(out, Violation{RuleRegion, "region", "region required"})
	}
	if r.RequireVenue && strings.TrimSpace(x.Venue) == "" {
		out = append(out, Violation{RuleVenue, "venue", "venue required"})
	}
	if x.Score < r.MinScore || x.Score > r.MaxScore {
		out = append(out, Violation{RuleScore, "score", fmt.Sprintf("score %d outside range", x.Score)})
	}
	ok := false
	for _, status := range r.AllowedStatuses {
		if status == x.Status {
			ok = true
			break
		}
	}
	if !ok {
		out = append(out, Violation{RuleStatus, "status", "unknown status"})
	}
	if x.Version < 1 {
		out = append(out, Violation{RuleVersion, "version", "version must be positive"})
	}
	return out
}
func (r RuleSet) Valid(x Record) bool { return len(r.Check(x)) == 0 }
func (r RuleSet) Explain(x Record) string {
	v := r.Check(x)
	parts := make([]string, 0, len(v))
	for _, item := range v {
		parts = append(parts, item.Field+": "+item.Message)
	}
	return strings.Join(parts, "; ")
}
func MergeRules(base, override RuleSet) RuleSet {
	out := base
	if override.MinScore != 0 {
		out.MinScore = override.MinScore
	}
	if override.MaxScore != 0 {
		out.MaxScore = override.MaxScore
	}
	if override.AllowedStatuses != nil {
		out.AllowedStatuses = append([]string{}, override.AllowedStatuses...)
	}
	if override.RequireRegion {
		out.RequireRegion = true
	}
	if override.RequireVenue {
		out.RequireVenue = true
	}
	if override.RequireTitle {
		out.RequireTitle = true
	}
	return out
}
func StatusOrder(status string) int {
	switch status {
	case StatusDraft:
		return 1
	case StatusReview:
		return 2
	case StatusApproved:
		return 3
	case StatusArchived:
		return 4
	}
	return 0
}
func CompareStatus(a, b string) int {
	aa, bb := StatusOrder(a), StatusOrder(b)
	if aa < bb {
		return -1
	}
	if aa > bb {
		return 1
	}
	return 0
}
