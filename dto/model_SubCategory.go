package dto

import (
	"time"
)

type SubCategory struct {
	SubCategoryId string    `bson:"subCategoryId" json:"subCategoryId"`
	Name          string    `bson:"name" json:"name"`
	BrandID       string    `bson:"brandId" json:"brandId"`
	Deleted       bool      `bson:"deleted" json:"deleted"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}
