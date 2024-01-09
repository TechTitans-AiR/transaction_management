package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

func (service *TransactionService) GetTransactionByID(id string, token string) (*models.Transaction, error) {
	userRole, err := service.CheckUserRoleFromToken(token)
	if err != nil {
		return nil, err
	}

	if userRole != "admin" {
		return nil, errors.New("only admin users can perform this action")
	}

	return service.transactionRepo.GetByID(id)
}


func (service *TransactionService) GetTransactionsByMerchantID(token, requestedMerchantID string) ([]models.Transaction, error) {
	userRole, err := service.CheckUserRoleFromToken(token)
	if err != nil {
		return nil, err
	}

	if userRole == "admin" {
		return service.transactionRepo.GetByMerchantID(requestedMerchantID)
	}

	merchantID, err := service.GetMerchantIDFromToken(token)
	if err != nil {
		return nil, err
	}

	if requestedMerchantID != merchantID {
		return nil, errors.New("you can only retrieve transactions for your own merchant ID")
	}

	return service.transactionRepo.GetByMerchantID(requestedMerchantID)
}

func (service *TransactionService) SearchTransactions(merchantID, description string, createdAt time.Time) ([]models.Transaction, error) {
	return service.transactionRepo.Search(merchantID, description, createdAt)
}



func (service *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	if transaction.MerchantID == "" || transaction.Description == "" || transaction.Amount == 0 || transaction.Currency == "" {
		return errors.New("missing required fields")
	}

	fmt.Println("Transaction ID before inserting:", transaction.ID.Hex())
	fmt.Println(transaction)
	err := service.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (service *TransactionService) GetAllTransactions(token string) ([]models.Transaction, error) {

	userRole, err := service.CheckUserRoleFromToken(token)
	if err != nil {
		return nil, err
	}

	if userRole != "admin" {
		return nil, errors.New("only admin users can perform this action")
	}


	transactions, err := service.transactionRepo.GetAllTransactions()
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (service *TransactionService) CheckUserRoleFromToken(tokenString string) (string, error) {
    if tokenString == "" {
        return "", errors.New("token is empty")
    }

    tokenParts := strings.Split(tokenString, ".")
    if len(tokenParts) != 3 {
        return "", errors.New("invalid token format")
    }

    payload, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
    if err != nil {
        return "", errors.New("error decoding token payload")
    }

    var payloadData map[string]interface{}
    if err := json.Unmarshal(payload, &payloadData); err != nil {
        return "", errors.New("error parsing token payload")
    }

    role, ok := payloadData["role"].(string)
    if !ok {
        return "", errors.New("role not found in token payload")
    }

    return role, nil
}

func (service *TransactionService) GetMerchantIDFromToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token is empty")
	}

	tokenParts := strings.Split(tokenString, ".")
	if len(tokenParts) != 3 {
		return "", errors.New("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return "", errors.New("error decoding token payload")
	}

	var payloadData map[string]interface{}
	if err := json.Unmarshal(payload, &payloadData); err != nil {
		return "", errors.New("error parsing token payload")
	}

	userID, ok := payloadData["userId"].(string)
	if !ok {
		return "", errors.New("userID not found in token payload")
	}

	return userID, nil
}