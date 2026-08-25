package domain

import "sort"

type Change struct{ Field, Before, After string }

func Diff(a, b Record) []Change {
	out := []Change{}
	if a.Title != b.Title {
		out = append(out, Change{"title", a.Title, b.Title})
	}
	if a.Venue != b.Venue {
		out = append(out, Change{"venue", a.Venue, b.Venue})
	}
	if a.Region != b.Region {
		out = append(out, Change{"region", a.Region, b.Region})
	}
	if a.Score != b.Score {
		out = append(out, Change{"score", itoa(a.Score), itoa(b.Score)})
	}
	if a.Status != b.Status {
		out = append(out, Change{"status", a.Status, b.Status})
	}
	return out
}
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := make([]byte, 0, 12)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
func SortChanges(c []Change) { sort.Slice(c, func(i, j int) bool { return c[i].Field < c[j].Field }) }
func ApplyChanges(r *Record, changes []Change) {
	for _, c := range changes {
		switch c.Field {
		case "title":
			r.Title = c.After
		case "venue":
			r.Venue = c.After
		case "region":
			r.Region = c.After
		case "score":
			r.Score = parseInt(c.After)
		case "status":
			r.Status = c.After
		}
	}
}
func parseInt(s string) int {
	n := 0
	for _, c := range []byte(s) {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
