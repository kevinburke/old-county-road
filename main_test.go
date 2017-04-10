package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	mux := NewServeMux()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("GET /: got code %d, want 200", w.Code)
	}
	if body := w.Body.String(); !strings.Contains(body, "Hello World") {
		t.Errorf("GET /: expected 'Hello World' in body, got %s", body)
	}
}
