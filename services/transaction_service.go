package services

import (
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
