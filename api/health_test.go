package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	router := NewServer()

	// Test GET /health
	t.Run("GET /health", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"status": "ok"}`, w.Body.String())
	})

	// Test HEAD /health
	t.Run("HEAD /health", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("HEAD", "/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		// httptest records the body written by the handler, even if it's a HEAD request.
		// In a real server, net/http strips the body.
		// We just verify we get a 200 OK, which fixes the reported 404.
	})
}
