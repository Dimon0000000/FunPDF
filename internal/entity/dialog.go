package entity

import (
	"FunPDF/internal/entity/models"
)

type Dialog struct {
	ID        string         `gorm:"primary_key" json:"id"`
	SessionID string         `gorm:"index" json:"session_id"`
	Message   models.Message `gorm:"serializer:json" json:"message"`
	Status    int            `gorm:"index" json:"status"` // check if this dialog is success
	BaseModel
}
