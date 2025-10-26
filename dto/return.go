package dto

type ReturnProduct struct {
	ProductID string  `json:"productId"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"` // Damaged, Defective, Unwanted, Wrong Item
}

type ReturnDTO struct {
	ID                 string          `json:"id"`
	CustomerName       string          `json:"customerName"`
	ContactNumber      string          `json:"contactNumber"`
	OriginalBillNumber string          `json:"originalBillNumber,omitempty"`
	Products           []ReturnProduct `json:"products"`
	AdditionalNotes    string          `json:"additionalNotes,omitempty"`
	CreatedAt          string          `json:"createdAt"`
}
