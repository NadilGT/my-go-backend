package dto

import (
	"time"
)

type GRNItem struct {
	ProductId   string     `bson:"productId" json:"productId" validate:"required"`
	ProductName string     `bson:"productName" json:"productName"`
	ExpectedQty int        `bson:"expectedQty" json:"expectedQty" validate:"required,min=1"`
	ReceivedQty int        `bson:"receivedQty" json:"receivedQty" validate:"required,min=0"`
	UnitCost    float64    `bson:"unitCost" json:"unitCost" validate:"required,min=0"`
	TotalCost   float64    `bson:"totalCost" json:"totalCost"`
	ExpiryDate  *time.Time `bson:"expiryDate,omitempty" json:"expiryDate,omitempty"`
	BatchNumber string     `bson:"batchNumber,omitempty" json:"batchNumber,omitempty"`
	Remarks     string     `bson:"remarks,omitempty" json:"remarks,omitempty"`
}

type GRN struct {
	GRNId         string     `bson:"grnId" json:"grnId"`
	GRNNumber     string     `bson:"grnNumber" json:"grnNumber" validate:"required"`
	SupplierId    string     `bson:"supplierId" json:"supplierId" validate:"required"`
	SupplierName  string     `bson:"supplierName" json:"supplierName"`
	ReceivedDate  time.Time  `bson:"receivedDate" json:"receivedDate" validate:"required"`
	InvoiceNumber string     `bson:"invoiceNumber,omitempty" json:"invoiceNumber,omitempty"`
	InvoiceDate   *time.Time `bson:"invoiceDate,omitempty" json:"invoiceDate,omitempty"`
	Items         []GRNItem  `bson:"items" json:"items" validate:"required,min=1,dive"`
	TotalAmount   float64    `bson:"totalAmount" json:"totalAmount"`
	Status        string     `bson:"status" json:"status" validate:"required,oneof=pending completed partial_received"`
	ReceivedBy    string     `bson:"receivedBy" json:"receivedBy" validate:"required"`
	Notes         string     `bson:"notes" json:"notes" validate:"required"`
	Deleted       bool       `bson:"deleted" json:"deleted"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `bson:"updated_at" json:"updated_at"`
}
