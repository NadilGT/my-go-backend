package dto

import (
	"time"
)

type Batch struct {
	BatchId      string     `bson:"batchId" json:"batchId"`
	StockQty     int        `bson:"stockQty" json:"stockQty"`
	ExpiryDate   *time.Time `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	CostPrice    float64    `bson:"costPrice" json:"costPrice"`
	SellingPrice float64    `bson:"sellingPrice" json:"sellingPrice"`
	CreatedAt    time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `bson:"updated_at" json:"updated_at"`
}
