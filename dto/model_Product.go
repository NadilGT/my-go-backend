package dto

import (
	"time"
)

type Product struct {
	ProductId     string     `bson:"productId" json:"productId"`
	Name          string     `bson:"name" json:"name"`
	Barcode       string     `bson:"barcode" json:"barcode"`
	CategoryID    string     `bson:"categoryId" json:"categoryId"`
	BrandID       string     `bson:"brandId" json:"brandId"`
	SubCategoryID string     `bson:"subcategoryId" json:"subcategoryId"`
	CostPrice     float64    `bson:"costPrice" json:"costPrice"`
	SellingPrice  float64    `bson:"sellingPrice" json:"sellingPrice"`
	StockQty      int        `bson:"stockQty" json:"stockQty"`
	ExpiryDate    *time.Time `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Deleted       bool       `bson:"deleted" json:"deleted"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `bson:"updated_at" json:"updated_at"`
}
