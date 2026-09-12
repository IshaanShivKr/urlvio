package service

import "errors"

var (
	ErrURLRequired  = errors.New("url is required")
	ErrInvalidURL   = errors.New("invalid url")
	ErrNotFound     = errors.New("link not found")
)
