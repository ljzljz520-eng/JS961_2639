package domain

type Metrics struct {
	Count, Approved, Archived int
	ScoreTotal                int
	ScoreMin, ScoreMax        int
}

func ComputeMetrics(records []Record) Metrics {
	m := Metrics{ScoreMin: 101}
	for _, r := range records {
		m.Count++
		m.ScoreTotal += r.Score
		if r.Score < m.ScoreMin {
			m.ScoreMin = r.Score
		}
		if r.Score > m.ScoreMax {
			m.ScoreMax = r.Score
		}
		if r.Status == StatusApproved {
			m.Approved++
		}
		if r.Status == StatusArchived {
			m.Archived++
		}
	}
	if m.Count == 0 {
		m.ScoreMin = 0
	}
	return m
}
func (m Metrics) Average() float64 {
	if m.Count == 0 {
		return 0
	}
	return float64(m.ScoreTotal) / float64(m.Count)
}
func (m Metrics) CompletionRate() float64 {
	if m.Count == 0 {
		return 0
	}
	return float64(m.Approved+m.Archived) / float64(m.Count)
}
