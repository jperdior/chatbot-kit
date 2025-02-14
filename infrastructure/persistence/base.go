package persistence

import (
	"time"

	"github.com/jperdior/chatbot-kit/domain"
	"gorm.io/gorm"
)

type Base struct {
	ID        []byte    `gorm:"type:binary(16);primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBase(id domain.UuidValueObject) (*Base, error) {
	currentTime := time.Now()
	return &Base{
		ID:        id.Bytes(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}, nil
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// ExecuteTransaction runs the provided function within a transaction.
func (r *TransactionRepository) ExecuteTransaction(txFunc func(tx *gorm.DB) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return txFunc(tx)
	})
}
