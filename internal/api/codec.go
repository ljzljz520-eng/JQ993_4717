package api

import (
	"aroma-maintenance/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ErrorBody struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
type RecordRequest struct {
	ID       string   `json:"id"`
	Batch    string   `json:"batch"`
	Name     string   `json:"name"`
	Scent    string   `json:"scent"`
	Material string   `json:"material"`
	Owner    string   `json:"owner"`
	Tags     []string `json:"tags"`
}

func DecodeRecord(r io.Reader) (domain.Record, error) {
	var input RecordRequest
	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return domain.Record{}, err
	}
	out := domain.NewRecord(input.ID, input.Batch, input.Name, input.Scent, input.Material, input.Owner)
	out.Tags = input.Tags
	return out, nil
}
func EncodeError(err error) ErrorBody {
	code := "invalid_request"
	if err == domain.ErrNotFound {
		code = "not_found"
	}
	if err == domain.ErrConflict {
		code = "conflict"
	}
	return ErrorBody{Error: err.Error(), Code: code}
}
func WriteError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(EncodeError(err))
}
func ReadActor(r *http.Request) string {
	actor := strings.TrimSpace(r.Header.Get("x-actor"))
	if actor == "" {
		return "operator"
	}
	return actor
}
func ParsePath(path string) (string, string, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[0] != "records" {
		return "", "", false
	}
	return parts[1], parts[2], true
}
func RequestID(r *http.Request) string { return strings.TrimSpace(r.Header.Get("x-request-id")) }
func RequireJSON(r *http.Request) error {
	if !strings.Contains(r.Header.Get("content-type"), "application/json") {
		return fmt.Errorf("content type must be application/json")
	}
	return nil
}
