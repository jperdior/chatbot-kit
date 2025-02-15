package gorm

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"github.com/jperdior/chatbot-kit/domain"
)

type Adapter interface {
	Scan(value interface{}) error
	Value() (driver.Value, error)
}

type UUIDAdapter struct {
	ValueObject domain.UUIDValueObjectInterface
}

func (uidAdapter *UUIDAdapter) Scan(value interface{}) error {
	if value == nil {
		uidAdapter.ValueObject = domain.UUIDValueObject{} // If value is nil, set it to the zero value
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UUID: expected []byte but got %T", value)
	}

	parsedUUID, err := uuid.FromBytes(bytes)
	if err != nil {
		return fmt.Errorf("failed to parse UUID from bytes: %w", err)
	}

	uidAdapter.ValueObject, _ = domain.NewUuidValueObject(parsedUUID.String())
	return nil
}

func (uidAdapter *UUIDAdapter) Value() (driver.Value, error) {
	if uidAdapter.ValueObject == nil {
		return nil, nil
	}
	return uidAdapter.ValueObject.Value(), nil
}

func (uidAdapter *UUIDAdapter) GormDataType() string {
	return "uuid"
}
