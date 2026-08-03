package dashboard

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("health handler returned status %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "ok") {
		t.Fatalf("health handler body = %q, want ok", w.Body.String())
	}
}

func TestMatchEndpoint(t *testing.T) {
	payload := `{"pattern":"*.go","input":"main.go"}`
	r := httptest.NewRequest(http.MethodPost, "/match", strings.NewReader(payload))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(matchHandler)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("match handler returned status %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "matched") {
		t.Fatalf("match handler body = %q, want match output", w.Body.String())
	}
}
