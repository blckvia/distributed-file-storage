package storage

import "errors"

var (
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrFileNotFound      = errors.New("file not found")
	ErrAppNotFound       = errors.New("app not found")
)
