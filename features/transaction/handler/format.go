package handler

import "time"

type Total struct {
	TotalProduct     int `json:"total_product"`
	TotalUser        int `json:"total_user"`
	TotalTransaction int `json:"total_transaction"`
}

type AllProduct struct {
	ID       uint   `json:"id"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Picture  string `json:"picture"`
}

type AdminDashboard struct {
	TotalProduct     int          `json:"total_product"`
	TotalUser        int          `json:"total_user"`
	TotalTransaction int          `json:"total_transaction"`
	Product          []AllProduct `json:"product"`
}

type TransactionRequest struct {
	ProductID  int `json:"product_id"`
	TotalPrice int `json:"total_price"`
}

type TransactionResponse struct {
	ID         int    `json:"transaction_id"`
	ProductID  int    `json:"product_id"`
	TotalPrice int    `json:"total_price"`
	Status     string `json:"status"`
}

type TransactionDetail struct {
	TransactionID int       `json:"transaction_id"`
	ProductID     int       `json:"product_id"`
	TotalPrice    int       `json:"total_price"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
}
