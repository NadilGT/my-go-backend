package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductId  string             `bson:"productId" json:"productId"`
	Name       string             `bson:"name" json:"name"`
	StockQty   int                `bson:"stockQty" json:"stockQty"`
	ExpiryDate *time.Time         `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
