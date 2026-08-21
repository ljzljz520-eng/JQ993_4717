package domain

import "strings"

type Filter struct {
	Query  string
	Status Status
	Tag    string
	Owner  string
}

func NormalizeFilter(f Filter) Filter {
	f.Query = strings.ToLower(strings.TrimSpace(f.Query))
	f.Tag = strings.ToLower(strings.TrimSpace(f.Tag))
	f.Owner = strings.TrimSpace(f.Owner)
	return f
}
func MatchFilter(r Record, f Filter) bool {
	f = NormalizeFilter(f)
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Owner != "" && r.Owner != f.Owner {
		return false
	}
	if f.Tag != "" && !r.HasTag(f.Tag) {
		return false
	}
	if f.Query == "" {
		return true
	}
	haystack := strings.ToLower(r.Name + " " + r.Batch + " " + r.Scent + " " + r.Material + " " + strings.Join(r.Notes, " "))
	return strings.Contains(haystack, f.Query)
}

func SanitizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	seen := map[string]bool{}
	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag != "" && !seen[tag] {
			out = append(out, tag)
			seen[tag] = true
		}
	}
	return out
}
func AddNote(r *Record, note string) bool {
	if r == nil || strings.TrimSpace(note) == "" {
		return false
	}
	r.Notes = append(r.Notes, strings.TrimSpace(note))
	return true
}
func AddTag(r *Record, tag string) bool {
	if r == nil {
		return false
	}
	normalized := SanitizeTags(append(r.Tags, tag))
	if len(normalized) == len(r.Tags) {
		return false
	}
	r.Tags = normalized
	return true
}
func Attach(r *Record, a Attachment) error {
	if r == nil || strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.Name) == "" {
		return ErrInvalid
	}
	a.RecordID = r.ID
	r.Attachments = append(r.Attachments, a)
	return nil
}
func ValidateAttachment(a Attachment) error {
	if strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.Name) == "" || strings.TrimSpace(a.Kind) == "" {
		return ErrInvalid
	}
	return nil
}
