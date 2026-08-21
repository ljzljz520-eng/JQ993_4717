package review

import (
	"aroma-maintenance/internal/domain"
	"sort"
	"strings"
)

type QueueItem struct {
	Record   domain.Record
	Priority int
	Reason   string
}

func BuildQueue(rows []domain.Record, policy Policy) []QueueItem {
	out := []QueueItem{}
	for _, r := range rows {
		if r.Status != domain.PendingReview && r.Status != domain.Withdrawn {
			continue
		}
		priority := 1
		reason := "pending review"
		if r.Status == domain.Withdrawn {
			priority = 3
			reason = "withdrawn retry"
		}
		if len(CheckPolicy(r, policy)) > 0 {
			priority += 2
			reason += "; metadata incomplete"
		}
		out = append(out, QueueItem{r, priority, reason})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].Record.Batch < out[j].Record.Batch
		}
		return out[i].Priority > out[j].Priority
	})
	return out
}
func QueueIDs(items []QueueItem) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.Record.ID)
	}
	return out
}
func FilterQueue(items []QueueItem, owner string) []QueueItem {
	out := []QueueItem{}
	for _, item := range items {
		if owner == "" || strings.EqualFold(item.Record.Owner, owner) {
			out = append(out, item)
		}
	}
	return out
}
func (s *Service) Queue(policy Policy) ([]QueueItem, error) {
	rows, err := s.store.ListRecords()
	if err != nil {
		return nil, err
	}
	return BuildQueue(rows, policy), nil
}
func (s *Service) WithdrawReason(id string) (string, error) {
	events, err := s.store.ListAudits(id)
	if err != nil {
		return "", err
	}
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Action == "withdraw" {
			return events[i].Detail, nil
		}
	}
	return "", domain.ErrNotFound
}
