package search

import (
	"aroma-maintenance/internal/domain"
	"sort"
)

type Facet struct {
	Value string
	Count int
}

func OwnerFacets(rows []domain.Record) []Facet {
	counts := map[string]int{}
	for _, r := range rows {
		counts[r.Owner]++
	}
	return facetList(counts)
}
func TagFacets(rows []domain.Record) []Facet {
	counts := map[string]int{}
	for _, r := range rows {
		for _, tag := range r.Tags {
			counts[tag]++
		}
	}
	return facetList(counts)
}
func StatusFacets(rows []domain.Record) []Facet {
	counts := map[string]int{}
	for _, r := range rows {
		counts[string(r.Status)]++
	}
	return facetList(counts)
}
func facetList(counts map[string]int) []Facet {
	out := make([]Facet, 0, len(counts))
	for value, count := range counts {
		out = append(out, Facet{value, count})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count == out[j].Count {
			return out[i].Value < out[j].Value
		}
		return out[i].Count > out[j].Count
	})
	return out
}
func FacetValues(facets []Facet) []string {
	out := make([]string, 0, len(facets))
	for _, facet := range facets {
		out = append(out, facet.Value)
	}
	return out
}
func LimitFacets(facets []Facet, limit int) []Facet {
	if limit < 1 {
		return []Facet{}
	}
	if len(facets) > limit {
		return facets[:limit]
	}
	return facets
}
