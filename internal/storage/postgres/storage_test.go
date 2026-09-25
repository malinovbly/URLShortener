package postgres

import (
	"URLShortener/internal/config"
	storagePackage "URLShortener/internal/storage"
	"errors"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestStorage_SaveURL(t *testing.T) {
	storage := newTestStorage(t)

	alias, err := storage.SaveURL(
		"abc123",
		"https://google.com",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if alias != "abc123" {
		t.Errorf("expected alias %q, got %q", "abc123", alias)
	}

	clearDatabase(t, storage)
}

func TestStorage_SaveURL_ExistingAlias(t *testing.T) {
	storage := newTestStorage(t)

	_, err := storage.SaveURL(
		"abc123",
		"https://google.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = storage.SaveURL(
		"abc123",
		"https://github.com",
	)

	if !errors.Is(err, storagePackage.ErrAliasAlreadyExists) {
		t.Errorf(
			"expected ErrAliasAlreadyExists, got %v",
			err,
		)
	}

	clearDatabase(t, storage)
}

func TestStorage_SaveURL_ExistingURL(t *testing.T) {
	storage := newTestStorage(t)

	firstAlias, err := storage.SaveURL(
		"abc123",
		"https://google.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secondAlias, err := storage.SaveURL(
		"xyz789",
		"https://google.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if secondAlias != firstAlias {
		t.Errorf(
			"expected existing alias %q, got %q",
			firstAlias,
			secondAlias,
		)
	}

	clearDatabase(t, storage)
}

func TestStorage_GetURL(t *testing.T) {
	postgresStorage := newTestStorage(t)

	_, err := postgresStorage.SaveURL(
		"abc123",
		"https://google.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	originalURL, err := postgresStorage.GetURL("abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if originalURL != "https://google.com" {
		t.Errorf(
			"expected URL %q, got %q",
			"https://google.com",
			originalURL,
		)
	}

	clearDatabase(t, postgresStorage)
}

func TestStorage_GetURL_NotFound(t *testing.T) {
	postgresStorage := newTestStorage(t)

	_, err := postgresStorage.GetURL("abc123")

	if !errors.Is(err, storagePackage.ErrURLNotFound) {
		t.Errorf(
			"expected ErrURLNotFound, got %v",
			err,
		)
	}

	clearDatabase(t, postgresStorage)
}

func TestStorage_GetAllURLs(t *testing.T) {
	postgresStorage := newTestStorage(t)

	expected := map[string]string{
		"abc123": "https://google.com",
		"xyz789": "https://github.com",
	}

	for alias, originalURL := range expected {
		_, err := postgresStorage.SaveURL(alias, originalURL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	actual, err := postgresStorage.GetAllURLs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(
			"expected URLs %v, got %v",
			expected,
			actual,
		)
	}

	clearDatabase(t, postgresStorage)
}

func newTestStorage(t *testing.T) *Storage {
	t.Helper()

	cfg := config.PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "url_shortener",
	}

	storage, err := NewStorage(cfg)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	t.Cleanup(func() {
		if err := storage.Close(); err != nil {
			t.Errorf("failed to close storage: %v", err)
		}
	})

	clearDatabase(t, storage)

	return storage
}

func clearDatabase(t *testing.T, storage *Storage) {
	t.Helper()

	err := storage.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Exec("TRUNCATE TABLE urls").
		Error

	if err != nil {
		t.Fatalf("failed to clear database: %v", err)
	}
}
