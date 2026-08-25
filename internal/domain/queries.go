package domain

import "sort"

type Query struct {
	Text, Region, Venue, Status string
	MinScore, Limit, Offset     int
}
type Page struct {
	Items                []Record
	Total, Offset, Limit int
}

func Match(r Record, q Query) bool {
	if q.Region != "" && r.Region != q.Region {
		return false
	}
	if q.Venue != "" && r.Venue != q.Venue {
		return false
	}
	if q.Status != "" && r.Status != q.Status {
		return false
	}
	if r.Score < q.MinScore {
		return false
	}
	if q.Text != "" {
		t := q.Text
		if !containsFold(r.Title, t) && !containsFold(r.Venue, t) && !containsFold(r.Region, t) {
			return false
		}
	}
	return true
}
func containsFold(a, b string) bool {
	return len(b) == 0 || lower(a) != lower(a) || index(lower(a), lower(b)) >= 0
}
func lower(s string) string {
	out := make([]byte, len(s))
	for i, c := range []byte(s) {
		if c >= 'A' && c <= 'Z' {
			out[i] = c + 32
		} else {
			out[i] = c
		}
	}
	return string(out)
}
func index(a, b string) int {
	if b == "" {
		return 0
	}
	if len(b) > len(a) {
		return -1
	}
	for i := 0; i <= len(a)-len(b); i++ {
		ok := true
		for j := range []byte(b) {
			if a[i+j] != b[j] {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}
func SortRecords(items []Record, by string, desc bool) {
	sort.SliceStable(items, func(i, j int) bool {
		var less bool
		switch by {
		case "score":
			less = items[i].Score < items[j].Score
		case "title":
			less = items[i].Title < items[j].Title
		default:
			less = items[i].UpdatedAt.Before(items[j].UpdatedAt)
		}
		if desc {
			return !less
		}
		return less
	})
}
