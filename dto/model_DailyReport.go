package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DailyReportDocument represents a saved daily sales report in the database
type DailyReportDocument struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	ReportDate      time.Time            `bson:"reportDate" json:"reportDate"`
	Month           int                  `bson:"month" json:"month"` // Month number (1-12)
	Year            int                  `bson:"year" json:"year"`   // Year
	TotalSales      int                  `bson:"totalSales" json:"totalSales"`
	TotalRevenue    float64              `bson:"totalRevenue" json:"totalRevenue"`
	TotalDiscount   float64              `bson:"totalDiscount" json:"totalDiscount"`
	TotalTax        float64              `bson:"totalTax" json:"totalTax"`
	CashSales       int                  `bson:"cashSales" json:"cashSales"`
	CardSales       int                  `bson:"cardSales" json:"cardSales"`
	CashRevenue     float64              `bson:"cashRevenue" json:"cashRevenue"`
	CardRevenue     float64              `bson:"cardRevenue" json:"cardRevenue"`
	ProductsSold    []ProductSoldSummary `bson:"productsSold" json:"productsSold"`
	TopSellingItems []ProductSoldSummary `bson:"topSellingItems" json:"topSellingItems"`
	CreatedAt       time.Time            `bson:"createdAt" json:"createdAt"`
	ExpiresAt       time.Time            `bson:"expiresAt" json:"expiresAt"` // TTL for auto-deletion
}
