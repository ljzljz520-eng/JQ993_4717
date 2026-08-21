package domain

import "testing"

func TestRecordValidationAndTransitions(t *testing.T) {
	r := NewRecord("r1", "b1", "Cedar", "wood", "wax", "ops")
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := r.Transition(PendingReview); err != nil {
		t.Fatal(err)
	}
	if err := r.Transition(Archived); err == nil {
		t.Fatal("expected invalid transition")
	}
	if !MatchFilter(r, Filter{Query: "cedar", Status: PendingReview}) {
		t.Fatal("filter mismatch")
	}
}
func TestMetadataHelpers(t *testing.T) {
	r := NewRecord("r1", "b1", "Cedar", "wood", "wax", "ops")
	if !AddNote(&r, " note ") {
		t.Fatal("note")
	}
	if !AddTag(&r, "Safety") {
		t.Fatal("tag")
	}
	if !r.HasTag("safety") {
		t.Fatal("normalized tag")
	}
	if err := Attach(&r, Attachment{ID: "a1", Name: "label", Kind: "pdf"}); err != nil {
		t.Fatal(err)
	}
}
