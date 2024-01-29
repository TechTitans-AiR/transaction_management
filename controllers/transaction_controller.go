package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

var banks []models.Bank

func (controller *TransactionController) CreateTransactionWithCardHandler(w http.ResponseWriter, r *http.Request) {
	if err := loadBanksData(); err != nil {
		fmt.Println("Error loading banks data:", err)
	}
	var transaction models.TransactionWithCard

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&transaction); err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	cardPayload := models.Card{
		CardNumber:     transaction.Card.CardNumber,
		ExpirationDate: transaction.Card.ExpirationDate,
		Balance:        transaction.Card.Balance,
		CVC:            transaction.Card.CVC,
	}

	cardPayloadBytes, err := json.Marshal(cardPayload)
	if err != nil {
		fmt.Println("Error encoding card payload:", err)
		http.Error(w, "Error encoding card payload", http.StatusInternalServerError)
		return
	}

	bankEndpoint, err := chooseBankEndpoint(transaction.Card.CardNumber)
	if err != nil {
		fmt.Println("Error choosing bank endpoint:", err)
		http.Error(w, "Bank with provided card number is not recognized", http.StatusInternalServerError)
		return
	}

	fmt.Println("ODABRAO BANKU", bankEndpoint)
	//DOCKER
	bankHostResponse, err := http.Post(bankEndpoint, "application/json", bytes.NewBuffer(cardPayloadBytes))
	//LOCALHOST
	//bankHostResponse, err := http.Post("http://localhost:8084/sendToBankHostSimulator", "application/json", bytes.NewBuffer(cardPayloadBytes))
	if err != nil {
		fmt.Println("Error sending request to bank host:", err)
		http.Error(w, "Error sending request to bank host", http.StatusInternalServerError)
		return
	}
	defer bankHostResponse.Body.Close()

	var bankHostResponseBody map[string]string
	if err := json.NewDecoder(bankHostResponse.Body).Decode(&bankHostResponseBody); err != nil {
		fmt.Println("Error decoding bank host response:", err)
		http.Error(w, "Error decoding bank host response", http.StatusInternalServerError)
		return
	}

	status, ok := bankHostResponseBody["status"]
	if !ok {
		fmt.Println("Status not found in bank host response")
		http.Error(w, "Status not found in bank host response", http.StatusInternalServerError)
		return
	}
	fmt.Println("STATUS:", status)
	if status == "declined" {
		http.Error(w, "Payment declined", http.StatusBadRequest)
		return
	} else {
		currentTime := time.Now()
		transaction.CreatedAt = currentTime
		transaction.UpdatedAt = currentTime

		transaction.ID = primitive.NilObjectID
		fmt.Println("transaction u transaction_controller->", &transaction, "%n")
		fmt.Println("transaction u transaction_controller2->", transaction, "%n")
		err = controller.transactionService.CreateTransactionWithCard(&transaction)
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

}

func chooseBankEndpoint(cardNumber string) (string, error) {
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")

	firstSixDigits, err := strconv.Atoi(cardNumber[:6])
	if err != nil {
		return "", fmt.Errorf("Error converting first six digits to integer: %v", err)
	}

	for _, bank := range banks {
		for _, mark := range bank.Mark {
			if mark == firstSixDigits {
				println("BANKA ->", bank.Name)
				return bank.Endpoint, nil
			}
		}
	}

	return "", fmt.Errorf("There is no designated bank for that card - %s", firstSixDigits)
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

func loadBanksData() error {
	file, err := ioutil.ReadFile("banks.json")
	if err != nil {
		return fmt.Errorf("Error reading banks.json: %v", err)
	}

	if err := json.Unmarshal(file, &banks); err != nil {
		return fmt.Errorf("Error unmarshaling banks.json: %v", err)
	}

	return nil
}
