package domain

import (
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

type UUIDValueObjectInterface interface {
	Value() string
	RandomUUIDValueObject() UUIDValueObjectInterface
}

// UUIDValueObject represents a value object for UUIDs
type UUIDValueObject uuid.UUID

func (uuidValueObject *UUIDValueObject) Value() string {
	return uuid.UUID(*uuidValueObject).String()
}

func (uuidValueObject *UUIDValueObject) RandomUUIDValueObject() UUIDValueObjectInterface {
	newUUID := UUIDValueObject(uuid.New())
	return &newUUID
}

func NewUuidValueObject(value string) (UUIDValueObjectInterface, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return &UUIDValueObject{}, err
	}
	uuidValueObject := UUIDValueObject(uid)
	return &uuidValueObject, nil
}

func NewUuidValueObjectFromBytes(value []byte) (UUIDValueObjectInterface, error) {
	uid, err := uuid.FromBytes(value)
	if err != nil {
		return &UUIDValueObject{}, err
	}
	uuidValueObject := UUIDValueObject(uid)
	return &uuidValueObject, nil
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
