package search

import (
	"aroma-maintenance/internal/domain"
	"sort"
	"strings"
)

type Hit struct {
	Record  domain.Record
	Score   int
	Reasons []string
}

func Rank(r domain.Record, q string) Hit {
	query := strings.ToLower(strings.TrimSpace(q))
	h := Hit{Record: r}
	if query == "" {
		return h
	}
	fields := []struct{ name, value string }{{"name", r.Name}, {"batch", r.Batch}, {"scent", r.Scent}, {"material", r.Material}, {"owner", r.Owner}}
	for _, field := range fields {
		value := strings.ToLower(field.value)
		if value == query {
			h.Score += 10
			h.Reasons = append(h.Reasons, field.name+" exact")
		} else if strings.Contains(value, query) {
			h.Score += 4
			h.Reasons = append(h.Reasons, field.name+" contains")
		}
	}
	for _, tag := range r.Tags {
		if strings.EqualFold(tag, query) {
			h.Score += 6
			h.Reasons = append(h.Reasons, "tag match")
		}
	}
	return h
}
func RankRows(rows []domain.Record, q string) []Hit {
	out := make([]Hit, 0, len(rows))
	for _, r := range rows {
		h := Rank(r, q)
		if q == "" || h.Score > 0 {
			out = append(out, h)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Record.ID < out[j].Record.ID
		}
		return out[i].Score > out[j].Score
	})
	return out
}
func Highlight(value, query string) string {
	if query == "" {
		return value
	}
	lower := strings.ToLower(value)
	pos := strings.Index(lower, strings.ToLower(query))
	if pos < 0 {
		return value
	}
	return value[:pos] + "[" + value[pos:pos+len(query)] + "]" + value[pos+len(query):]
}
func ExtractTerms(query string) []string {
	parts := strings.Fields(strings.ToLower(query))
	seen := map[string]bool{}
	out := []string{}
	for _, part := range parts {
		if !seen[part] {
			seen[part] = true
			out = append(out, part)
		}
	}
	return out
}
