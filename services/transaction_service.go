package services

import (
	"errors"
	"fmt"
	"transaction_management/models"
	"transaction_management/repositories"
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionService(transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

func (service *TransactionService) GetTransactionsByMerchantID(merchantID string) ([]models.Transaction, error) {
	return service.transactionRepo.GetByMerchantID(merchantID)
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
