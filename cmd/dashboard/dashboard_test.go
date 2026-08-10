package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestEventsEndpointStreamsSSE(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/events", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(eventsHandler)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("events handler returned status %d", w.Code)
	}

	body := w.Body.String()
	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Fatalf("events handler content type = %q, want text/event-stream", w.Header().Get("Content-Type"))
	}
	if !strings.Contains(body, "event:") {
		t.Fatalf("events handler body = %q, want SSE event payload", body)
	}
}

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

func TestStatusEndpoint(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(statusHandler)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status handler returned status %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "repository") {
		t.Fatalf("status handler body = %q, want workflow metadata", body)
	}

	if !strings.Contains(body, "passing") {
		t.Fatalf("status handler body = %q, want passing status", body)
	}
}

func TestWorkflowRunEndpoint(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/workflows/run", strings.NewReader(`{"workflow":"production-ci"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(workflowRunHandler)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("workflow run handler returned status %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "runId") {
		t.Fatalf("workflow run body = %q, want runId", body)
	}
}

func TestWorkflowDefinitionEndpoint(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/workflow-file", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(workflowFileHandler)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("workflow file handler returned status %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "name: Production CI") {
		t.Fatalf("workflow file body = %q, want workflow definition", body)
	}
}

func TestDashboardServesInteractiveWorkflowSteps(t *testing.T) {
	path := filepath.Join("static", "index.html")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	http.ServeFile(w, r, path)

	if w.Code != http.StatusOK {
		t.Fatalf("dashboard page returned status %d", w.Code)
	}

	body := w.Body.String()
	for _, marker := range []string{"data-step=\"checkout\"", "data-step=\"benchmarks\"", "data-step=\"artifacts\""} {
		if !strings.Contains(body, marker) {
			t.Fatalf("dashboard page body = %q, want marker %q", body, marker)
		}
	}
}

func TestDashboardServesEngineeringConsoleSections(t *testing.T) {
	path := filepath.Join("static", "index.html")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	http.ServeFile(w, r, path)

	if w.Code != http.StatusOK {
		t.Fatalf("dashboard page returned status %d", w.Code)
	}

	body := w.Body.String()
	for _, marker := range []string{"Validation Trace", "Reference Validation", "Release Readiness", "Validation Matrix", "Artifact Bundle"} {
		if !strings.Contains(body, marker) {
			t.Fatalf("dashboard page body = %q, want marker %q", body, marker)
		}
	}
}
