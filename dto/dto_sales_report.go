package dto

import "time"

// DailySalesSummary represents the summary of sales for a specific date
type DailySalesSummary struct {
	ReportDate      time.Time            `json:"reportDate"`
	TotalSales      int                  `json:"totalSales"`
	TotalRevenue    float64              `json:"totalRevenue"`
	TotalDiscount   float64              `json:"totalDiscount"`
	TotalTax        float64              `json:"totalTax"`
	CashSales       int                  `json:"cashSales"`
	CardSales       int                  `json:"cardSales"`
	CashRevenue     float64              `json:"cashRevenue"`
	CardRevenue     float64              `json:"cardRevenue"`
	ProductsSold    []ProductSoldSummary `json:"productsSold"`
	TopSellingItems []ProductSoldSummary `json:"topSellingItems"`
}

// ProductSoldSummary represents the summary of a product sold during the day
type ProductSoldSummary struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	TotalAmount float64 `json:"totalAmount"`
}
