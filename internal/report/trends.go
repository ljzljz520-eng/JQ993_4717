package report

import (
	"aroma-maintenance/internal/domain"
	"sort"
)

type Trend struct {
	Batch     string
	Total     int
	Confirmed int
	Archived  int
	Attention int
}

func BuildTrends(rows []domain.Record) []Trend {
	groups := map[string]*Trend{}
	for _, r := range rows {
		t := groups[r.Batch]
		if t == nil {
			t = &Trend{Batch: r.Batch}
			groups[r.Batch] = t
		}
		t.Total++
		switch r.Status {
		case domain.Confirmed:
			t.Confirmed++
		case domain.Archived:
			t.Archived++
		case domain.Draft, domain.PendingReview, domain.Withdrawn:
			t.Attention++
		}
	}
	out := make([]Trend, 0, len(groups))
	for _, t := range groups {
		out = append(out, *t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Batch < out[j].Batch })
	return out
}
func TrendRate(t Trend) float64 {
	if t.Total == 0 {
		return 0
	}
	return float64(t.Confirmed+t.Archived) / float64(t.Total)
}
func NeedsEscalation(t Trend) bool { return t.Attention > t.Confirmed+t.Archived }
func HighestAttention(rows []Trend) Trend {
	if len(rows) == 0 {
		return Trend{}
	}
	out := append([]Trend(nil), rows...)
	sort.Slice(out, func(i, j int) bool { return out[i].Attention > out[j].Attention })
	return out[0]
}
func TrendLabels(rows []Trend) []string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.Batch)
	}
	return out
}
