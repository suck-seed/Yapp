package utils

import (
	"github.com/google/uuid"
)

// ParseUUID parses a string into a UUID and returns a pointer to uuid.UUID
func ParseUUID(idStr string) (*uuid.UUID, error) {
	u, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
