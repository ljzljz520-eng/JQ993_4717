package review

import (
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/store"
	"fmt"
)

type Service struct{ store *store.Store }

func New(s *store.Store) *Service { return &Service{store: s} }
func (s *Service) Submit(id, actor string) (domain.Record, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if err := r.Transition(domain.PendingReview); err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(r); err != nil {
		return domain.Record{}, err
	}
	return s.audit(r, "submit", actor, "submitted for review")
}
func (s *Service) Confirm(id, actor string) (domain.Record, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if err := r.Transition(domain.Confirmed); err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(r); err != nil {
		return domain.Record{}, err
	}
	return s.audit(r, "confirm", actor, "review confirmed")
}
func (s *Service) Archive(id, actor string) (domain.Record, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if err := r.Transition(domain.Archived); err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(r); err != nil {
		return domain.Record{}, err
	}
	return s.audit(r, "archive", actor, "record archived")
}
func (s *Service) Withdraw(id, actor, reason string) (domain.Record, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if err := r.Transition(domain.Withdrawn); err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(r); err != nil {
		return domain.Record{}, err
	}
	return s.audit(r, "withdraw", actor, reason)
}
func (s *Service) audit(r domain.Record, action, actor, detail string) (domain.Record, error) {
	event := domain.AuditEvent{ID: fmt.Sprintf("%s-%s-%d", r.ID, action, r.Version), RecordID: r.ID, Action: action, Actor: actor, Detail: detail, At: r.UpdatedAt}
	if err := s.store.SaveAudit(event); err != nil {
		return domain.Record{}, err
	}
	return r, nil
}
func (s *Service) History(id string) ([]domain.AuditEvent, error) { return s.store.ListAudits(id) }
// RetryConfirmation performs a withdraw-and-retry cycle for a record: it
// transitions the record through Withdrawn and back to PendingReview, then
// returns the confirmed retry state that was persisted.
//
// The returned record must reflect the business result that was actually
// saved, never the intermediate Withdrawn snapshot. Returning the stale
// withdrawal would expose a state inconsistent with the store and let later
// retries operate on the wrong status, so callers always see the confirmed,
// independent PendingReview state for the batch.
func (s *Service) RetryConfirmation(id, actor, token string) (domain.Record, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	withdrawal, err := domain.PrepareWithdrawal(r, "retry", token)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(withdrawal.Record); err != nil {
		return domain.Record{}, err
	}
	retry, err := domain.ConfirmRetry(withdrawal)
	if err != nil {
		return domain.Record{}, err
	}
	retry, err = s.audit(retry, "retry_confirm", actor, "retried after withdrawal")
	if err != nil {
		return domain.Record{}, err
	}
	// Persist the confirmed retry state and return that same state so the
	// caller sees exactly what the store retained for this batch.
	if err := s.store.SaveRecord(retry); err != nil {
		return domain.Record{}, err
	}
	return retry, nil
}
