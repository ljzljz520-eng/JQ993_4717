package importer

import "aroma-maintenance/internal/domain"

type Stats struct {
	Rows       int
	Valid      int
	Invalid    int
	Duplicates int
	Tags       int
}

func Measure(rows []domain.Record, errs []error) Stats {
	stats := Stats{Rows: len(rows), Invalid: len(errs)}
	seen := map[string]bool{}
	for _, r := range rows {
		if r.Validate() == nil {
			stats.Valid++
		}
		if seen[r.ID] {
			stats.Duplicates++
		}
		seen[r.ID] = true
		stats.Tags += len(r.Tags)
	}
	return stats
}
func (s Stats) AcceptanceRate() float64 {
	if s.Rows == 0 {
		return 0
	}
	return float64(s.Valid) / float64(s.Rows)
}
func (s Stats) HasErrors() bool { return s.Invalid > 0 || s.Duplicates > 0 }
func (s Stats) Status() string {
	if s.Rows == 0 {
		return "empty"
	}
	if s.HasErrors() {
		return "review"
	}
	return "ready"
}
func RequiredHeaders(headers []string) []string {
	required := []string{"id", "batch", "name", "scent"}
	missing := []string{}
	for _, need := range required {
		found := false
		for _, got := range headers {
			if got == need {
				found = true
			}
		}
		if !found {
			missing = append(missing, need)
		}
	}
	return missing
}
func IsCompatible(r domain.Record) bool {
	return r.ID != "" && r.Batch != "" && r.Name != "" && r.Scent != ""
}
