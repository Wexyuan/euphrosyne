package id

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
)

// NewUUID returns a new unique UUID string.
func NewUUID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("[id] generate uuid error: %w", err)
	}
	return id.String(), nil
}
