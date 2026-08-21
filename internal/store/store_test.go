package store

import (
	"aroma-maintenance/internal/domain"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	r := domain.NewRecord("r1", "993-27", "Cedar", "wood", "wax", "ops")
	if err := s.CreateRecord(r); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveAudit(domain.AuditEvent{ID: "e1", RecordID: "r1", Action: "create"}); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveWorkflow(domain.Workflow{ID: "w1", RecordID: "r1", Name: "review"}); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveAttachment(domain.Attachment{ID: "a1", RecordID: "r1", Name: "label", Kind: "pdf"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if _, err := s.GetRecord("r1"); err != nil {
		t.Fatal(err)
	}
	if rows, _ := s.ListAudits("r1"); len(rows) != 1 {
		t.Fatal("audit missing")
	}
	if rows, _ := s.ListWorkflows("r1"); len(rows) != 1 {
		t.Fatal("workflow missing")
	}
	if rows, _ := s.ListAttachments("r1"); len(rows) != 1 {
		t.Fatal("attachment missing")
	}
}
func TestSnapshotRoundTrip(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "x.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	r := domain.NewRecord("r1", "b", "n", "s", "m", "o")
	if err := s.CreateRecord(r); err != nil {
		t.Fatal(err)
	}
	data, err := s.ExportJSON("r1")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.ImportJSON(data); err != nil {
		t.Fatal(err)
	}
}
