package transaction

import (
	"BE-hi-SPEC/features/product"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	ID         int    `json:"transaction_id"`
	Nota       string `json:"nota"`
	ProductID  int    `json:"product_id"`
	TotalPrice int    `json:"total_price"`
	Status     string `json:"status"`
	Token      string `json:"token"`
	Url        string `json:"url"`
}

type TransactionDashboard struct {
	TotalProduct     int
	TotalUser        int
	TotalTransaction int
	Product          []product.Product
}

type TransactionList struct {
	TransactionID int       `json:"transaction_id"`
	Nota          string    `json:"nota"`
	ProductID     int       `json:"product_id"`
	TotalPrice    int       `json:"total_price"`
	Status        string    `json:"status"`
	Timestamp     time.Time `json:"timestamp"`
	Token         string    `json:"token"`
	Url           string    `json:"url"`
}

type Handler interface {
	AdminDashboard() echo.HandlerFunc
	Checkout() echo.HandlerFunc
	TransactionList() echo.HandlerFunc
	GetTransaction() echo.HandlerFunc
	MidtransCallback() echo.HandlerFunc
}

type Repository interface {
	AdminDashboard() (TransactionDashboard, error)
	Checkout(userID uint, ProductID int, ProductPrice int) (Transaction, error)
	TransactionList(page, limit int) ([]TransactionList, int, error)
	GetTransaction(transactionID uint) (*TransactionList, error)
	MidtransCallback(transactionID string) (*TransactionList, error)
}

type Service interface {
	AdminDashboard() (TransactionDashboard, error)
	Checkout(token *jwt.Token, ProductID int, ProductPrice int) (Transaction, error)
	TransactionList(page, limit int) ([]TransactionList, int, error)
	GetTransaction(transactionID uint) (TransactionList, error)
	MidtransCallback(transactionID string) (TransactionList, error)
}
