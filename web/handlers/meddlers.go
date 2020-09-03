package handlers

import "net/http"

// Meddlers holds different middleware implementations and provide some components for use in implementations.
type Meddlers struct{}

// NewMeddlers returns a newly created and ready to use Meddlers.
func NewMeddlers() *Meddlers {
	return &Meddlers{}
}

// JSONContentTypeHeaderMiddleware adds application/json content type header to responses.
func (ms *Meddlers) JSONContentTypeHeaderMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		handler.ServeHTTP(w, r)
	})
}
