package dto

import "time"

type SupplierProduct struct {
	SupplierID   string    `bson:"supplierId" json:"supplierId"`
	SupplierName string    `bson:"supplierName" json:"supplierName"`
	ProductID    string    `bson:"productId" json:"productId"`
	ProductName  string    `bson:"productName" json:"productName"`
	AssignedAt   time.Time `bson:"assignedAt" json:"assignedAt"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}
