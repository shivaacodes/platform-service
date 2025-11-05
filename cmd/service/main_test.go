package main

import (
"net/http"
"net/http/httptest"
"testing"
"time"
)

// MockCacheClient is a mock implementation of the cache client.
type MockCacheClient struct {
GetFunc   func(key string) (string, error)
SetFunc   func(key string, value interface{}, ttl time.Duration) error
CloseFunc func() error
PingFunc  func() error
}

func (m *MockCacheClient) Get(key string) (string, error) {
return m.GetFunc(key)
}

func (m *MockCacheClient) Set(key string, value interface{}, ttl time.Duration) error {
return m.SetFunc(key, value, ttl)
}

func (m *MockCacheClient) Close() error {
return m.CloseFunc()
}

func (m *MockCacheClient) Ping() error {
return m.PingFunc()
}

func TestDataHandler(t *testing.T) {
// Create a mock cache client.
mockCache := &MockCacheClient{}

// Set the mock functions.
mockCache.GetFunc = func(key string) (string, error) {
return "", nil // Simulate a cache miss.
}
mockCache.SetFunc = func(key string, value interface{}, ttl time.Duration) error {
return nil
}
mockCache.PingFunc = func() error {
return nil
}

// Create a new request and response recorder.
req := httptest.NewRequest(http.MethodGet, "/api/v1/data", nil)
w := httptest.NewRecorder()

// Create the handler and serve the request.
http.HandlerFunc(dataHandler(mockCache)).ServeHTTP(w, req)

// Check the response.
res := w.Result()
if res.StatusCode != http.StatusOK {
t.Fatalf("expected 200 got %d", res.StatusCode)
}
if ct := res.Header.Get("Content-Type"); ct != "application/json" {
t.Fatalf("expected application/json got %s", ct)
}
}