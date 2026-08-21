package report

import (
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Service struct{ store *store.Store }

func New(s *store.Store) *Service { return &Service{store: s} }
func (s *Service) Summary() (domain.Summary, error) {
	rows, err := s.store.ListRecords()
	if err != nil {
		return domain.Summary{}, err
	}
	sum := domain.Summary{}
	for _, r := range rows {
		sum.Add(r)
	}
	return sum, nil
}
func (s *Service) JSONSummary() ([]byte, error) {
	sum, err := s.Summary()
	if err != nil {
		return nil, err
	}
	return json.Marshal(sum)
}
func (s *Service) CSV(records []domain.Record) (string, error) {
	var b strings.Builder
	w := csv.NewWriter(&b)
	if err := w.Write([]string{"id", "batch", "name", "scent", "status", "owner", "version"}); err != nil {
		return "", err
	}
	for _, r := range records {
		if err := w.Write([]string{r.ID, r.Batch, r.Name, r.Scent, string(r.Status), r.Owner, fmt.Sprint(r.Version)}); err != nil {
			return "", err
		}
	}
	w.Flush()
	return b.String(), w.Error()
}
func (s *Service) AuditReport(id string) ([]domain.AuditEvent, error) {
	events, err := s.store.ListAudits(id)
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].At < events[j].At })
	return events, nil
}
func (s *Service) WorkQueue() ([]domain.Record, error) {
	rows, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	out := []domain.Record{}
	for _, r := range rows {
		if r.Status == domain.PendingReview || r.Status == domain.Withdrawn {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Batch < out[j].Batch })
	return out, nil
}
func (s *Service) GroupOwners() map[string]int {
	rows, _ := s.store.ListRecords()
	out := map[string]int{}
	for _, r := range rows {
		out[r.Owner]++
	}
	return out
}
