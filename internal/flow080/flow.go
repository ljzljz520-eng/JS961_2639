package flow080

import (
	"fmt"
	"time"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/service"
)

func RunCreateReviewArchive(s *service.Service, id string) (domain.Record, error) {
	if _, e := s.Review(id, "reviewer"); e != nil {
		return domain.Record{}, e
	}
	if _, e := s.Approve(id, "approver"); e != nil {
		return domain.Record{}, e
	}
	return s.Archive(id, "archiver")
}
func RunSearchUpdatePublish(s *service.Service, q domain.Query, id string) (domain.Page, domain.Record, error) {
	p, e := s.Search(q)
	if e != nil {
		return p, domain.Record{}, e
	}
	if _, e = s.Review(id, "reviewer"); e != nil {
		return p, domain.Record{}, e
	}
	r, e := s.Revise(id, "Updated template", "Updated venue", "updated", 88)
	if e != nil {
		return p, r, e
	}
	approved, e := s.Approve(id, "publisher")
	if e != nil {
		return p, r, e
	}
	return p, approved, nil
}
func RunImportReport(s *service.Service, rows []domain.ImportRow) (domain.ImportResult, service.Summary, error) {
	r, e := s.Import(rows, "importer")
	if e != nil {
		return r, service.Summary{}, e
	}
	sum, e := s.Summary()
	return r, sum, e
}
func DeterministicClock() service.FixedClock {
	return service.FixedClock{T: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}
}
func ValidateBatch(result domain.ImportResult) error {
	if len(result.Accepted) == 0 {
		return fmt.Errorf("batch has no accepted records")
	}
	return nil
}
