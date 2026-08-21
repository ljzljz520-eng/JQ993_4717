package store

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"aroma-maintenance/internal/domain"
	"go.etcd.io/bbolt"
)

var buckets = [][]byte{[]byte("records"), []byte("audits"), []byte("workflows"), []byte("attachments")}

type Store struct {
	db  *bbolt.DB
	mu  sync.RWMutex
	now func() time.Time
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, now: func() time.Time { return time.Unix(1700000000, 0) }}
	if err := s.init(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func OpenMemory(path string) (*Store, error) { return Open(path) }
func (s *Store) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return err
			}
		}
		return nil
	})
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) timestamp() string    { return s.now().UTC().Format(time.RFC3339) }
func encode(v any) ([]byte, error)    { return json.Marshal(v) }
func decode(data []byte, v any) error { return json.Unmarshal(data, v) }
func key(id string) []byte            { return []byte(id) }
func put(tx *bbolt.Tx, bucket []byte, id string, v any) error {
	data, err := encode(v)
	if err != nil {
		return err
	}
	return tx.Bucket(bucket).Put(key(id), data)
}
func get(tx *bbolt.Tx, bucket []byte, id string, v any) error {
	data := tx.Bucket(bucket).Get(key(id))
	if data == nil {
		return domain.ErrNotFound
	}
	return decode(data, v)
}
func (s *Store) SaveRecord(r domain.Record) error {
	if err := r.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error {
		if r.CreatedAt == "" {
			r.CreatedAt = s.timestamp()
		}
		r.UpdatedAt = s.timestamp()
		return put(tx, []byte("records"), r.ID, r)
	})
}
func (s *Store) CreateRecord(r domain.Record) error {
	if err := r.Validate(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error {
		if tx.Bucket([]byte("records")).Get(key(r.ID)) != nil {
			return domain.ErrConflict
		}
		r.CreatedAt = s.timestamp()
		r.UpdatedAt = r.CreatedAt
		return put(tx, []byte("records"), r.ID, r)
	})
}
func (s *Store) GetRecord(id string) (domain.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var r domain.Record
	err := s.db.View(func(tx *bbolt.Tx) error { return get(tx, []byte("records"), id, &r) })
	return r, err
}
func (s *Store) ListRecords() ([]domain.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Record{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, v []byte) error {
			var r domain.Record
			if err := decode(v, &r); err != nil {
				return err
			}
			out = append(out, r)
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, err
}
func (s *Store) DeleteRecord(id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte("records")).Delete(key(id)) })
}

func (s *Store) SaveAudit(e domain.AuditEvent) error {
	if e.ID == "" {
		return domain.ErrInvalid
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, []byte("audits"), e.ID, e) })
}
func (s *Store) ListAudits(recordID string) ([]domain.AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.AuditEvent{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("audits")).ForEach(func(_, v []byte) error {
			var e domain.AuditEvent
			if err := decode(v, &e); err != nil {
				return err
			}
			if recordID == "" || e.RecordID == recordID {
				out = append(out, e)
			}
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, err
}
func (s *Store) SaveWorkflow(w domain.Workflow) error {
	if w.ID == "" || w.RecordID == "" {
		return domain.ErrInvalid
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, []byte("workflows"), w.ID, w) })
}
func (s *Store) GetWorkflow(id string) (domain.Workflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var w domain.Workflow
	err := s.db.View(func(tx *bbolt.Tx) error { return get(tx, []byte("workflows"), id, &w) })
	return w, err
}
func (s *Store) ListWorkflows(recordID string) ([]domain.Workflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Workflow{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("workflows")).ForEach(func(_, v []byte) error {
			var w domain.Workflow
			if err := decode(v, &w); err != nil {
				return err
			}
			if recordID == "" || w.RecordID == recordID {
				out = append(out, w)
			}
			return nil
		})
	})
	return out, err
}
func (s *Store) SaveAttachment(a domain.Attachment) error {
	if err := domain.ValidateAttachment(a); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, []byte("attachments"), a.ID, a) })
}
func (s *Store) ListAttachments(recordID string) ([]domain.Attachment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Attachment{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("attachments")).ForEach(func(_, v []byte) error {
			var a domain.Attachment
			if err := decode(v, &a); err != nil {
				return err
			}
			if recordID == "" || a.RecordID == recordID {
				out = append(out, a)
			}
			return nil
		})
	})
	return out, err
}
func (s *Store) RequireOpen() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return nil
}
