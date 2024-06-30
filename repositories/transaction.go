package repositories

import (
	"synapsis/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(order *models.Transaction) error
	UpsertTransaction(order *models.Transaction) error
	GetTransactionsByUserID(userID uint) ([]models.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(order *models.Transaction) error {
	return r.db.Create(order).Error
}

func (r *transactionRepository) UpsertTransaction(order *models.Transaction) error {
	return r.db.Save(order).Error
}

func (r *transactionRepository) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}
