package service

import (
	"fmt"
	"sort"
	"weddingtemplates/internal/domain"
)

type Notification struct {
	Recipient, Subject, Body string
	Priority                 int
}

func BuildNotification(r domain.Record, action string) Notification {
	subject := fmt.Sprintf("Template %s: %s", r.ID, action)
	body := fmt.Sprintf("%s is now %s with score %d", r.Title, r.Status, r.Score)
	priority := 1
	if r.Status == domain.StatusApproved {
		priority = 2
	}
	if r.Status == domain.StatusArchived {
		priority = 3
	}
	return Notification{Recipient: r.Region, Subject: subject, Body: body, Priority: priority}
}
func SortNotifications(items []Notification) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].Subject < items[j].Subject
		}
		return items[i].Priority > items[j].Priority
	})
}
func (s *Service) PendingNotifications() ([]Notification, error) {
	p, e := s.Search(domain.Query{Status: domain.StatusReview, Limit: 100000})
	if e != nil {
		return nil, e
	}
	out := make([]Notification, 0, len(p.Items))
	for _, r := range p.Items {
		out = append(out, BuildNotification(r, "review required"))
	}
	SortNotifications(out)
	return out, nil
}
func (s *Service) ApprovalNotification(id string) (Notification, error) {
	r, e := s.Repo.FindRecord(id)
	if e != nil {
		return Notification{}, e
	}
	return BuildNotification(r, "approval update"), nil
}
