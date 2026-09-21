package url

import (
	"URLShortener/internal/http/response"
	servicePackage "URLShortener/internal/service"
	"URLShortener/internal/storage/memory"
	"URLShortener/internal/testutil"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type mockService struct {
	saveURL    func(string) (string, error)
	getAllURLs func() (map[string]string, error)
}

func (m *mockService) SaveURL(originalURL string) (string, error) {
	return m.saveURL(originalURL)
}

func (m *mockService) GetAllURLs() (map[string]string, error) {
	return m.getAllURLs()
}

func TestHandler_CreateURL_Valid(t *testing.T) {
	storage := memory.NewStorage()
	service := servicePackage.NewService(storage)

	handler := NewHandler(service)

	body := bytes.NewBufferString(`{"url":"https://google.com"}`)
	request := httptest.NewRequest(http.MethodPost, "/url", body)
	rr := httptest.NewRecorder()

	handler.CreateURL(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusCreated)
	testutil.CheckContentType(t, rr, "application/json")

	responseBody := testutil.DecodeJSON[CreateURLResponse](t, rr.Body)

	if responseBody.Alias == "" {
		t.Error("expected alias, got empty")
	}
}

func TestHandler_CreateURL_InvalidJSON(t *testing.T) {
	storage := memory.NewStorage()
	service := servicePackage.NewService(storage)

	handler := NewHandler(service)

	body := strings.NewReader(`{"url":`)
	request := httptest.NewRequest(http.MethodPost, "/url", body)
	rr := httptest.NewRecorder()

	handler.CreateURL(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusBadRequest)
	testutil.CheckContentType(t, rr, "application/json")

	responseBody := testutil.DecodeJSON[response.ErrorResponse](t, rr.Body)

	testutil.CheckError(t, responseBody, response.ErrorInvalidRequestBody)
}

func TestHandler_CreateURL_InvalidURL(t *testing.T) {
	mockService := &mockService{
		saveURL: func(string) (string, error) {
			return "", fmt.Errorf(
				"validation failed: %w",
				servicePackage.ErrInvalidURL,
			)
		},
	}

	handler := NewHandler(mockService)

	body := bytes.NewBufferString(`{"url":"https://google.com"}`)
	request := httptest.NewRequest(http.MethodPost, "/url", body)
	rr := httptest.NewRecorder()

	handler.CreateURL(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusBadRequest)
	testutil.CheckContentType(t, rr, "application/json")

	responseBody := testutil.DecodeJSON[response.ErrorResponse](t, rr.Body)
	expectedError := "validation failed: " + servicePackage.ErrInvalidURL.Error()

	testutil.CheckError(t, responseBody, expectedError)
}

func TestHandler_CreateURL_ServiceError(t *testing.T) {
	mockService := &mockService{
		saveURL: func(string) (string, error) {
			return "", errors.New("database error")
		},
	}

	handler := NewHandler(mockService)

	body := bytes.NewBufferString(`{"url":"https://google.com"}`)
	request := httptest.NewRequest(http.MethodPost, "/url", body)
	rr := httptest.NewRecorder()

	handler.CreateURL(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusInternalServerError)
	testutil.CheckContentType(t, rr, "application/json")

	responseBody := testutil.DecodeJSON[response.ErrorResponse](t, rr.Body)

	testutil.CheckError(t, responseBody, response.ErrorInternalServer)
}

func TestHandler_AllURLs_Valid(t *testing.T) {
	expectedURLs := map[string]string{
		"abc123": "https://google.com",
		"xyz789": "https://github.com",
	}

	mockService := &mockService{
		getAllURLs: func() (map[string]string, error) {
			return expectedURLs, nil
		},
	}

	handler := NewHandler(mockService)

	request := httptest.NewRequest(http.MethodGet, "/urls", nil)
	rr := httptest.NewRecorder()

	handler.AllURLs(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusOK)
	testutil.CheckContentType(t, rr, "application/json")

	responseBody := testutil.DecodeJSON[AllURLsResponse](t, rr.Body)

	if !reflect.DeepEqual(responseBody.URLs, expectedURLs) {
		t.Errorf(
			"expected URLs %v, got %v",
			expectedURLs,
			responseBody.URLs,
		)
	}
}

func TestHandler_AllURLs_ServiceError(t *testing.T) {
	mockService := &mockService{
		getAllURLs: func() (map[string]string, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(mockService)

	request := httptest.NewRequest(http.MethodGet, "/urls", nil)
	rr := httptest.NewRecorder()

	handler.AllURLs(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusInternalServerError)
	testutil.CheckContentType(t, rr, "application/json")

	responseBody := testutil.DecodeJSON[response.ErrorResponse](t, rr.Body)

	testutil.CheckError(t, responseBody, response.ErrorInternalServer)
}
