package handler

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

type TransactionDashboard struct {
	TotalProduct     int `json:"total_product"`
	TotalUser        int `json:"total_user"`
	TotalTransaction int `json:"total_transaction"`
	Product          []AllProduct
}
