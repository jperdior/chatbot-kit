package domain

import (
	"regexp"

	"github.com/google/uuid"
)

// EmailValueObject represents a value object for emails
type EmailValueObject struct {
	value string
}

func NewEmailValueObject(value string) (*EmailValueObject, error) {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(value) {
		return nil, NewDomainError("invalid email", "email.invalid")
	}
	return &EmailValueObject{value: value}, nil
}

func (emailValueObject *EmailValueObject) Value() string {
	return emailValueObject.value
}

type UUIDValueObjectInterface interface {
	Value() uuid.UUID
}

// UUIDValueObject represents a value object for UUIDs
type UUIDValueObject struct {
	value uuid.UUID
}

func (u UUIDValueObject) Value() uuid.UUID {
	return u.value
}

func NewRandomUUIDValueObject() *UUIDValueObject {
	uid := uuid.New()
	return &UUIDValueObject{value: uid}
}

func NewUuidValueObject(value string) (*UUIDValueObject, error) {
	uid, err := uuid.Parse(value)
	if err != nil {
		return &UUIDValueObject{}, err
	}
	return &UUIDValueObject{value: uid}, nil
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
