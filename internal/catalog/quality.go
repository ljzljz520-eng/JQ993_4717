package catalog

import (
	"aroma-maintenance/internal/domain"
	"strings"
)

type Quality struct {
	Score    int
	Warnings []string
	Ready    bool
}

func AssessQuality(r domain.Record) Quality {
	q := Quality{}
	q.Score = domain.RecordCompleteness(r) * 10
	if r.Status == domain.Archived {
		q.Score += 5
	}
	if len(r.Scent) < 3 {
		q.Warnings = append(q.Warnings, "scent description is short")
	}
	if len(r.Material) < 3 {
		q.Warnings = append(q.Warnings, "material description is short")
	}
	if len(r.Tags) == 0 {
		q.Warnings = append(q.Warnings, "no tags")
	}
	q.Ready = len(domain.MissingFields(r)) == 0 && len(q.Warnings) == 0
	return q
}
func NormalizeOwner(owner string) string { return strings.TrimSpace(strings.ToLower(owner)) }
func NormalizeBatch(batch string) string { return strings.TrimSpace(strings.ToUpper(batch)) }
func IsBatchFormat(batch string) bool {
	parts := strings.Split(NormalizeBatch(batch), "-")
	return len(parts) == 2 && len(parts[0]) > 0 && len(parts[1]) > 0
}
func (s *Service) Quality(id string) (Quality, error) {
	r, err := s.Get(id)
	if err != nil {
		return Quality{}, err
	}
	return AssessQuality(r), nil
}
