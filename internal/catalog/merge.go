package catalog

import (
	"aroma-maintenance/internal/domain"
	"sort"
)

func MergeRecords(primary, secondary domain.Record) domain.Record {
	out := primary.Clone()
	if out.Name == "" {
		out.Name = secondary.Name
	}
	if out.Scent == "" {
		out.Scent = secondary.Scent
	}
	if out.Material == "" {
		out.Material = secondary.Material
	}
	if out.Owner == "" {
		out.Owner = secondary.Owner
	}
	out.Tags = domain.SanitizeTags(append(out.Tags, secondary.Tags...))
	out.Notes = append(out.Notes, secondary.Notes...)
	out.Attachments = append(out.Attachments, secondary.Attachments...)
	return out
}
func Deduplicate(rows []domain.Record) []domain.Record {
	byID := map[string]domain.Record{}
	for _, r := range rows {
		if old, ok := byID[r.ID]; ok {
			byID[r.ID] = MergeRecords(old, r)
		} else {
			byID[r.ID] = r
		}
	}
	out := make([]domain.Record, 0, len(byID))
	for _, r := range byID {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func SelectNewest(rows []domain.Record) domain.Record {
	out := Deduplicate(rows)
	if len(out) == 0 {
		return domain.Record{}
	}
	return out[len(out)-1]
}
