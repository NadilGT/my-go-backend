package dto

import "time"

type SaleItem struct {
	ProductID   string  `bson:"productId" json:"productId"`
	ProductName string  `bson:"productName" json:"productName"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	UnitPrice   float64 `bson:"unitPrice" json:"unitPrice"`
	TotalPrice  float64 `bson:"totalPrice" json:"totalPrice"`
}

type Sale struct {
	SaleID         string     `bson:"saleId" json:"saleId"`
	CustomerName   string     `bson:"customerName,omitempty" json:"customerName,omitempty"`
	MobileNumber   string     `bson:"mobileNumber,omitempty" json:"mobileNumber,omitempty"`
	Items          []SaleItem `bson:"items" json:"items"`
	Subtotal       float64    `bson:"subtotal" json:"subtotal"`
	Tax            float64    `bson:"tax" json:"tax"`
	TaxPercentage  float64    `bson:"taxPercentage" json:"taxPercentage"`
	Discount       float64    `bson:"discount" json:"discount"`
	DiscountType   string     `bson:"discountType" json:"discountType"` // "percentage" or "fixed"
	Total          float64    `bson:"total" json:"total"`
	PaymentMethod  string     `bson:"paymentMethod" json:"paymentMethod"` // "cash" or "card"
	AmountReceived float64    `bson:"amountReceived,omitempty" json:"amountReceived,omitempty"`
	Change         float64    `bson:"change,omitempty" json:"change,omitempty"`
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `bson:"updated_at" json:"updated_at"`
}

// Request DTOs
type CreateSaleRequest struct {
	CustomerName   string     `json:"customerName,omitempty"`
	MobileNumber   string     `json:"mobileNumber,omitempty"`
	Items          []SaleItem `json:"items" binding:"required"`
	Tax            float64    `json:"tax"`
	TaxPercentage  float64    `json:"taxPercentage"`
	Discount       float64    `json:"discount"`
	DiscountType   string     `json:"discountType"`                     // "percentage" or "fixed"
	PaymentMethod  string     `json:"paymentMethod" binding:"required"` // "cash" or "card"
	AmountReceived float64    `json:"amountReceived,omitempty"`
}

type CalculateOrderSummaryRequest struct {
	Items         []SaleItem `json:"items" binding:"required"`
	TaxPercentage float64    `json:"taxPercentage"`
	Discount      float64    `json:"discount"`
	DiscountType  string     `json:"discountType"` // "percentage" or "fixed"
}

type OrderSummaryResponse struct {
	Subtotal float64 `json:"subtotal"`
	Tax      float64 `json:"tax"`
	Discount float64 `json:"discount"`
	Total    float64 `json:"total"`
}
