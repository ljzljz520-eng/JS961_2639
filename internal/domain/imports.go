package domain

import "time"

type ImportRow struct {
	ExternalID, Title, Venue, Region string
	Score                            int
}
type ImportResult struct {
	Accepted []Record
	Rejected []string
	Events   []AuditEvent
}

func BuildImport(rows []ImportRow, actor string, now time.Time) ImportResult {
	out := ImportResult{}
	for _, row := range rows {
		if row.ExternalID == "" || row.Title == "" || row.Score < 0 || row.Score > 100 {
			out.Rejected = append(out.Rejected, row.ExternalID)
			continue
		}
		r := NewRecord(row.ExternalID, row.Title, row.Venue, row.Region, row.Score, now)
		out.Accepted = append(out.Accepted, r)
		out.Events = append(out.Events, AuditEvent{ID: "import-" + row.ExternalID, RecordID: r.ID, Action: "import", Actor: actor, Detail: "accepted", At: now})
	}
	return out
}
