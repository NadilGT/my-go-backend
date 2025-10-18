package dto

import (
	"time"
)

type SaleItem struct {
	ProductId   string  `bson:"productId" json:"productId" validate:"required"`
	ProductName string  `bson:"productName" json:"productName"`
	Quantity    int     `bson:"quantity" json:"quantity" validate:"required,min=1"`
	UnitPrice   float64 `bson:"unitPrice" json:"unitPrice" validate:"required,min=0"`
	TotalPrice  float64 `bson:"totalPrice" json:"totalPrice"`
	Discount    float64 `bson:"discount,omitempty" json:"discount,omitempty"` // Optional item-level discount
}

type Sale struct {
	SaleId     string `bson:"saleId" json:"saleId"`
	SaleNumber string `bson:"saleNumber" json:"saleNumber" validate:"required"`

	// Optional Customer Information
	CustomerName  string `bson:"customerName,omitempty" json:"customerName,omitempty"`
	CustomerPhone string `bson:"customerPhone,omitempty" json:"customerPhone,omitempty"`

	// Sale Details
	Items         []SaleItem `bson:"items" json:"items" validate:"required,min=1,dive"`
	SubTotal      float64    `bson:"subTotal" json:"subTotal"`
	TotalDiscount float64    `bson:"totalDiscount,omitempty" json:"totalDiscount,omitempty"` // Overall discount
	TaxAmount     float64    `bson:"taxAmount,omitempty" json:"taxAmount,omitempty"`         // Tax if applicable
	GrandTotal    float64    `bson:"grandTotal" json:"grandTotal"`

	// Payment Information
	PaymentMethod string  `bson:"paymentMethod" json:"paymentMethod" validate:"required,oneof=cash card transfer"`
	PaidAmount    float64 `bson:"paidAmount" json:"paidAmount" validate:"required,min=0"`
	ChangeAmount  float64 `bson:"changeAmount" json:"changeAmount"`

	// Metadata
	SaleDate  time.Time `bson:"saleDate" json:"saleDate" validate:"required"`
	SoldBy    string    `bson:"soldBy" json:"soldBy" validate:"required"` // Staff/cashier name
	Notes     string    `bson:"notes,omitempty" json:"notes,omitempty"`
	Status    string    `bson:"status" json:"status" validate:"required,oneof=completed refunded partial_refund"`
	Deleted   bool      `bson:"deleted" json:"deleted"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
