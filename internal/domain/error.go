package domain

import "errors"

var (
	ErrNotFound         = errors.New("entity not found")
	ErrAlreadyExists    = errors.New("entity already exists")
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrCollectionFailed = errors.New("collection failed")
	ErrSourceDisabled   = errors.New("source is disabled")
	ErrProfileNotFound  = errors.New("profile not found")
	ErrSourceNotFound   = errors.New("source not found")
	ErrTrendNotFound    = errors.New("trend not found")
)
