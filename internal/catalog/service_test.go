package catalog

import (
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
	"path/filepath"
	"testing"
)

func TestRegisterUpdateAndAttach(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	svc := New(s)
	r, err := svc.Register(domain.NewRecord("r1", "b", "name", "scent", "wax", "owner"))
	if err != nil {
		t.Fatal(err)
	}
	r, err = svc.AddMetadata(r.ID, "note", "tag")
	if err != nil || len(r.Notes) != 1 {
		t.Fatal(err)
	}
	if _, err = svc.AddAttachment(r.ID, domain.Attachment{ID: "a", Name: "label", Kind: "image"}); err != nil {
		t.Fatal(err)
	}
}
func TestRegisterManyAndSummary(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	svc := New(s)
	out := svc.RegisterMany([]domain.Record{domain.NewRecord("1", "b1", "n", "s", "m", "o"), domain.NewRecord("", "b2", "n", "s", "m", "o")})
	if len(out.Created) != 1 || len(out.Failed) != 1 {
		t.Fatal("batch result")
	}
	sum, err := svc.Summary()
	if err != nil || sum.Total != 1 {
		t.Fatal(err)
	}
}
