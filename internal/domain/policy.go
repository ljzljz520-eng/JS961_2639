package domain

import "fmt"

type ReviewPolicy struct {
	MinimumScore   int
	RequiredRegion string
	RequireVenue   bool
}

func (p ReviewPolicy) Evaluate(r Record) []string {
	var reasons []string
	if r.Score < p.MinimumScore {
		reasons = append(reasons, fmt.Sprintf("score below %d", p.MinimumScore))
	}
	if p.RequiredRegion != "" && r.Region != p.RequiredRegion {
		reasons = append(reasons, "region mismatch")
	}
	if p.RequireVenue && r.Venue == "" {
		reasons = append(reasons, "venue required")
	}
	return reasons
}
func (p ReviewPolicy) Eligible(r Record) bool { return len(p.Evaluate(r)) == 0 }

type Permission string

const (
	PermissionRegister Permission = "register"
	PermissionReview   Permission = "review"
	PermissionApprove  Permission = "approve"
	PermissionArchive  Permission = "archive"
)

func Allowed(actor string, p Permission) bool {
	switch actor {
	case "registrar":
		return p == PermissionRegister
	case "reviewer":
		return p == PermissionReview
	case "approver":
		return p == PermissionApprove
	case "archiver":
		return p == PermissionArchive
	case "api":
		return true
	}
	return false
}
