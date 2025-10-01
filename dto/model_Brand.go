package dto

import (
	"time"
)

type Brand struct {
	BrandId    string    `bson:"brandId" json:"brandId"`
	Name       string    `bson:"name" json:"name"`
	CategoryID string    `bson:"categoryId" json:"categoryId"`
	Deleted    bool      `bson:"deleted" json:"deleted"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
