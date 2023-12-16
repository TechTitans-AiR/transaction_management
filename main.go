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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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

	router.HandleFunc("/api/v1/transactions/{merchantID}", transactionController.GetTransactionsByMerchantIDHandler).Methods("GET")
	router.HandleFunc("/api/v1/transactions/create", transactionController.CreateTransactionHandler).Methods("POST")

	addr := ":" + port
	fmt.Printf("Listening on %s...\n", addr)
	http.ListenAndServe(addr, nil)
}
