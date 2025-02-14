package domain

import (
	"database/sql/driver"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

// EmailValueObject represents a value object for emails
type EmailValueObject string

func NewEmailValueObject(value string) (EmailValueObject, error) {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(value) {
		return "", NewDomainError("invalid email", "email.invalid")
	}
	return EmailValueObject(value), nil
}

func (emailValueObject *EmailValueObject) Value() string {
	return string(*emailValueObject)
}

// UuidValueObject represents a value object for UUIDs
type UuidValueObject uuid.UUID

func NewUuidValueObject(value string) (UuidValueObject, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return UuidValueObject{}, err
	}
	return UuidValueObject(uid), nil
}

func NewUuidValueObjectFromBytes(value []byte) (UuidValueObject, error) {
	uid, err := uuid.FromBytes(value)
	if err != nil {
		return UuidValueObject{}, err
	}
	return UuidValueObject(uid), nil
}

func RandomUuidValueObject() UuidValueObject {
	return UuidValueObject(uuid.New())
}

func (uidValueObject *UuidValueObject) Scan(value interface{}) error {
	if value == nil {
		*uidValueObject = UuidValueObject{} // If value is nil, set it to the zero value
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

	*uidValueObject = UuidValueObject(parsedUUID)
	return nil
}

func (uidValueObject *UuidValueObject) Value() (driver.Value, error) {
	return uuid.UUID(*uidValueObject).MarshalBinary()
}

func (uidValueObject *UuidValueObject) String() string {
	return uuid.UUID(*uidValueObject).String()
}

func (uidValueObject *UuidValueObject) Bytes() []byte {
	bytes, _ := uuid.UUID(*uidValueObject).MarshalBinary()
	return bytes
}

type SortDirValueObject string

func NewSortDirValueObject(value string) (SortDirValueObject, error) {
	if value == "" {
		return SortDirValueObject("desc"), nil
	}
	if value != "asc" && value != "desc" {
		return "", NewDomainError("invalid sort direction", "sort.invalid")
	}
	return SortDirValueObject(value), nil
}

func (sortDirValueObject *SortDirValueObject) Value() string {
	return string(*sortDirValueObject)
}

type PageValueObject int

func NewPageValueObject(value int) (PageValueObject, error) {
	if value < 1 {
		return PageValueObject(1), nil
	}
	return PageValueObject(value), nil
}

func (pageValueObject *PageValueObject) Value() int {
	return int(*pageValueObject)
}

type PageSizeValueObject int

func NewPageSizeValueObject(value int) (PageSizeValueObject, error) {
	if value < 1 {
		return PageSizeValueObject(25), nil
	}
	if value > 100 {
		return -1, NewDomainError("page size must be less than or equal to 100", "page_size.invalid")
	}
	return PageSizeValueObject(value), nil
}

func (pageSizeValueObject *PageSizeValueObject) Value() int {
	return int(*pageSizeValueObject)
}
