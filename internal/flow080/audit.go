package flow080

import (
	"sort"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/service"
)

type AuditSummary struct {
	Actions []string
	Count   int
	Latest  string
}

func SummarizeAudit(events []domain.AuditEvent) AuditSummary {
	sort.SliceStable(events, func(i, j int) bool { return events[i].At.Before(events[j].At) })
	out := AuditSummary{Count: len(events)}
	seen := map[string]bool{}
	for _, e := range events {
		if !seen[e.Action] {
			out.Actions = append(out.Actions, e.Action)
			seen[e.Action] = true
		}
		out.Latest = e.Action
	}
	return out
}
func AuditWorkflow(s *service.Service, id string) (AuditSummary, error) {
	events, e := s.History(id)
	if e != nil {
		return AuditSummary{}, e
	}
	return SummarizeAudit(events), nil
}
func AuditComplete(summary AuditSummary) bool    { return summary.Count >= 3 && summary.Latest != "" }
func AuditActions(summary AuditSummary) []string { return append([]string{}, summary.Actions...) }
