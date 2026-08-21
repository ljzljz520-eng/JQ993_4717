package domain

import "strings"

type Change struct {
	Field  string
	Before string
	After  string
}

func Diff(a, b Record) []Change {
	out := []Change{}
	add := func(field, before, after string) {
		if before != after {
			out = append(out, Change{field, before, after})
		}
	}
	add("batch", a.Batch, b.Batch)
	add("name", a.Name, b.Name)
	add("scent", a.Scent, b.Scent)
	add("material", a.Material, b.Material)
	add("owner", a.Owner, b.Owner)
	add("status", string(a.Status), string(b.Status))
	add("version", itoa(a.Version), itoa(b.Version))
	add("notes", strings.Join(a.Notes, "|"), strings.Join(b.Notes, "|"))
	add("tags", strings.Join(a.Tags, "|"), strings.Join(b.Tags, "|"))
	return out
}
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
func HasChange(changes []Change, field string) bool {
	for _, c := range changes {
		if c.Field == field {
			return true
		}
	}
	return false
}
func ChangeSummary(changes []Change) string {
	parts := make([]string, 0, len(changes))
	for _, c := range changes {
		parts = append(parts, c.Field+":"+c.Before+">"+c.After)
	}
	return strings.Join(parts, ", ")
}
