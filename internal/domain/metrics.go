package domain

import "sort"

type StatusCount struct {
	Status Status
	Count  int
}
type Timeline struct {
	RecordID string
	First    string
	Last     string
	Actions  []string
}

func Statuses() []Status { return []Status{Draft, PendingReview, Confirmed, Withdrawn, Archived} }
func IsKnownStatus(s Status) bool {
	for _, known := range Statuses() {
		if known == s {
			return true
		}
	}
	return false
}
func StatusRank(s Status) int {
	for i, known := range Statuses() {
		if known == s {
			return i
		}
	}
	return -1
}
func SortStatuses(rows []Record) []Record {
	out := append([]Record(nil), rows...)
	sort.SliceStable(out, func(i, j int) bool { return StatusRank(out[i].Status) < StatusRank(out[j].Status) })
	return out
}
func CountStatuses(rows []Record) []StatusCount {
	counts := map[Status]int{}
	for _, r := range rows {
		counts[r.Status]++
	}
	out := make([]StatusCount, 0, len(counts))
	for _, s := range Statuses() {
		if counts[s] > 0 {
			out = append(out, StatusCount{s, counts[s]})
		}
	}
	return out
}
func RecordCompleteness(r Record) int {
	score := 0
	if r.ID != "" {
		score++
	}
	if r.Batch != "" {
		score++
	}
	if r.Name != "" {
		score++
	}
	if r.Scent != "" {
		score++
	}
	if r.Material != "" {
		score++
	}
	if r.Owner != "" {
		score++
	}
	if len(r.Tags) > 0 {
		score++
	}
	if len(r.Notes) > 0 {
		score++
	}
	if len(r.Attachments) > 0 {
		score++
	}
	return score
}
func MissingFields(r Record) []string {
	missing := []string{}
	if r.ID == "" {
		missing = append(missing, "id")
	}
	if r.Batch == "" {
		missing = append(missing, "batch")
	}
	if r.Name == "" {
		missing = append(missing, "name")
	}
	if r.Scent == "" {
		missing = append(missing, "scent")
	}
	if r.Material == "" {
		missing = append(missing, "material")
	}
	if r.Owner == "" {
		missing = append(missing, "owner")
	}
	return missing
}
func VersionAdvanced(old, new Record) bool { return new.Version > old.Version && new.ID == old.ID }
