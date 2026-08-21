package report

import (
	"aroma-maintenance/internal/domain"
	"sort"
)

type Dashboard struct {
	Summary    domain.Summary
	Queue      int
	Owners     map[string]int
	Tags       map[string]int
	Completion string
}

func BuildDashboard(rows []domain.Record) Dashboard {
	sum := domain.Summary{}
	owners := map[string]int{}
	tags := map[string]int{}
	queue := 0
	for _, r := range rows {
		sum.Add(r)
		owners[r.Owner]++
		for _, tag := range r.Tags {
			tags[tag]++
		}
		if r.Status == domain.PendingReview || r.Status == domain.Withdrawn {
			queue++
		}
	}
	return Dashboard{Summary: sum, Queue: queue, Owners: owners, Tags: tags, Completion: Progress(sum)}
}
func TopOwners(values map[string]int, limit int) []string {
	type pair struct {
		name  string
		count int
	}
	pairs := []pair{}
	for name, count := range values {
		pairs = append(pairs, pair{name, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].name < pairs[j].name
		}
		return pairs[i].count > pairs[j].count
	})
	if limit < 1 || limit > len(pairs) {
		limit = len(pairs)
	}
	out := []string{}
	for _, pair := range pairs[:limit] {
		out = append(out, pair.name)
	}
	return out
}
func RiskLevel(d Dashboard) string {
	if d.Queue > 5 {
		return "high"
	}
	if d.Queue > 0 {
		return "medium"
	}
	return "low"
}
func DashboardRows(d Dashboard) []string {
	return []string{"completion=" + d.Completion, "queue=" + itoa(d.Queue), "risk=" + RiskLevel(d), "owners=" + itoa(len(d.Owners)), "tags=" + itoa(len(d.Tags))}
}
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	out := ""
	for n > 0 {
		out = string(byte('0'+n%10)) + out
		n /= 10
	}
	return out
}
