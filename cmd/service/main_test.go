package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDataHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/data", nil)
	w := httptest.NewRecorder()
	dataHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json got %s", ct)
	}
}

