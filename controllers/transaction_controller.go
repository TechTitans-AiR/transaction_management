package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"transaction_management/models"
	"transaction_management/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionController struct {
	transactionService *services.TransactionService
}

func NewTransactionController(transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}
func (controller *TransactionController) GetTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	vars := mux.Vars(r)
	id := vars["id"]

	transaction, err := controller.transactionService.GetTransactionByID(id, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		fmt.Printf("Transaction: %v\n", transaction)
		return
	}
}

func (controller *TransactionController) SearchTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		MerchantID  string `json:"merchantId"`
		Description string `json:"description"`
		CreatedAt   string `json:"createdAt"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil && err != io.EOF {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdAt, err := time.Parse("2006-01-02", requestBody.CreatedAt)
	if err != nil && requestBody.CreatedAt != "" {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	transactions, err := controller.transactionService.SearchTransactions(requestBody.MerchantID, requestBody.Description, createdAt)
	if err != nil {
		http.Error(w, "Error searching transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		fmt.Printf("Transactions: %v\n", transactions)
		return
	}
}
func (controller *TransactionController) GetTransactionsByMerchantIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	merchantID := vars["merchantID"]

	token := r.Header.Get("Authorization")

	transactions, err := controller.transactionService.GetTransactionsByMerchantID(token, merchantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		fmt.Printf("Transactions: %v\n", transactions)
		return
	}
}

func (controller *TransactionController) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	currentTime := time.Now()
	transaction.CreatedAt = currentTime
	transaction.UpdatedAt = currentTime

	transaction.ID = primitive.NilObjectID
	fmt.Println("transaction u transaction_controller->", &transaction)
	err = controller.transactionService.CreateTransaction(&transaction)
	if err != nil {
		http.Error(w, "Error creating transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (controller *TransactionController) GetAllTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	transactions, err := controller.transactionService.GetAllTransactions(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		fmt.Printf("Transactions: %v\n", transactions)
		return
	}
}
