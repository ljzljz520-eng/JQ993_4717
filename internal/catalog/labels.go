package catalog

import (
	"aroma-maintenance/internal/domain"
	"sort"
)

func GroupByTag(rows []domain.Record) map[string][]domain.Record {
	out := map[string][]domain.Record{}
	for _, r := range rows {
		for _, t := range r.Tags {
			out[t] = append(out[t], r)
		}
	}
	for key := range out {
		sort.Slice(out[key], func(i, j int) bool { return out[key][i].ID < out[key][j].ID })
	}
	return out
}
func UniqueOwners(rows []domain.Record) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, r := range rows {
		if !seen[r.Owner] {
			seen[r.Owner] = true
			out = append(out, r.Owner)
		}
	}
	sort.Strings(out)
	return out
}
func CountByStatus(rows []domain.Record) map[domain.Status]int {
	out := map[domain.Status]int{}
	for _, r := range rows {
		out[r.Status]++
	}
	return out
}
