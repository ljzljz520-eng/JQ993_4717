package main

import (
	"aroma-maintenance/internal/catalog"
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/review"
	"aroma-maintenance/internal/search"
	"aroma-maintenance/internal/store"
	"path/filepath"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	c := catalog.New(s)
	r, _ := c.Register(domain.NewRecord("r1", "993-01", "Cedar", "wood", "wax", "ops"))
	v := review.New(s)
	if _, err := v.Submit(r.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	if _, err := v.Confirm(r.ID, "reviewer"); err != nil {
		t.Fatal(err)
	}
	archived, err := v.Archive(r.ID, "reviewer")
	if err != nil || archived.Status != domain.Archived {
		t.Fatal(err)
	}
}
func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	c := catalog.New(s)
	r, _ := c.Register(domain.NewRecord("r1", "993-02", "Rose", "floral", "wax", "ops"))
	if _, err := c.AddMetadata(r.ID, "ready", "label"); err != nil {
		t.Fatal(err)
	}
	rows, err := search.New(s).Search(domain.Filter{Query: "rose"})
	if err != nil || len(rows) != 1 {
		t.Fatal(err)
	}
}
func TestWorkflowImportReport(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	c := catalog.New(s)
	out := c.RegisterMany([]domain.Record{domain.NewRecord("r1", "993-03", "Mint", "fresh", "wax", "ops")})
	if len(out.Created) != 1 {
		t.Fatal("import")
	}
	sum, err := c.Summary()
	if err != nil || sum.Total != 1 {
		t.Fatal(err)
	}
}
func TestBusiness27Regression(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	c := catalog.New(s)
	r, _ := c.Register(domain.NewRecord("r27", "993-27", "Cedar", "wood", "wax", "duty"))
	v := review.New(s)
	if _, err := v.Submit(r.ID, "duty"); err != nil {
		t.Fatal(err)
	}
	if _, err := v.Confirm(r.ID, "duty"); err != nil {
		t.Fatal(err)
	}
	got, err := v.RetryConfirmation(r.ID, "duty", "retry-27")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.PendingReview {
		t.Fatalf("expected confirmed retry state, got %s", got.Status)
	}
}
