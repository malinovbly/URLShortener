package testutil

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
)

func CheckStatusCode(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	t.Helper()

	actual := response.Code

	if actual != expected {
		t.Errorf(
			"expected status code %d, got %d",
			expected,
			actual,
		)
	}
}

func CheckContentType(t *testing.T, response *httptest.ResponseRecorder, expected string) {
	t.Helper()

	actual := response.Header().Get("Content-Type")

	if actual != expected {
		t.Errorf(
			"expected Content-Type %q, got %q",
			expected,
			actual,
		)
	}
}

func DecodeJSON[T any](t *testing.T, body io.Reader) T {
	t.Helper()

	var result T

	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	return result
}
