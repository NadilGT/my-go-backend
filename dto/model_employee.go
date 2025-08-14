package dto

type Employee struct {
	ID        string `json:"id" bson:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Type      string `json:"type"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Deleted   bool   `json:"deleted" bson:"deleted"`
}
