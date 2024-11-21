package service

import (
	"GoMJTrainingCamp/dbs/dbConnection"
	"GoMJTrainingCamp/dbs/models"
	"fmt"
)

type TransactionServiceInterface interface {
	CreateTransaction(transaction *models.Transaction) error
}
type TransactionService struct{}

func NewTransactionService() TransactionServiceInterface {
	return &TransactionService{}
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	if err := dbConnection.DB.Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}
