package review

import (
	"aroma-maintenance/internal/catalog"
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
	"path/filepath"
	"testing"
)

func TestReviewLifecycle(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	c := catalog.New(s)
	r, _ := c.Register(domain.NewRecord("r1", "b", "n", "s", "m", "o"))
	svc := New(s)
	if _, err := svc.Submit(r.ID, "a"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Confirm(r.ID, "a"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Archive(r.ID, "a"); err != nil {
		t.Fatal(err)
	}
	events, _ := svc.History(r.ID)
	if len(events) != 3 {
		t.Fatal("history")
	}
}
func TestAssignmentPolicy(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	c := catalog.New(s)
	r, _ := c.Register(domain.NewRecord("r1", "b", "n", "s", "m", "o"))
	svc := New(s)
	w, err := svc.Assign(r.ID, "reviewer", "tomorrow")
	if err != nil || w.Assignee != "reviewer" {
		t.Fatal(err)
	}
	if _, err := svc.CompleteAssignment(r.ID); err != nil {
		t.Fatal(err)
	}
	if Eligible(r, DefaultPolicy()) {
		t.Fatal("missing policy tags")
	}
}
