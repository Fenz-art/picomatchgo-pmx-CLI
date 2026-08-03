package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	picomatch "github.com/debayansamal/port-mortem-picomatch-go"
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/match", matchHandler)
	mux.HandleFunc("/scan", scanHandler)
	mux.HandleFunc("/parse", parseHandler)
	mux.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("dashboard listening on http://localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Picomatch Go Dashboard")
	_, _ = fmt.Fprintln(w, "Available endpoints: /health, /match, /scan, /parse")
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
