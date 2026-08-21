package store

import (
	"aroma-maintenance/internal/domain"
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
)

type Snapshot struct {
	Record      domain.Record
	Audits      []domain.AuditEvent
	Workflows   []domain.Workflow
	Attachments []domain.Attachment
}

func (s *Store) Snapshot(recordID string) (Snapshot, error) {
	r, err := s.GetRecord(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	a, err := s.ListAudits(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	w, err := s.ListWorkflows(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	at, err := s.ListAttachments(recordID)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Record: r, Audits: a, Workflows: w, Attachments: at}, nil
}
func (s *Store) SaveSnapshot(x Snapshot) error {
	if err := s.RequireOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		if err := put(tx, []byte("records"), x.Record.ID, x.Record); err != nil {
			return err
		}
		for _, e := range x.Audits {
			if err := put(tx, []byte("audits"), e.ID, e); err != nil {
				return err
			}
		}
		for _, w := range x.Workflows {
			if err := put(tx, []byte("workflows"), w.ID, w); err != nil {
				return err
			}
		}
		for _, a := range x.Attachments {
			if err := put(tx, []byte("attachments"), a.ID, a); err != nil {
				return err
			}
		}
		return nil
	})
}
func (s *Store) ExportJSON(recordID string) ([]byte, error) {
	x, err := s.Snapshot(recordID)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(x, "", "  ")
}
func (s *Store) ImportJSON(data []byte) error {
	var x Snapshot
	if err := json.Unmarshal(data, &x); err != nil {
		return fmt.Errorf("decode snapshot: %w", err)
	}
	return s.SaveSnapshot(x)
}
func withTx(db *bbolt.DB, writable bool, fn func(*bbolt.Tx) error) error {
	if writable {
		return db.Update(fn)
	}
	return db.View(fn)
}
