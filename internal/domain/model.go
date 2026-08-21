package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound   = errors.New("aroma record not found")
	ErrInvalid    = errors.New("invalid aroma record")
	ErrConflict   = errors.New("aroma record conflict")
	ErrTransition = errors.New("invalid workflow transition")
)

type Record struct {
	ID          string       `json:"id"`
	Batch       string       `json:"batch"`
	Name        string       `json:"name"`
	Scent       string       `json:"scent"`
	Material    string       `json:"material"`
	Owner       string       `json:"owner"`
	Status      Status       `json:"status"`
	Version     int          `json:"version"`
	Notes       []string     `json:"notes,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

type Status string

const (
	Draft         Status = "draft"
	PendingReview Status = "pending_review"
	Confirmed     Status = "confirmed"
	Withdrawn     Status = "withdrawn"
	Archived      Status = "archived"
)

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Detail   string `json:"detail"`
	At       string `json:"at"`
}

type Workflow struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	Stage     string `json:"stage"`
	Assignee  string `json:"assignee"`
	Completed bool   `json:"completed"`
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Digest   string `json:"digest"`
}

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.Batch) == "" {
		return fmt.Errorf("%w: id and batch required", ErrInvalid)
	}
	if strings.TrimSpace(r.Name) == "" || strings.TrimSpace(r.Scent) == "" {
		return fmt.Errorf("%w: name and scent required", ErrInvalid)
	}
	if r.Version < 1 {
		return fmt.Errorf("%w: version must be positive", ErrInvalid)
	}
	return nil
}

func (r Record) IsTerminal() bool { return r.Status == Archived }
func (r Record) CanEdit() bool    { return r.Status == Draft || r.Status == Withdrawn }
func (r Record) HasTag(tag string) bool {
	for _, t := range r.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
func (r Record) DisplayName() string { return r.Batch + " / " + r.Name }
func (r Record) Clone() Record {
	c := r
	c.Notes = append([]string(nil), r.Notes...)
	c.Tags = append([]string(nil), r.Tags...)
	c.Attachments = append([]Attachment(nil), r.Attachments...)
	return c
}

func NewRecord(id, batch, name, scent, material, owner string) Record {
	return Record{ID: id, Batch: batch, Name: name, Scent: scent, Material: material, Owner: owner, Status: Draft, Version: 1}
}

func TransitionAllowed(from, to Status) bool {
	switch from {
	case Draft:
		return to == PendingReview
	case PendingReview:
		return to == Confirmed || to == Withdrawn
	case Confirmed:
		return to == Withdrawn || to == Archived
	case Withdrawn:
		return to == PendingReview || to == Archived
	case Archived:
		return false
	default:
		return false
	}
}

func (r *Record) Transition(to Status) error {
	if r == nil || !TransitionAllowed(r.Status, to) {
		return ErrTransition
	}
	r.Status = to
	r.Version++
	return nil
}
