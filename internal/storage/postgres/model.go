package postgres

type URL struct {
	Alias       string `gorm:"primaryKey;column:alias"`
	OriginalURL string `gorm:"uniqueIndex;column:original_url"`
}

func (URL) TableName() string {
	return "urls"
}
