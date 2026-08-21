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
	if err := s.store.SaveRecord(retry); err != nil {
		return domain.Record{}, err
	}
	return withdrawal.Record, nil
}
func (s *Service) RetryConfirmationFixed(id, actor, token string) (domain.Record, error) {
	r, err := s.store.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	w, err := domain.PrepareWithdrawal(r, "retry", token)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(w.Record); err != nil {
		return domain.Record{}, err
	}
	retry, err := domain.ConfirmRetry(w)
	if err != nil {
		return domain.Record{}, err
	}
	retry, err = s.audit(retry, "retry_confirm", actor, "retried after withdrawal")
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.store.SaveRecord(retry); err != nil {
		return domain.Record{}, err
	}
	return retry, nil
}
