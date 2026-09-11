package storage

import "errors"

var (
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrURLNotFound        = errors.New("URL not found")
)
