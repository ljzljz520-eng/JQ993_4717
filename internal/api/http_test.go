package api

import (
	"aroma-maintenance/internal/catalog"
	"aroma-maintenance/internal/importer"
	"aroma-maintenance/internal/report"
	"aroma-maintenance/internal/review"
	"aroma-maintenance/internal/search"
	"aroma-maintenance/internal/store"
	"bytes"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHTTPCreateList(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "x.db"))
	defer s.Close()
	h := NewHandler(catalog.New(s), review.New(s), search.New(s), importer.New(), report.New(s))
	req := httptest.NewRequest("POST", "/records", bytes.NewBufferString(`{"id":"r1","batch":"b","name":"n","scent":"s","material":"m","owner":"o"}`))
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	if res.Code != 201 {
		t.Fatalf("create %d", res.Code)
	}
	res = httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest("GET", "/records?q=n", nil))
	if res.Code != 200 {
		t.Fatalf("list %d", res.Code)
	}
}
