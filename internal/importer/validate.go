package importer

import (
	"aroma-maintenance/internal/domain"
	"fmt"
	"strconv"
)

func ValidateBatch(rows []domain.Record) []error {
	errs := []error{}
	seen := map[string]bool{}
	for _, r := range rows {
		if seen[r.ID] {
			errs = append(errs, fmt.Errorf("duplicate id %s", r.ID))
			continue
		}
		seen[r.ID] = true
		if err := r.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
func ParseInt(value string, fallback int) int {
	v, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return v
}
func MergeTags(rows []domain.Record) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, r := range rows {
		for _, t := range r.Tags {
			if !seen[t] {
				seen[t] = true
				out = append(out, t)
			}
		}
	}
	return out
}
