package repository

import (
	"time"
	"weddingtemplates/internal/domain"
)

func (r *Repository) AddAttachment(a domain.Attachment) error       { return r.Store.PutAttachment(a) }
func (r *Repository) Audits(id string) ([]domain.AuditEvent, error) { return r.Store.Audits(id) }
func (r *Repository) Attachments(id string) ([]domain.Attachment, error) {
	return r.Store.Attachments(id)
}
func (r *Repository) BeginWorkflow(id, owner string, now time.Time) error {
	return r.Store.PutWorkflow(domain.Workflow{ID: "workflow-" + id, RecordID: id, State: domain.StatusDraft, Owner: owner, Revision: 1, UpdatedAt: now})
}
func (r *Repository) Workflow(id string) (domain.Workflow, error) {
	return r.Store.Workflow("workflow-" + id)
}
