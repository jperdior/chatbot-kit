package domain

type DomainError interface {
	GetKey() string
	Error() string
}
