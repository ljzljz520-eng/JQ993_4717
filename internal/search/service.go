package search

import (
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
	"sort"
)

type Service struct{ store *store.Store }

func New(s *store.Store) *Service { return &Service{store: s} }
func (s *Service) Search(f domain.Filter) ([]domain.Record, error) {
	rows, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	out := []domain.Record{}
	for _, r := range rows {
		if domain.MatchFilter(r, f) {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt > out[j].UpdatedAt })
	return out, nil
}
func (s *Service) FindExact(batch string) ([]domain.Record, error) {
	return s.Search(domain.Filter{Query: batch})
}
func (s *Service) Select(rows []domain.Record, index int) (domain.Record, bool) {
	if index < 0 || index >= len(rows) {
		return domain.Record{}, false
	}
	return rows[index], true
}
func (s *Service) Suggestions(query string) ([]string, error) {
	rows, err := s.Search(domain.Filter{Query: query})
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	out := []string{}
	for _, r := range rows {
		if !seen[r.Name] {
			seen[r.Name] = true
			out = append(out, r.Name)
		}
	}
	return out, nil
}
