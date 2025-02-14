package kit

import (
	"errors"
	"github.com/google/uuid"
)

func UuidStringToBinary(id string) ([]byte, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid UUID format")
	}
	return parsedUUID.MarshalBinary()
}

func ContainsString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
