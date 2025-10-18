package dto

import (
	"time"
)

// SaleReport represents a daily sales report snapshot
// This data persists for 30 days while actual sales are deleted after 24 hours
type SaleReport struct {
	ReportId   string    `bson:"reportId" json:"reportId"`
	ReportDate time.Time `bson:"reportDate" json:"reportDate"` // The day this report represents

	// Aggregate Data
	TotalSales     int     `bson:"totalSales" json:"totalSales"`         // Number of sales
	TotalRevenue   float64 `bson:"totalRevenue" json:"totalRevenue"`     // Total revenue
	TotalDiscount  float64 `bson:"totalDiscount" json:"totalDiscount"`   // Total discounts given
	TotalTax       float64 `bson:"totalTax" json:"totalTax"`             // Total tax collected
	TotalItemsSold int     `bson:"totalItemsSold" json:"totalItemsSold"` // Total quantity of items

	// Payment Method Breakdown
	CashSales       int     `bson:"cashSales" json:"cashSales"`
	CardSales       int     `bson:"cardSales" json:"cardSales"`
	TransferSales   int     `bson:"transferSales" json:"transferSales"`
	CashRevenue     float64 `bson:"cashRevenue" json:"cashRevenue"`
	CardRevenue     float64 `bson:"cardRevenue" json:"cardRevenue"`
	TransferRevenue float64 `bson:"transferRevenue" json:"transferRevenue"`

	// Sales with customer info
	SalesWithCustomer    int `bson:"salesWithCustomer" json:"salesWithCustomer"`
	SalesWithoutCustomer int `bson:"salesWithoutCustomer" json:"salesWithoutCustomer"`

	// Top Products (optional - can store top 10 products sold)
	TopProducts []TopProductReport `bson:"topProducts,omitempty" json:"topProducts,omitempty"`

	// Metadata
	GeneratedAt time.Time `bson:"generatedAt" json:"generatedAt"`
	GeneratedBy string    `bson:"generatedBy" json:"generatedBy"` // System or user who generated
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type TopProductReport struct {
	ProductId    string  `bson:"productId" json:"productId"`
	ProductName  string  `bson:"productName" json:"productName"`
	QuantitySold int     `bson:"quantitySold" json:"quantitySold"`
	Revenue      float64 `bson:"revenue" json:"revenue"`
}
