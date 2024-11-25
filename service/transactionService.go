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
	GetTransactionAll() ([]models.Transaction, error)
	GetTransactionById(id uint) (*GetTransactionById, error)
	GetTransactionByUser(userID uint) ([]GetTransactionById, error)
}
type TransactionService struct{}

func NewTransactionService() TransactionServiceInterface {
	return &TransactionService{}
}

type GetTransactionById struct {
	TrasactionID uint                `json:"idTransaction"`
	Transaction  *models.Transaction `json:"transaction"`
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction) (uint, error) {
	if err := dbConnection.DB.Create(transaction).Error; err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

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
func (s *TransactionService) GetTransactionAll() ([]models.Transaction, error) {
	var transaction []models.Transaction
	if err := dbConnection.DB.Find(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get transaction failed")
		}
	}
	return transaction, nil
}
func (s *TransactionService) GetTransactionById(id uint) (*GetTransactionById, error) {
	var transaction models.Transaction

	// Preload related VisitPackage and Membership
	if err := dbConnection.DB.
		Preload("VisitPackage").Preload("Membership").
		First(&transaction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}
	response := &GetTransactionById{
		TrasactionID: transaction.IDTransaction,
		Transaction:  &transaction,
	}
	return response, nil
}
func (s *TransactionService) GetTransactionByUser(userID uint) ([]GetTransactionById, error) {
	var transactions []models.Transaction

	err := dbConnection.DB.Raw(`
		SELECT t.* 
		FROM transactions t
		LEFT JOIN memberships m ON t.membership_id = m.id_membership
		LEFT JOIN visit_packages v ON t.Visit_id = v.id_visit_package
		WHERE m.user_id = ? OR v.user_id = ?
	`, userID, userID).Scan(&transactions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions for user ID %d: %w", userID, err)
	}

	var results []GetTransactionById
	for _, transaction := range transactions {
		response := GetTransactionById{
			TrasactionID: transaction.IDTransaction,
			Transaction:  &transaction,
		}
		results = append(results, response)
	}

	return results, nil
}
