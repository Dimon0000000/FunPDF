package engine

import (
	"context"
	"encoding/json"
)

type Translator interface {
	Translate(ctx context.Context, from, to, q string, params json.RawMessage) (string, error)
	Healthy(ctx context.Context) bool
	Name(ctx context.Context) string
}
