package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handleRoot(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, appName) {
		t.Errorf("body should contain %q, got: %s", appName, body)
	}
}

func TestHandleReadiness(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()

	handleReadiness(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if body := rec.Body.String(); body != "Ready" {
		t.Errorf("expected 'Ready', got %q", body)
	}
}

func TestHandleInfo(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rec := httptest.NewRecorder()

	handleInfo(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("expected JSON content-type, got %q", ct)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "kiada-go") {
		t.Errorf("body should contain app name, got: %s", body)
	}
}

func TestHandleProxyQuote_Unconfigured(t *testing.T) {
	t.Setenv("QUOTE_URL", "")
	req := httptest.NewRequest(http.MethodGet, "/proxy/quote", nil)
	rec := httptest.NewRecorder()

	handleProxyQuote(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 when QUOTE_URL not set, got %d", rec.Code)
	}
}

func TestHandleProxyQuiz_Unconfigured(t *testing.T) {
	t.Setenv("QUIZ_URL", "")
	req := httptest.NewRequest(http.MethodGet, "/proxy/quiz", nil)
	rec := httptest.NewRecorder()

	handleProxyQuiz(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 when QUIZ_URL not set, got %d", rec.Code)
	}
}

func TestGetEnv_Fallback(t *testing.T) {
	got := getEnv("__KIADA_GO_NONEXISTENT_VAR__", "default")
	if got != "default" {
		t.Errorf("expected 'default', got %q", got)
	}
}

func TestGetEnv_Set(t *testing.T) {
	t.Setenv("__KIADA_GO_TEST_VAR__", "hello")
	got := getEnv("__KIADA_GO_TEST_VAR__", "default")
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}
