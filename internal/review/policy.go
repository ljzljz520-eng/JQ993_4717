package review

import "aroma-maintenance/internal/domain"

type Policy struct {
	RequiredTags    []string
	RequireMaterial bool
}

func DefaultPolicy() Policy {
	return Policy{RequiredTags: []string{"safety", "label"}, RequireMaterial: true}
}
func CheckPolicy(r domain.Record, p Policy) []string {
	missing := []string{}
	for _, tag := range p.RequiredTags {
		if !r.HasTag(tag) {
			missing = append(missing, tag)
		}
	}
	if p.RequireMaterial && r.Material == "" {
		missing = append(missing, "material")
	}
	return missing
}
func Eligible(r domain.Record, p Policy) bool { return len(CheckPolicy(r, p)) == 0 }
func NextAction(r domain.Record, p Policy) string {
	if r.Status == domain.Draft && !Eligible(r, p) {
		return "complete metadata"
	}
	if r.Status == domain.Draft {
		return "submit"
	}
	if r.Status == domain.PendingReview {
		return "confirm"
	}
	if r.Status == domain.Confirmed {
		return "archive"
	}
	return "inspect"
}
