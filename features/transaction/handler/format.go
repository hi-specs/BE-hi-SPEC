package handler

import (
	"BE-hi-SPEC/features/user/handler"
)

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
	ProductID  int `json:"product_id" form:"product_id"`
	TotalPrice int `json:"total_price" form:"total_price"`
}

type TransactionResponse struct {
	ID         int    `json:"transaction_id"`
	Nota       string `json:"nota"`
	ProductID  int    `json:"product_id"`
	TotalPrice int    `json:"total_price"`
	Status     string `json:"status"`
	Token      string `json:"token"`
	Url        string `json:"url"`
}

type MidtransCallBack struct {
	OrderID string `json:"order_id"`
}

type UserNota struct {
	ID         int    `json:"transaction_id"`
	Nota       string `json:"nota"`
	Product    string `json:"product_name"`
	TotalPrice int    `json:"total_price"`
	Status     string `json:"status"`
	Token      string `json:"token"`
	Url        string `json:"url"`
}

type UserTransactionResponse struct {
	User handler.GetUserResponse `json:"user" form:"user"`
	Nota []UserNota              `json:"nota" form:"nota"`
}
