package service

import (
	"URLShortener/internal/storage"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"strings"
)

const (
	aliasLength = 6
	chars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) SaveURL(originalURL string) (string, error) {
	originalURL = strings.TrimSpace(originalURL)

	err := validateURL(originalURL)
	if err != nil {
		return "", err
	}

	alias, err := generateAlias()
	if err != nil {
		return "", err
	}

	savedAlias, err := s.storage.SaveURL(alias, originalURL)
	if err != nil {
		return "", err
	}

	return savedAlias, nil
}

func (s *Service) GetURL(alias string) (string, error) {
	return s.storage.GetURL(alias)
}

func (s *Service) GetAllURLs() (map[string]string, error) {
	return s.storage.GetAllURLs()
}

func generateAlias() (string, error) {
	result := make([]byte, aliasLength)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}

		result[i] = chars[n.Int64()]
	}

	return string(result), nil
}

func validateURL(originalURL string) error {
	if originalURL == "" {
		return fmt.Errorf("%w: URL cannot be empty", ErrInvalidURL)
	}

	u, err := url.Parse(originalURL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: unsupported URL scheme", ErrInvalidURL)
	}

	if u.Hostname() == "" {
		return fmt.Errorf("%w: URL must contain a host", ErrInvalidURL)
	}

	return nil
}
