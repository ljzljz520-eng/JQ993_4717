package api

import (
	"net/http"
	"strings"
)

type Route struct {
	Method  string
	Prefix  string
	Handler http.HandlerFunc
}
type Router struct {
	routes   []Route
	fallback http.Handler
}

func NewRouter(fallback http.Handler) *Router { return &Router{fallback: fallback} }
func (r *Router) Add(method, prefix string, h http.HandlerFunc) {
	r.routes = append(r.routes, Route{method, prefix, h})
}
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if req.Method == route.Method && strings.HasPrefix(req.URL.Path, route.Prefix) {
			route.Handler(w, req)
			return
		}
	}
	if r.fallback != nil {
		r.fallback.ServeHTTP(w, req)
		return
	}
	http.NotFound(w, req)
}
func MethodAllowed(method string, allowed []string) bool {
	for _, candidate := range allowed {
		if method == candidate {
			return true
		}
	}
	return false
}
func NormalizePath(path string) string {
	clean := strings.TrimSpace(path)
	if clean == "" {
		return "/"
	}
	if !strings.HasPrefix(clean, "/") {
		clean = "/" + clean
	}
	return strings.TrimRight(clean, "/")
}
func QueryValue(req *http.Request, key string) string {
	return strings.TrimSpace(req.URL.Query().Get(key))
}
