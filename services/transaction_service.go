package services

import (
	"errors"
	"fmt"
	"time"
	"transaction_management/models"
	"transaction_management/repositories"
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionService(transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

func (service *TransactionService) GetTransactionByID(id string) (*models.Transaction, error) {
	return service.transactionRepo.GetByID(id)
}

func (service *TransactionService) GetTransactionsByMerchantID(merchantID string) ([]models.Transaction, error) {
	return service.transactionRepo.GetByMerchantID(merchantID)
}
func (service *TransactionService) SearchTransactions(merchantID, description string, createdAt time.Time) ([]models.Transaction, error) {
	return service.transactionRepo.Search(merchantID, description, createdAt)
}

func (service *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	if transaction.MerchantID == "" || transaction.Description == "" || transaction.Amount == 0 || transaction.Currency == "" {
		return errors.New("Missing required fields")
	}

	fmt.Println("Transaction ID before inserting:", transaction.ID.Hex())
	fmt.Println(transaction)
	err := service.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (service *TransactionService) GetAllTransactions() ([]models.Transaction, error) {
	return service.transactionRepo.GetAllTransactions()
}
