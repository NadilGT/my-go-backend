package dto

import (
	"time"
)

type Supplier struct {
	SupplierId string    `bson:"supplierId" json:"supplierId"`
	Name       string    `bson:"name" json:"name"`
	Contact    string    `bson:"contact" json:"contact"`
	Email      string    `bson:"email" json:"email"`
	Address    string    `bson:"address" json:"address"`
	Status     string    `bson:"status" json:"status"` // "active" or "inactive"
	Deleted    bool      `json:"deleted" bson:"deleted"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
