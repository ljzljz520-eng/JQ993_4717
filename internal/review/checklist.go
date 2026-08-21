package review

import (
	"aroma-maintenance/internal/domain"
	"strings"
)

type Check struct {
	ID     string
	Label  string
	Passed bool
	Detail string
}

func Checklist(r domain.Record, p Policy) []Check {
	checks := []Check{}
	checks = append(checks, Check{"identity", "batch and id", r.ID != "" && r.Batch != "", ""})
	checks = append(checks, Check{"scent", "scent profile", r.Scent != "", ""})
	checks = append(checks, Check{"material", "material declaration", r.Material != "", ""})
	for _, tag := range p.RequiredTags {
		checks = append(checks, Check{"tag-" + tag, "required tag " + tag, r.HasTag(tag), ""})
	}
	for i := range checks {
		if !checks[i].Passed {
			checks[i].Detail = "missing requirement"
		}
	}
	return checks
}
func PassedChecks(checks []Check) int {
	count := 0
	for _, c := range checks {
		if c.Passed {
			count++
		}
	}
	return count
}
func ChecklistReady(checks []Check) bool {
	return len(checks) > 0 && PassedChecks(checks) == len(checks)
}
func CheckLabels(checks []Check) string {
	parts := []string{}
	for _, c := range checks {
		state := "fail"
		if c.Passed {
			state = "pass"
		}
		parts = append(parts, c.Label+"="+state)
	}
	return strings.Join(parts, ";")
}
func (s *Service) Checklist(id string, p Policy) ([]Check, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return nil, err
	}
	return Checklist(r, p), nil
}
