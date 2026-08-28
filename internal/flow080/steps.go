package flow080

import (
	"fmt"
	"weddingtemplates/internal/domain"
	"weddingtemplates/internal/service"
)

type StepResult struct {
	Step     string
	Complete bool
	Message  string
}

func step(name string, ok bool, msg string) StepResult {
	return StepResult{Step: name, Complete: ok, Message: msg}
}
func RegisterStep(s *service.Service, id, title string) (StepResult, error) {
	_, e := s.Register(id, title, "", "", 70)
	if e != nil {
		return step("register", false, e.Error()), e
	}
	return step("register", true, "record created"), nil
}
func ReviewStep(s *service.Service, id string) (StepResult, error) {
	_, e := s.Review(id, "reviewer")
	if e != nil {
		return step("review", false, e.Error()), e
	}
	return step("review", true, "awaiting approval"), nil
}
func ApproveStep(s *service.Service, id string) (StepResult, error) {
	_, e := s.Approve(id, "approver")
	if e != nil {
		return step("approve", false, e.Error()), e
	}
	return step("approve", true, "approved"), nil
}
func ArchiveStep(s *service.Service, id string) (StepResult, error) {
	_, e := s.Archive(id, "archiver")
	if e != nil {
		return step("archive", false, e.Error()), e
	}
	return step("archive", true, "archived"), nil
}
func ImportStep(s *service.Service, rows []domain.ImportRow) (StepResult, error) {
	r, e := s.Import(rows, "importer")
	if e != nil {
		return step("import", false, e.Error()), e
	}
	if len(r.Accepted) == 0 {
		return step("import", false, "empty"), fmt.Errorf("empty import")
	}
	return step("import", true, fmt.Sprintf("%d accepted", len(r.Accepted))), nil
}
