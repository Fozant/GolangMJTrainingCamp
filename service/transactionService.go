package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type TransactionServiceInterface interface {
	CreateTransaction(transaction *models.Transaction) (uint, error)
	UpdateTransaction(transaction *models.Transaction) error
}
type TransactionService struct{}

func NewTransactionService() TransactionServiceInterface {
	return &TransactionService{}
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction) (uint, error) {
	if err := dbConnection.DB.Create(transaction).Error; err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}
	fmt.Println(transaction.IDTransaction)
	fmt.Println(transaction.IDTransaction)
	fmt.Println(transaction.IDTransaction)
	fmt.Println(transaction.IDTransaction)
	return transaction.IDTransaction, nil
}
func (s *TransactionService) GetById(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := dbConnection.DB.First(&transaction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve transaction: %w", err)
	}
	return &transaction, nil
}

func (s *TransactionService) UpdateTransaction(transaction *models.Transaction) error {
	// Check if the transaction exists
	existingTransaction, err := s.GetById(transaction.IDTransaction)
	if err != nil {
		return fmt.Errorf("transaction not found or failed to retrieve: %w", err)
	}
	// Log or debug if needed
	fmt.Printf("Updating transaction: %+v\n", existingTransaction)

	// Update the transaction record
	if err := dbConnection.DB.Model(&models.Transaction{}).Where("id_transaction = ?", transaction.IDTransaction).Updates(transaction).Error; err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}
