package entity

import "encoding/json"

type Translator struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
	BaseModel
}
