package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// newRouter registers all routes and returns the handler.
// Route design follows the kiada app pattern across Chapters 2–11.
func newRouter() http.Handler {
	mux := http.NewServeMux()

	// Meta / health endpoints
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/healthz/ready", handleReadiness)
	mux.HandleFunc("/info", handleInfo)

	// Proxy endpoints — delegate to upstream services (Ch 11 pattern)
	mux.HandleFunc("/proxy/quote", handleProxyQuote)
	mux.HandleFunc("/proxy/quiz", handleProxyQuiz)

	// Fallback
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return mux
}

// handleRoot returns a plain-text greeting including the pod name (Ch 2/8 pattern).
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	hostname, _ := os.Hostname()
	podName := getEnv("POD_NAME", hostname)
	clientIP := r.RemoteAddr
	log.Printf("GET / from %s", clientIP)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello from %s v%s!\n", appName, version)
	fmt.Fprintf(w, "Pod: %s | Node: %s\n", podName, getEnv("NODE_NAME", "unknown"))
	fmt.Fprintf(w, "Pod IP: %s | Node IP: %s\n", getEnv("POD_IP", "0.0.0.0"), getEnv("NODE_IP", "0.0.0.0"))
	msg := getEnv("INITIAL_STATUS_MESSAGE", "")
	if msg != "" {
		fmt.Fprintf(w, "Status: %s\n", msg)
	}
}

// handleReadiness is the readiness probe endpoint (Ch 6/11 pattern).
// Returns 200 OK when the service is ready to receive traffic.
func handleReadiness(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /healthz/ready from %s", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Ready")
}

// handleInfo returns a JSON object with pod metadata (Ch 8 downward API pattern).
func handleInfo(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	info := map[string]string{
		"app":      appName,
		"version":  version,
		"pod":      getEnv("POD_NAME", hostname),
		"podIP":    getEnv("POD_IP", "0.0.0.0"),
		"node":     getEnv("NODE_NAME", "unknown"),
		"nodeIP":   getEnv("NODE_IP", "0.0.0.0"),
		"status":   getEnv("INITIAL_STATUS_MESSAGE", ""),
		"hostname": hostname,
	}
	log.Printf("GET /info from %s", r.RemoteAddr)
	jsonResponse(w, http.StatusOK, info)
}

// handleProxyQuote forwards requests to the quote upstream service (Ch 11 proxy pattern).
func handleProxyQuote(w http.ResponseWriter, r *http.Request) {
	quoteURL := getEnv("QUOTE_URL", "")
	if quoteURL == "" {
		jsonError(w, http.StatusServiceUnavailable, "QUOTE_URL not configured")
		return
	}
	log.Printf("GET /proxy/quote → %s", quoteURL)
	body, statusCode, err := httpGet(quoteURL)
	if err != nil {
		log.Printf("ERROR proxying quote: %v", err)
		jsonError(w, http.StatusBadGateway, fmt.Sprintf("upstream error: %v", err))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

// handleProxyQuiz forwards requests to the quiz upstream service (Ch 11 proxy pattern).
func handleProxyQuiz(w http.ResponseWriter, r *http.Request) {
	quizURL := getEnv("QUIZ_URL", "")
	if quizURL == "" {
		jsonError(w, http.StatusServiceUnavailable, "QUIZ_URL not configured")
		return
	}
	upstream := quizURL + "/questions/random"
	log.Printf("GET /proxy/quiz → %s", upstream)
	body, statusCode, err := httpGet(upstream)
	if err != nil {
		log.Printf("ERROR proxying quiz: %v", err)
		jsonError(w, http.StatusBadGateway, fmt.Sprintf("upstream error: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

// httpGet performs a simple HTTP GET and returns the response body and status code.
func httpGet(url string) (string, int, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	return string(data), resp.StatusCode, nil
}

// jsonResponse writes a JSON-encoded body with the given status code.
func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// jsonError writes a JSON error response.
func jsonError(w http.ResponseWriter, statusCode int, message string) {
	jsonResponse(w, statusCode, map[string]string{"error": message})
}
