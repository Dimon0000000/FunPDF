package enetity

import "time"

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Data []byte `json:"data"`

	Base
}

type Base struct {
	CreateAt time.Time
	UpdateAt time.Time
}
