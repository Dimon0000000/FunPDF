package service

import "errors"

var (
	ErrProviderNotFound       = errors.New("provider not found")
	ErrProviderIDRequired     = errors.New("provider id is empty")
	ErrProviderNameRequired   = errors.New("provider name is empty")
	ErrProviderURLSuffix      = errors.New("provider chat and models url suffix are required")
	ErrModelNameRequired      = errors.New("model name is empty")
	ErrUnsupportedProvider    = errors.New("unsupported provider")
)
