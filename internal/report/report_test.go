package report

import (
	"aroma-maintenance/internal/catalog"
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
	"path/filepath"
	"testing"
)

func TestSummaryAndFormats(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	c := catalog.New(s)
	r, _ := c.Register(domain.NewRecord("1", "b", "n", "s", "m", "o"))
	sum, err := New(s).Summary()
	if err != nil || sum.Total != 1 {
		t.Fatal(err)
	}
	if RenderRecord(r) == "" || RenderSummary(sum) == "" || Progress(sum) == "" {
		t.Fatal("format")
	}
	if _, err := New(s).CSV([]domain.Record{r}); err != nil {
		t.Fatal(err)
	}
}
