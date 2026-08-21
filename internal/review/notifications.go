package review

import (
	"aroma-maintenance/internal/domain"
	"fmt"
	"sort"
)

type Notification struct {
	Recipient string
	Subject   string
	Body      string
	Priority  int
}

func BuildNotification(r domain.Record, action, recipient string) Notification {
	priority := 1
	if action == "withdraw" {
		priority = 3
	}
	return Notification{Recipient: recipient, Subject: "aroma record " + action + ": " + r.DisplayName(), Body: "Record " + r.ID + " is now " + string(r.Status), Priority: priority}
}
func SortNotifications(rows []Notification) []Notification {
	out := append([]Notification(nil), rows...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].Recipient < out[j].Recipient
		}
		return out[i].Priority > out[j].Priority
	})
	return out
}
func NotificationBatch(rows []domain.Record, action, recipient string) []Notification {
	out := make([]Notification, 0, len(rows))
	for _, r := range rows {
		out = append(out, BuildNotification(r, action, recipient))
	}
	return SortNotifications(out)
}
func NotificationText(n Notification) string {
	return fmt.Sprintf("[%d] %s -> %s", n.Priority, n.Recipient, n.Subject)
}
func (s *Service) Notify(id, action, recipient string) (Notification, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return Notification{}, err
	}
	return BuildNotification(r, action, recipient), nil
}
