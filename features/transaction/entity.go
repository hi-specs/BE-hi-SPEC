package transaction

import (
	"BE-hi-SPEC/features/product"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	ID         int `json:"transaction_id"`
	ProductID  int `json:"product_id"`
	TotalPrice int `json:"total_price"`
}
type TransactionDashboard struct {
	TotalProduct     int
	TotalUser        int
	TotalTransaction int
	Product          []product.Product
}

type Handler interface {
	AdminDashboard() echo.HandlerFunc
	Checkout() echo.HandlerFunc
}

type Repository interface {
	AdminDashboard() (TransactionDashboard, error)
	Checkout(userID uint, ProductID int, ProductPrice int) (Transaction, error)
}

type Service interface {
	AdminDashboard() (TransactionDashboard, error)
	Checkout(token *jwt.Token, ProductID int, ProductPrice int) (Transaction, error)
}
