package repository

import "errors"

var (
	ErrNotFound      = errors.New("link not found")
	ErrAlreadyExists = errors.New("link already exists")
)
