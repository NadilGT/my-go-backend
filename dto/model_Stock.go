package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ProductId  string             `bson:"productId" json:"productId"`
	BatchId    string             `bson:"batchId,omitempty" json:"batchId,omitempty"`
	Name       string             `bson:"name" json:"name"`
	StockQty   int                `bson:"stockQty" json:"stockQty"`
	Status     string             `bson:"-" json:"status"` // Not stored in DB, calculated dynamically
	ExpiryDate *time.Time         `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// CalculateStatus calculates the stock status based on quantity
// - Low Stock: quantity < 10
// - Average Stock: quantity >= 10 and < 25
// - Good Stock: quantity >= 25
func (s *Stock) CalculateStatus() {
	if s.StockQty < 10 {
		s.Status = "Low Stock"
	} else if s.StockQty >= 10 && s.StockQty < 25 {
		s.Status = "Average Stock"
	} else {
		s.Status = "Good Stock"
	}
}
