package service

import (
	"URLShortener/internal/storage/memory"
	"strings"
	"testing"
)

func makeService() *Service {
	return NewService(memory.NewStorage())
}

func TestService_SaveURL_ValidURL(t *testing.T) {
	service := makeService()

	urls := []string{
		"https://google.com",
		"http://http.badssl.com",
	}

	for _, url := range urls {
		alias, err := service.SaveURL(url)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", url, err)
		}

		if len(alias) != aliasLength {
			t.Errorf("expected alias length %d, got %d",
				aliasLength, len(alias))
		}
	}
}

func TestService_SaveURL_URLLength(t *testing.T) {
	service := makeService()

	validURL := "https://" + strings.Repeat("a", 2036) + ".com"
	_, err := service.SaveURL(validURL)
	if err != nil {
		t.Errorf("expected URL with length <= 2048 to be valid: %v", err)
	}

	invalidURL := "https://" + strings.Repeat("a", 2037) + ".com"
	_, err = service.SaveURL(invalidURL)
	if err == nil {
		t.Error("expected error for URL longer than 2048")
	}
}

func TestService_SaveURL_AliasNotRepeat(t *testing.T) {
	service := makeService()

	url1 := "https://google.com"
	url2 := "https://github.com"

	alias1, err := service.SaveURL(url1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	alias2, err := service.SaveURL(url2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if alias1 == alias2 {
		t.Errorf("expected different alias, but got same alias: %v", alias1)
	}
}

func TestService_SaveURL_TrimsSpaces(t *testing.T) {
	service := makeService()

	originalURL := "  https://google.com  "
	expectedURL := "https://google.com"

	alias, err := service.SaveURL(originalURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if alias == "" {
		t.Error("expected non-empty alias")
	}

	url, err := service.GetURL(alias)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if url != expectedURL {
		t.Errorf("expected trimmed URL, got %q", url)
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid https",
			url:     "https://google.com",
			wantErr: false,
		},
		{
			name:    "valid http",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "empty",
			url:     "",
			wantErr: true,
		},
		{
			name:    "localhost",
			url:     "http://localhost",
			wantErr: true,
		},
		{
			name:    "localhostIP",
			url:     "http://127.0.0.1",
			wantErr: true,
		},
		{
			name:    "onlyScheme",
			url:     "https://",
			wantErr: true,
		},
		{
			name:    "oneWord",
			url:     "google",
			wantErr: true,
		},
		{
			name:    "noTLD",
			url:     "https://google",
			wantErr: true,
		},
		{
			name:    "noScheme",
			url:     "google.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr = %v",
					err, tt.wantErr)
			}
		})
	}
}
