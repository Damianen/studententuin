package domain

import "errors"

// Sentinel errors the infra adapters translate Docker errors into, so the
// service and HTTP layers never import Docker packages.
var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("already exists")
	ErrInvalid      = errors.New("invalid input")
	ErrImageMissing = errors.New("image not present on host")
)
