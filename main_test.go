package main

import (
	"io"
	"io/ioutil"
	"net/http"
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

func BenchmarkHomepage(b *testing.B) {
	mux := NewServeMux()
	s := httptest.NewServer(mux)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res, err := http.Get(s.URL)
		if err != nil {
			b.Fatal(err)
		}
		if res.StatusCode != 200 {
			b.Fatalf("GET /: expected code 200, got %d", res.StatusCode)
		}
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}
}
