package api

import (
	"aroma-maintenance/internal/catalog"
	"aroma-maintenance/internal/domain"
	"aroma-maintenance/internal/importer"
	"aroma-maintenance/internal/report"
	"aroma-maintenance/internal/review"
	"aroma-maintenance/internal/search"
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	catalog  *catalog.Service
	review   *review.Service
	search   *search.Service
	importer *importer.Parser
	report   *report.Service
}

func NewHandler(c *catalog.Service, r *review.Service, q *search.Service, i *importer.Parser, p *report.Service) *Handler {
	return &Handler{catalog: c, review: r, search: q, importer: i, report: p}
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	path := strings.Trim(req.URL.Path, "/")
	if req.Method == http.MethodGet && path == "records" {
		h.list(w, req)
		return
	}
	if req.Method == http.MethodPost && path == "records" {
		h.create(w, req)
		return
	}
	if req.Method == http.MethodPost && path == "import" {
		h.importCSV(w, req)
		return
	}
	if req.Method == http.MethodGet && path == "reports/summary" {
		sum, err := h.report.Summary()
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, sum)
		return
	}
	parts := strings.Split(path, "/")
	if len(parts) == 3 && parts[0] == "records" && req.Method == http.MethodPost {
		h.transition(w, req, parts[1], parts[2])
		return
	}
	http.NotFound(w, req)
}
func (h *Handler) list(w http.ResponseWriter, req *http.Request) {
	rows, err := h.search.Search(domain.Filter{Query: req.URL.Query().Get("q"), Owner: req.URL.Query().Get("owner"), Tag: req.URL.Query().Get("tag"), Status: domain.Status(req.URL.Query().Get("status"))})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rows)
}
func (h *Handler) create(w http.ResponseWriter, req *http.Request) {
	var r domain.Record
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		writeErr(w, err)
		return
	}
	saved, err := h.catalog.Register(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}
func (h *Handler) transition(w http.ResponseWriter, req *http.Request, id, action string) {
	actor := req.Header.Get("x-actor")
	if actor == "" {
		actor = "operator"
	}
	var r domain.Record
	var err error
	switch action {
	case "submit":
		r, err = h.review.Submit(id, actor)
	case "confirm":
		r, err = h.review.Confirm(id, actor)
	case "archive":
		r, err = h.review.Archive(id, actor)
	case "retry":
		r, err = h.review.RetryConfirmation(id, actor, "http-retry")
	default:
		http.NotFound(w, req)
		return
	}
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, r)
}
func (h *Handler) importCSV(w http.ResponseWriter, req *http.Request) {
	rows, invalid, err := h.importer.ParseCSV(req.Body)
	if err != nil {
		writeErr(w, err)
		return
	}
	result := h.catalog.RegisterMany(rows)
	writeJSON(w, http.StatusOK, map[string]any{"created": result.Created, "failed": append(invalid, result.Failed...)})
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}
