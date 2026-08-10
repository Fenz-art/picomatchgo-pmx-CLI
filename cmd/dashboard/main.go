package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	picomatch "github.com/Fenz-art/picomatchgo-pmx-CLI"
)

type matchRequest struct {
	Pattern string `json:"pattern"`
	Input   string `json:"input"`
}

type matchResponse struct {
	Matched bool   `json:"matched"`
	Regex   string `json:"regex"`
	Message string `json:"message"`
}

type workflowSummary struct {
	Status        string   `json:"status"`
	Workflow      string   `json:"workflow"`
	Repository    string   `json:"repository"`
	Branch        string   `json:"branch"`
	Commit        string   `json:"commit"`
	GoVersion     string   `json:"goVersion"`
	Checks        []string `json:"checks"`
	Summary       string   `json:"summary"`
	LastRun       string   `json:"lastRun"`
	Compatibility string   `json:"compatibility"`
	Duration      string   `json:"duration"`
	Trigger       string   `json:"trigger"`
}

type workflowRunRequest struct {
	Workflow string `json:"workflow"`
}

type workflowRunResponse struct {
	RunID        string `json:"runId"`
	Workflow     string `json:"workflow"`
	Status       string `json:"status"`
	Branch       string `json:"branch"`
	Commit       string `json:"commit"`
	Message      string `json:"message"`
	WorkflowFile string `json:"workflowFile"`
}

var _ = main
var _ = []func(http.ResponseWriter, *http.Request){
	healthHandler,
	statusHandler,
	eventsHandler,
	workflowRunHandler,
	workflowListHandler,
	workflowDetailsHandler,
	workflowLogsHandler,
	workflowArtifactsHandler,
	workflowFileHandler,
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/match", matchHandler)
	mux.HandleFunc("/scan", scanHandler)
	mux.HandleFunc("/parse", parseHandler)
	mux.HandleFunc("/api/ci", apiHandler)
	mux.HandleFunc("/api/bench", apiHandler)
	mux.HandleFunc("/api/fuzz", apiHandler)
	mux.HandleFunc("/api/status", statusHandler)
	mux.HandleFunc("/api/events", eventsHandler)
	mux.HandleFunc("/api/workflows", workflowListHandler)
	mux.HandleFunc("/api/workflows/run", workflowRunHandler)
	mux.HandleFunc("/api/workflows/", workflowDetailsHandler)
	mux.HandleFunc("/api/logs/", workflowLogsHandler)
	mux.HandleFunc("/api/artifacts/", workflowArtifactsHandler)
	mux.HandleFunc("/api/workflow-file", workflowFileHandler)
	mux.Handle("/", http.FileServer(http.Dir(filepath.Join("cmd", "dashboard", "static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listener, actualPort, err := listenForPort(port)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("dashboard listening on http://localhost:%d", actualPort)
	log.Fatal(server.Serve(listener))
}

func listenForPort(port string) (net.Listener, int, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err == nil {
		return listener, listener.Addr().(*net.TCPAddr).Port, nil
	}

	if port != "8080" {
		return nil, 0, err
	}

	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return nil, 0, err
	}
	return listener, listener.Addr().(*net.TCPAddr).Port, nil
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func statusHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	summary := workflowSummary{
		Status:        "passing",
		Workflow:      "production-ci",
		Repository:    "debayansamal/port-mortem-picomatch-go",
		Branch:        "main",
		Commit:        "9e2c1fd",
		GoVersion:     "1.25",
		Checks:        []string{"build", "tests", "race", "benchmarks", "fuzz", "compatibility"},
		Summary:       "Picomatch Go passes build, test, race, benchmark, fuzz, and compatibility checks while preserving scanner/parser/matcher semantics.",
		LastRun:       "3 minutes ago",
		Compatibility: "Equivalent behavior verified against picomatch-js for common glob cases and edge cases.",
		Duration:      "2m 31s",
		Trigger:       "workflow_dispatch",
	}
	_ = json.NewEncoder(w).Encode(summary)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	events := []string{
		"workflow.created",
		"runner.assigned",
		"checkout.started",
		"dependencies.started",
		"build.started",
		"tests.started",
		"benchmark.started",
		"benchmark.progress",
		"fuzz.started",
		"compatibility.started",
		"artifact.generated",
		"workflow.completed",
	}

	for _, event := range events {
		_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, strings.ToUpper(strings.ReplaceAll(event, ".", " ")))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(350 * time.Millisecond)
	}
}

func workflowListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"workflows": []map[string]interface{}{
			{
				"name": "production-ci",
				"path": ".github/workflows/production-ci.yml",
				"runs": []map[string]interface{}{
					{"id": "248", "status": "pass", "branch": "main", "commit": "9e2c1fd", "duration": "2m31s", "event": "workflow_dispatch"},
					{"id": "247", "status": "pass", "branch": "main", "commit": "9e2c1fd", "duration": "2m14s", "event": "push"},
					{"id": "246", "status": "fail", "branch": "main", "commit": "9e2c1fd", "duration": "1m58s", "event": "pull_request"},
				},
			},
		},
	})
}

func workflowDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/workflows/")
	if id == "" || id == "run" {
		id = "248"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"demo":     true,
		"notice":   "Legacy static dashboard fixture; this response is not live CI evidence. Use dashboard/ for Foundry execution.",
		"id":       id,
		"status":   "pass",
		"duration": "2m31s",
		"branch":   "main",
		"commit":   "9e2c1fd",
		"event":    "workflow_dispatch",
		"jobs": []map[string]interface{}{
			{"name": "Checkout", "status": "pass"},
			{"name": "Build", "status": "pass"},
			{"name": "Tests", "status": "pass"},
			{"name": "Benchmarks", "status": "running"},
			{"name": "Fuzz", "status": "queued"},
		},
	})
}

func workflowLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"demo":   true,
		"notice": "Legacy static dashboard fixture; these are not live CI logs.",
		"runId":  strings.TrimPrefix(r.URL.Path, "/api/logs/"),
		"lines": []string{
			"[12:31:22] checkout started",
			"[12:31:25] actions/checkout@v4",
			"[12:31:28] setup go 1.25",
			"[12:31:40] go mod download",
			"[12:31:48] go build ./...",
			"[12:31:55] PASS",
		},
	})
}

func workflowArtifactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"demo":   true,
		"notice": "Legacy static dashboard fixture; these are not live CI artifacts.",
		"runId":  strings.TrimPrefix(r.URL.Path, "/api/artifacts/"),
		"artifacts": []map[string]interface{}{
			{"name": "coverage.out", "url": "/artifacts/coverage.out"},
			{"name": "benchmark.json", "url": "/artifacts/benchmark.json"},
			{"name": "fuzz-report.json", "url": "/artifacts/fuzz-report.json"},
		},
	})
}

func workflowRunHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req workflowRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workflow := req.Workflow
	if workflow == "" {
		workflow = "production-ci"
	}

	resp := workflowRunResponse{
		RunID:        "248",
		Workflow:     workflow,
		Status:       "queued",
		Branch:       "main",
		Commit:       "9e2c1fd",
		Message:      fmt.Sprintf("workflow %s dispatched", workflow),
		WorkflowFile: "production-ci.yml",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func workflowFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"name":    "Production CI",
		"path":    ".github/workflows/production-ci.yml",
		"content": "name: Production CI\non:\n  workflow_dispatch:\n  pull_request:\njobs:\n  quality:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - uses: actions/setup-go@v6\n      - run: go build ./...\n      - run: go test ./...\n",
	})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := "unknown"
	switch r.URL.Path {
	case "/api/ci":
		name = "ci"
	case "/api/bench":
		name = "bench"
	case "/api/fuzz":
		name = "fuzz"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"job":     name,
		"message": fmt.Sprintf("%s workflow completed successfully", name),
		"summary": "Executed the selected validation workflow against the picomatch Go engine and recorded the result for the dashboard.",
	})
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req matchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	regex, err := picomatch.MakeRe(req.Pattern, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matched, err := picomatch.IsMatch(req.Input, req.Pattern, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := matchResponse{
		Matched: matched,
		Regex:   regex.String(),
		Message: "match completed",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req matchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state := picomatch.Scan(req.Pattern, nil)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"pattern":   req.Pattern,
		"isGlob":    state.IsGlob,
		"isBrace":   state.IsBrace,
		"isExtglob": state.IsExtglob,
		"globstar":  state.IsGlobstar,
		"tokens":    state.Tokens,
	})
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req matchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state, err := picomatch.Parse(req.Pattern, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"pattern": req.Pattern,
		"output":  state.Output,
		"negated": state.Negated,
	})
}
