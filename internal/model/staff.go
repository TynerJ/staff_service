package model

// Model for the staff records
type Staff struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Quote string `json:"quote"`
}
