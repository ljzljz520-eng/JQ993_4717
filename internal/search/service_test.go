package search

import (
	"aroma-maintenance/internal/catalog"
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
	"path/filepath"
	"testing"
)

func TestSearchSelectPaginate(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	c := catalog.New(s)
	_, _ = c.Register(domain.NewRecord("1", "b1", "Cedar", "wood", "wax", "o"))
	_, _ = c.Register(domain.NewRecord("2", "b2", "Rose", "floral", "wax", "o"))
	svc := New(s)
	rows, err := svc.Search(domain.Filter{Query: "cedar"})
	if err != nil || len(rows) != 1 {
		t.Fatal(err)
	}
	if _, ok := svc.Select(rows, 0); !ok {
		t.Fatal("select")
	}
	if len(Paginate(rows, 2, 10)) != 0 {
		t.Fatal("page")
	}
}
