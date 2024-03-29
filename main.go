package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"transaction_management/controllers"
	"transaction_management/repositories"
	"transaction_management/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	mongoURI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	transactionRepo := repositories.NewTransactionRepository(client.Database("transaction_management"))
	transactionService := services.NewTransactionService(transactionRepo)
	transactionController := controllers.NewTransactionController(transactionService)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/transactions/merchant/{merchantID}", transactionController.GetTransactionsByMerchantIDHandler).Methods("GET")
	router.HandleFunc("/api/v1/transactions/search", transactionController.SearchTransactionsHandler).Methods("POST")
	router.HandleFunc("/api/v1/transactions/{id}", transactionController.GetTransactionByIDHandler).Methods("GET")
	router.HandleFunc("/api/v1/transactions/create", transactionController.CreateTransactionHandler).Methods("POST")
	router.HandleFunc("/api/v1/transactions/create/card", transactionController.CreateTransactionWithCardHandler).Methods("POST")
	router.HandleFunc("/api/v1/transactions", transactionController.GetAllTransactionsHandler).Methods("GET")

	port := ":8082"
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
