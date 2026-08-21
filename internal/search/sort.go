package search

import (
	"aroma-maintenance/internal/domain"
	"sort"
)

func SortByBatch(rows []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rows...)
	sort.Slice(out, func(i, j int) bool { return out[i].Batch < out[j].Batch })
	return out
}
func SortByStatus(rows []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rows...)
	sort.Slice(out, func(i, j int) bool { return out[i].Status < out[j].Status })
	return out
}
func Paginate(rows []domain.Record, page, size int) []domain.Record {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	start := (page - 1) * size
	if start >= len(rows) {
		return []domain.Record{}
	}
	end := start + size
	if end > len(rows) {
		end = len(rows)
	}
	return rows[start:end]
}
