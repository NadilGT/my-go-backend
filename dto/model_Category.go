package dto

import (
	"time"
)

type Category struct {
	CategoryId string    `bson:"categoryId" json:"categoryId"`
	Name       string    `bson:"name" json:"name"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}
