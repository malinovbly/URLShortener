package redirect

import (
	"URLShortener/internal/http/response"
	"URLShortener/internal/storage"
	"URLShortener/internal/testutil"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockService struct {
	saveURL    func(string) (string, error)
	getAllURLs func() (map[string]string, error)
	getURL     func(string) (string, error)
}

func (m *mockService) SaveURL(originalURL string) (string, error) {
	return m.saveURL(originalURL)
}

func (m *mockService) GetAllURLs() (map[string]string, error) {
	return m.getAllURLs()
}

func (m *mockService) GetURL(alias string) (string, error) {
	return m.getURL(alias)
}

func TestHandler_Redirect_Valid(t *testing.T) {
	const (
		alias       = "abc123"
		originalURL = "https://google.com"
	)

	mockService := &mockService{
		getURL: func(gotAlias string) (string, error) {
			if gotAlias != alias {
				t.Errorf("expected alias %q, got %q", alias, gotAlias)
			}

			return originalURL, nil
		},
	}

	handler := NewHandler(mockService)

	request := httptest.NewRequest(http.MethodGet, "/"+alias, nil)
	request.SetPathValue("alias", alias)

	rr := httptest.NewRecorder()

	handler.Redirect(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusFound)
	testutil.CheckLocation(t, rr, originalURL)
}

func TestHandler_Redirect_NotFound(t *testing.T) {
	mockService := &mockService{
		getURL: func(string) (string, error) {
			return "", storage.ErrURLNotFound
		},
	}

	handler := NewHandler(mockService)

	request := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	request.SetPathValue("alias", "abc123")

	rr := httptest.NewRecorder()

	handler.Redirect(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusNotFound)
	testutil.CheckContentType(t, rr, "application/json")

	responseBody := testutil.DecodeJSON[response.ErrorResponse](t, rr.Body)

	testutil.CheckError(t, responseBody, storage.ErrURLNotFound.Error())
}

func TestHandler_Redirect_ServiceError(t *testing.T) {
	mockService := &mockService{
		getURL: func(string) (string, error) {
			return "", errors.New("database error")
		},
	}

	handler := NewHandler(mockService)

	request := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	request.SetPathValue("alias", "abc123")

	rr := httptest.NewRecorder()

	handler.Redirect(rr, request)

	testutil.CheckStatusCode(t, rr, http.StatusInternalServerError)
	testutil.CheckContentType(t, rr, "application/json")

	responseBody := testutil.DecodeJSON[response.ErrorResponse](t, rr.Body)

	testutil.CheckError(t, responseBody, response.ErrorInternalServer)
}
