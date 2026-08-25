package flow080

import (
	"fmt"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/service"
)

type Plan struct {
	ID     string
	Steps  []string
	Ready  bool
	Issues []string
}

func BuildPlan(s *service.Service, id string) Plan {
	p := Plan{ID: id, Steps: []string{"register", "review", "approve", "archive"}}
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		p.Issues = append(p.Issues, "record not found")
		return p
	}
	if r.Status == domain.StatusArchived {
		p.Ready = true
		return p
	}
	if r.Score < 60 {
		p.Issues = append(p.Issues, "score threshold")
	}
	p.Ready = len(p.Issues) == 0
	return p
}
func ExecutePlan(s *service.Service, p Plan) (domain.Record, error) {
	if !p.Ready {
		return domain.Record{}, fmt.Errorf("plan is not ready")
	}
	return RunCreateReviewArchive(s, p.ID)
}
func PlanForImport(rows []domain.ImportRow) Plan {
	p := Plan{ID: "import", Steps: []string{"import", "validate", "persist", "report"}}
	for _, r := range rows {
		if r.ExternalID == "" || r.Title == "" {
			p.Issues = append(p.Issues, "identity")
		}
	}
	p.Ready = len(p.Issues) == 0
	return p
}
