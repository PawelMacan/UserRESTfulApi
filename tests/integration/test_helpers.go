package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// makeRequest is a helper function to make HTTP requests in tests
func makeRequest(t *testing.T, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error
	
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)
	return w
}
