package common

import (
	"github.com/google/uuid"
)

// GenerateUUID generate UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUUIDv7 generate UUID v7 (based on timeline, ordered)
func GenerateUUIDv7() string {
	return uuid.Must(uuid.NewV7()).String()
}

// ParseUUID parse the UUID
func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
