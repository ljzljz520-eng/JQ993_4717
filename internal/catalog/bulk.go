package catalog

import (
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
)

type BatchResult struct {
	Created []domain.Record
	Failed  []string
}

func (s *Service) RegisterMany(records []domain.Record) BatchResult {
	out := BatchResult{Created: []domain.Record{}, Failed: []string{}}
	for _, r := range records {
		r = NormalizeRecord(r)
		saved, err := s.Register(r)
		if err != nil {
			out.Failed = append(out.Failed, r.ID)
			continue
		}
		out.Created = append(out.Created, saved)
	}
	return out
}
func (s *Service) Summary() (domain.Summary, error) {
	rows, err := s.List()
	if err != nil {
		return domain.Summary{}, err
	}
	sum := domain.Summary{}
	for _, r := range rows {
		sum.Add(r)
	}
	return sum, nil
}
func (s *Service) Store() *store.Store { return s.store }
