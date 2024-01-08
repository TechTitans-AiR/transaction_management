package controllers

import (
	"encoding/json"
	"fmt"
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
	vars := mux.Vars(r)
	id := vars["id"]

	transaction, err := controller.transactionService.GetTransactionByID(id)
	if err != nil {

		http.Error(w, "Error fetching transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		fmt.Printf("Transaction: %v\n", transaction)
		return
	}
}

func (controller *TransactionController) GetTransactionsByMerchantIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	merchantID := vars["merchantID"]

	transactions, err := controller.transactionService.GetTransactionsByMerchantID(merchantID)
	if err != nil {
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
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
	transactions, err := controller.transactionService.GetAllTransactions()
	if err != nil {
		http.Error(w, "Error fetching transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		fmt.Printf("Transactions: %v\n", transactions)
		return
	}
}
