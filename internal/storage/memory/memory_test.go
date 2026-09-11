package memory

import (
	storagePackage "URLShortener/internal/storage"
	"errors"
	"testing"
)

func TestStorage_SaveURL(t *testing.T) {
	storage := NewStorage()

	alias := "abc123"
	originalURL := "https://google.com"

	savedAlias, err := storage.SaveURL(alias, originalURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedAlias != alias {
		t.Fatalf("expected alias %q, got %q", alias, savedAlias)
	}

	url, err := storage.GetURL(alias)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != originalURL {
		t.Errorf("expected URL %q, got %q", originalURL, url)
	}
}

func TestStorage_SaveURL_ExistingAlias(t *testing.T) {
	storage := NewStorage()

	alias := "abc123"
	url1 := "https://google.com"
	url2 := "https://github.com"

	_, err := storage.SaveURL(alias, url1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = storage.SaveURL(alias, url2)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, storagePackage.ErrAliasAlreadyExists) {
		t.Errorf("expected ErrAliasAlreadyExists, got %v", err)
	}
}

func TestStorage_SaveURL_ExistingURL(t *testing.T) {
	storage := NewStorage()

	alias1 := "abc123"
	alias2 := "xyz789"
	originalURL := "https://google.com"

	savedAlias, err := storage.SaveURL(alias1, originalURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedAlias != alias1 {
		t.Fatalf("expected first alias %q, got %q", alias1, savedAlias)
	}

	savedAlias, err = storage.SaveURL(alias2, originalURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if savedAlias != alias1 {
		t.Errorf("expected existing alias %q, got %q", alias1, savedAlias)
	}
}

func TestStorage_SaveURL_DifferentURLs(t *testing.T) {
	storage := NewStorage()

	alias1, err := storage.SaveURL("abc123", "https://google.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	alias2, err := storage.SaveURL("xyz789", "https://github.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if alias1 == alias2 {
		t.Errorf("expected different aliases, got %q", alias1)
	}
}

func TestStorage_GetURL_NotFound(t *testing.T) {
	storage := NewStorage()

	_, err := storage.GetURL("unknown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, storagePackage.ErrURLNotFound) {
		t.Errorf("expected ErrURLNotFound, got %v", err)
	}
}

func TestStorage_GetAllURLs(t *testing.T) {
	storage := NewStorage()

	alias1, err := storage.SaveURL("abc123", "https://google.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	alias2, err := storage.SaveURL("xyz789", "https://github.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	urls, err := storage.GetAllURLs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(urls) != 2 {
		t.Errorf("expected 2 URLs, got %d", len(urls))
	}
	if _, ok := urls[alias1]; !ok {
		t.Errorf("expected alias %q to exist", alias1)
	}
	if _, ok := urls[alias2]; !ok {
		t.Errorf("expected alias %q to exist", alias2)
	}
}

func TestStorage_GetAllURLs_ChangeMapValue(t *testing.T) {
	storage := NewStorage()

	alias := "abc123"
	originalURL := "https://google.com"
	changedURL := "https://github.com"

	_, err := storage.SaveURL(alias, originalURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	urls, err := storage.GetAllURLs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	urls[alias] = changedURL
	urls, err = storage.GetAllURLs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if urls[alias] != originalURL {
		t.Errorf("expected URL %q, got %q", originalURL, urls[alias])
	}
}

func TestStorage_GetAllURLs_DuplicatedURLs(t *testing.T) {
	storage := NewStorage()

	alias1 := "abc123"
	alias2 := "xyz789"
	originalURL := "https://google.com"

	_, err := storage.SaveURL(alias1, originalURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = storage.SaveURL(alias2, originalURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	urls, err := storage.GetAllURLs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(urls) != 1 {
		t.Errorf("expected 1 URL, got %d", len(urls))
	}
	if _, ok := urls[alias2]; ok {
		t.Errorf("expected alias %q not to exist", alias2)
	}
}

func TestStorage_GetAllURLs_NewStorageIsEmpty(t *testing.T) {
	storage := NewStorage()

	urls, err := storage.GetAllURLs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(urls) != 0 {
		t.Errorf("expected 0 URLs, got %d", len(urls))
	}
}
