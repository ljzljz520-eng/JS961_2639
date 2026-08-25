package flow080

import (
	"time"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/service"
)

type BatchHandler struct {
	Svc *service.Service
	now time.Time
}

func NewBatchHandler(s *service.Service, now time.Time) *BatchHandler {
	return &BatchHandler{Svc: s, now: now}
}
func (h *BatchHandler) Sync(rows []domain.ImportRow) (domain.ImportResult, error) {
	result := domain.ImportResult{}
	for _, row := range rows {
		built := domain.BuildImport([]domain.ImportRow{{ExternalID: row.ExternalID, Title: row.Title, Venue: row.Venue, Region: row.Region, Score: row.Score}}, "sync", h.now)
		result.Accepted = append(result.Accepted, built.Accepted...)
		result.Rejected = append(result.Rejected, built.Rejected...)
		result.Events = append(result.Events, built.Events...)
	}
	if e := h.Svc.Repo.SaveImport(result); e != nil {
		return result, e
	}
	return result, nil
}
