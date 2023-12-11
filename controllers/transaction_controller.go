package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"transaction_management/services"

	"github.com/gorilla/mux"
)

type TransactionController struct {
	transactionService *services.TransactionService
}

func NewTransactionController(transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
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
