package repositories

import (
	"context"
	"transaction_management/models"

	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database) *TransactionRepository {
	collection := db.Collection("transactions")
	return &TransactionRepository{collection: collection}
}

func (repo *TransactionRepository) GetByMerchantID(merchantID string) ([]models.Transaction, error) {
	filter := bson.M{"merchantId": merchantID}
	fmt.Printf("Filter: %v\n", filter)
	var transactions []models.Transaction
	cursor, err := repo.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (repo *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	transaction.ID = primitive.NilObjectID

	_, err := repo.collection.InsertOne(context.TODO(), transaction)
	if err != nil {
		return err
	}

	return nil
}

func (repo *TransactionRepository) GetAllTransactions() ([]models.Transaction, error) {
    var transactions []models.Transaction
    cursor, err := repo.collection.Find(context.TODO(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())
    if err := cursor.All(context.TODO(), &transactions); err != nil {
        return nil, err
    }
    return transactions, nil
}
