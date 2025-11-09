package storage

import "errors"

var (
	ErrTTLIsExpired = errors.New("ttl is expired")
	ErrTextNotFound = errors.New("text is not found")
)
