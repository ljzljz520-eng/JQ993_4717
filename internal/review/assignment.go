package review

import (
	"aroma-maintenance/internal/domain"
)

type Assignment struct {
	Reviewer string
	Due      string
}

func (s *Service) Assign(id, reviewer, due string) (domain.Workflow, error) {
	if reviewer == "" || due == "" {
		return domain.Workflow{}, domain.ErrInvalid
	}
	if _, err := s.store.GetRecord(id); err != nil {
		return domain.Workflow{}, err
	}
	w := domain.Workflow{ID: id + "-review", RecordID: id, Name: "quality-review", Stage: "assigned", Assignee: reviewer}
	if err := s.store.SaveWorkflow(w); err != nil {
		return domain.Workflow{}, err
	}
	return w, nil
}
func (s *Service) CompleteAssignment(id string) (domain.Workflow, error) {
	w, err := s.store.GetWorkflow(id + "-review")
	if err != nil {
		return domain.Workflow{}, err
	}
	if w.Completed {
		return w, nil
	}
	w.Stage = "completed"
	w.Completed = true
	if err := s.store.SaveWorkflow(w); err != nil {
		return domain.Workflow{}, err
	}
	return w, nil
}
func (s *Service) Assignments(recordID string) ([]domain.Workflow, error) {
	return s.store.ListWorkflows(recordID)
}
