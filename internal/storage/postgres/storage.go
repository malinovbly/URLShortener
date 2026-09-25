package postgres

import (
	"errors"
	"fmt"

	"URLShortener/internal/config"
	"URLShortener/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(cfg config.PostgresConfig) (*Storage, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(alias string, originalURL string) (string, error) {
	var existing URL

	err := s.db.
		Where("original_url = ?", originalURL).
		First(&existing).Error

	if err == nil {
		return existing.Alias, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	url := URL{
		Alias:       alias,
		OriginalURL: originalURL,
	}

	err = s.db.Create(&url).Error
	if err != nil {
		if isUniqueViolation(err) {
			return "", storage.ErrAliasAlreadyExists
		}

		return "", err
	}

	return url.Alias, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	var url URL

	err := s.db.
		Where("alias = ?", alias).
		First(&url).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", storage.ErrURLNotFound
		}

		return "", err
	}

	return url.OriginalURL, nil
}

func (s *Storage) GetAllURLs() (map[string]string, error) {
	var records []URL

	if err := s.db.Find(&records).Error; err != nil {
		return nil, err
	}

	urls := make(map[string]string, len(records))

	for _, record := range records {
		urls[record.Alias] = record.OriginalURL
	}

	return urls, nil
}

func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
