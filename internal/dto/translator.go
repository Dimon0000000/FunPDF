package dto

import "encoding/json"

type CreateTranslatorsRequest struct {
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}

type TranslateRequest struct {
	From   *string         `json:"from"`
	To     *string         `json:"to"`
	Q      *string         `json:"q"`
	Params json.RawMessage `json:"params"`
}
