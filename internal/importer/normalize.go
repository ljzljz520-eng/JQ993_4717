package importer

import (
	"aroma-maintenance/internal/domain"
	"strings"
)

type Options struct {
	DefaultOwner    string
	DefaultMaterial string
	Strict          bool
}

func ApplyOptions(r domain.Record, o Options) domain.Record {
	if r.Owner == "" {
		r.Owner = o.DefaultOwner
	}
	if r.Material == "" {
		r.Material = o.DefaultMaterial
	}
	r.ID = strings.TrimSpace(r.ID)
	r.Batch = strings.TrimSpace(r.Batch)
	r.Name = strings.TrimSpace(r.Name)
	r.Scent = strings.TrimSpace(r.Scent)
	r.Owner = strings.TrimSpace(r.Owner)
	return r
}
func NormalizeRows(rows []domain.Record, o Options) ([]domain.Record, []error) {
	out := []domain.Record{}
	errs := []error{}
	for _, r := range rows {
		r = ApplyOptions(r, o)
		if err := r.Validate(); err != nil {
			errs = append(errs, err)
			if o.Strict {
				continue
			}
		}
		out = append(out, r)
	}
	return out, errs
}
func GroupByBatch(rows []domain.Record) map[string][]domain.Record {
	out := map[string][]domain.Record{}
	for _, r := range rows {
		out[r.Batch] = append(out[r.Batch], r)
	}
	return out
}
func BatchNames(rows []domain.Record) []string {
	groups := GroupByBatch(rows)
	out := make([]string, 0, len(groups))
	for batch := range groups {
		out = append(out, batch)
	}
	sortStrings(out)
	return out
}
func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}
func CSVHeaders() []string {
	return []string{"id", "batch", "name", "scent", "material", "owner", "tags"}
}
