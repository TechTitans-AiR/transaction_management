package repositories

import (
	"context"
	"transaction_management/models"

	"fmt"
	"time"

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
func (repo *TransactionRepository) Search(merchantID, description string, createdAt time.Time) ([]models.TransactionWithCard, error) {
	filter := bson.M{}

	if merchantID != "" {
		filter["merchantId"] = merchantID
	}

	if description != "" {
		filter["description"] = primitive.Regex{Pattern: description, Options: "i"}
	}

	if !createdAt.IsZero() {
		filter["createdAt"] = bson.M{"$gte": createdAt, "$lt": createdAt.Add(24 * time.Hour)}
	}

	var transactions []models.TransactionWithCard

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

func (repo *TransactionRepository) GetByID(id string) (*models.TransactionWithCard, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var transaction models.TransactionWithCard

	err = repo.collection.FindOne(context.TODO(), filter).Decode(&transaction)
	if err != nil {

		return nil, err
	}

	return &transaction, nil
}

func (repo *TransactionRepository) GetByMerchantID(merchantID string) ([]models.TransactionWithCard, error) {
	filter := bson.M{"merchantId": merchantID}
	fmt.Printf("Filter: %v\n", filter)
	var transactions []models.TransactionWithCard
	fmt.Print("----TRANSACTION ----", transactions)
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

func (repo *TransactionRepository) CreateTransactionWithCard(transaction *models.TransactionWithCard) error {
	transaction.ID = primitive.NilObjectID

	card := bson.M{
		"cardNumber":     transaction.Card.CardNumber,
		"expirationDate": transaction.Card.ExpirationDate,
		"balance":        transaction.Card.Balance,
		"cvc":            transaction.Card.CVC,
	}

	document := bson.M{
		"merchantId":  transaction.MerchantID,
		"description": transaction.Description,
		"amount":      transaction.Amount,
		"currency":    transaction.Currency,
		"card":        card,
		"createdAt":   transaction.CreatedAt,
		"updatedAt":   transaction.UpdatedAt,
	}

	_, err := repo.collection.InsertOne(context.TODO(), document)
	if err != nil {
		return err
	}

	return nil
}

func (repo *TransactionRepository) GetAllTransactions() ([]models.TransactionWithCard, error) {
	var transactions []models.TransactionWithCard
	cursor, err := repo.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &transactions); err != nil {
		return nil, err
	}
	fmt.Print("****TRANSAKCIJE****", transactions)
	return transactions, nil
}
