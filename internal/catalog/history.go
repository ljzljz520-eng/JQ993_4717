package catalog

import (
	"aroma-maintenance/internal/domain"
	"sort"
	"strings"
)

type Revision struct {
	Version int
	Summary string
	Actor   string
	At      string
}
type History struct {
	RecordID  string
	Revisions []Revision
}

func BuildHistory(events []domain.AuditEvent) History {
	h := History{}
	for _, e := range events {
		if h.RecordID == "" {
			h.RecordID = e.RecordID
		}
		h.Revisions = append(h.Revisions, Revision{Summary: e.Action + ": " + e.Detail, Actor: e.Actor, At: e.At})
	}
	sort.SliceStable(h.Revisions, func(i, j int) bool { return h.Revisions[i].At < h.Revisions[j].At })
	for i := range h.Revisions {
		h.Revisions[i].Version = i + 1
	}
	return h
}
func (h History) Latest() Revision {
	if len(h.Revisions) == 0 {
		return Revision{}
	}
	return h.Revisions[len(h.Revisions)-1]
}
func (h History) Contains(action string) bool {
	for _, r := range h.Revisions {
		if strings.Contains(r.Summary, action) {
			return true
		}
	}
	return false
}
func (h History) Actors() []string {
	seen := map[string]bool{}
	out := []string{}
	for _, r := range h.Revisions {
		if r.Actor != "" && !seen[r.Actor] {
			seen[r.Actor] = true
			out = append(out, r.Actor)
		}
	}
	sort.Strings(out)
	return out
}
func (h History) Count() int { return len(h.Revisions) }
func (s *Service) History(id string) (History, error) {
	events, err := s.store.ListAudits(id)
	if err != nil {
		return History{}, err
	}
	return BuildHistory(events), nil
}
