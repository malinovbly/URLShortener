package memory

import (
	"URLShortener/internal/storage"
	"maps"
	"sync"
)

type Storage struct {
	mu         sync.RWMutex
	urls       map[string]string
	urlToAlias map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		urls:       make(map[string]string),
		urlToAlias: make(map[string]string),
	}
}

func (s *Storage) SaveURL(alias string, originalURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingAlias, ok := s.urlToAlias[originalURL]
	if ok {
		return existingAlias, nil
	}
	if _, ok = s.urls[alias]; ok {
		return "", storage.ErrAliasAlreadyExists
	}

	s.urls[alias] = originalURL
	s.urlToAlias[originalURL] = alias

	return alias, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.urls[alias]
	if !ok {
		return "", storage.ErrURLNotFound
	}

	return url, nil
}

func (s *Storage) GetAllURLs() (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]string, len(s.urls))
	maps.Copy(result, s.urls)

	return result, nil
}
