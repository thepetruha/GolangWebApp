package model

type ToDo struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id,omitempty"`
	Text       string `json:"text"`
	Date       string `json:"date"`
}
