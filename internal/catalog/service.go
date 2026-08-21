package catalog

import (
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
	"fmt"
	"strings"
)

type Service struct{ store *store.Store }

func New(s *store.Store) *Service { return &Service{store: s} }
func (s *Service) Register(input domain.Record) (domain.Record, error) {
	if input.Version < 1 {
		input.Version = 1
	}
	input.Status = domain.Draft
	input.Tags = domain.SanitizeTags(input.Tags)
	if err := input.Validate(); err != nil {
		return domain.Record{}, err
	}
	if err := s.store.CreateRecord(input); err != nil {
		return domain.Record{}, err
	}
	return s.store.GetRecord(input.ID)
}
func (s *Service) Update(id string, mutate func(*domain.Record) error) (domain.Record, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if !r.CanEdit() {
		return domain.Record{}, fmt.Errorf("%w: status %s", domain.ErrConflict, r.Status)
	}
	if err := mutate(&r); err != nil {
		return domain.Record{}, err
	}
	r.Tags = domain.SanitizeTags(r.Tags)
	r.Version++
	if err := s.store.SaveRecord(r); err != nil {
		return domain.Record{}, err
	}
	return r, nil
}
func (s *Service) AddMetadata(id, note, tag string) (domain.Record, error) {
	return s.Update(id, func(r *domain.Record) error {
		if note != "" && !domain.AddNote(r, note) {
			return domain.ErrInvalid
		}
		if tag != "" {
			domain.AddTag(r, tag)
		}
		return nil
	})
}
func (s *Service) AddAttachment(id string, a domain.Attachment) (domain.Record, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if !r.CanEdit() {
		return domain.Record{}, domain.ErrConflict
	}
	if err := domain.Attach(&r, a); err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveAttachment(a); err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(r); err != nil {
		return domain.Record{}, err
	}
	return r, nil
}
func (s *Service) Get(id string) (domain.Record, error) { return s.store.GetRecord(id) }
func (s *Service) List() ([]domain.Record, error)       { return s.store.ListRecords() }
func (s *Service) Remove(id string) error {
	r, err := s.Get(id)
	if err != nil {
		return err
	}
	if !r.IsTerminal() {
		return domain.ErrConflict
	}
	return s.store.DeleteRecord(id)
}
func NormalizeRecord(r domain.Record) domain.Record {
	r.ID = strings.TrimSpace(r.ID)
	r.Batch = strings.TrimSpace(r.Batch)
	r.Name = strings.TrimSpace(r.Name)
	r.Scent = strings.TrimSpace(r.Scent)
	r.Material = strings.TrimSpace(r.Material)
	r.Owner = strings.TrimSpace(r.Owner)
	r.Tags = domain.SanitizeTags(r.Tags)
	return r
}
