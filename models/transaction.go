package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MerchantID  string             `json:"merchantId,omitempty" bson:"merchantId,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Amount      float64            `json:"amount,omitempty" bson:"amount,omitempty"`
	Currency    string             `json:"currency,omitempty" bson:"currency,omitempty"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
