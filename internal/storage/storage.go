package storage

type Storage interface {
	SaveURL(alias string, originalURL string) (string, error)
	GetURL(alias string) (string, error)
	GetAllURLs() (map[string]string, error)
}
