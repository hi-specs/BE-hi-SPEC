package transaction

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/features/user/repository"
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
	Users         []repository.UserModel
	Products      []product.Product
}

type UserTransaction struct {
	User        user.User
	Product     []product.Product
	Transaction []Transaction
}

type Handler interface {
	AdminDashboard() echo.HandlerFunc
	Checkout() echo.HandlerFunc
	TransactionList() echo.HandlerFunc
	GetTransaction() echo.HandlerFunc
	MidtransCallback() echo.HandlerFunc
	UserTransaction() echo.HandlerFunc
	DownloadTransaction() echo.HandlerFunc
}

type Repository interface {
	AdminDashboard(userID uint, page int, limit int) (TransactionDashboard, int, error)
	Checkout(userID uint, ProductID int, ProductPrice int) (Transaction, error)
	TransactionList(userId uint, page, limit int) ([]TransactionList, int, error)
	GetTransaction(userID uint, transactionID uint) (*TransactionList, error)
	MidtransCallback(transactionID string) (*TransactionList, error)
	UserTransaction(userId int, userID uint) (UserTransaction, error)
	DownloadTransaction(userID, transactionID uint) error
	UpdatePdfTransaction(link string, transactionID uint) error
}

type Service interface {
	AdminDashboard(token *jwt.Token, page int, limit int) (TransactionDashboard, int, error)
	Checkout(token *jwt.Token, ProductID int, ProductPrice int) (Transaction, error)
	TransactionList(token *jwt.Token, page, limit int) ([]TransactionList, int, error)
	GetTransaction(token *jwt.Token, transactionID uint) (TransactionList, error)
	MidtransCallback(transactionID string) (TransactionList, error)
	UserTransaction(token *jwt.Token, userID uint) (UserTransaction, error)
	DownloadTransaction(token *jwt.Token, transactionID uint) error
	UpdatePdfTransaction(link string, transactionID uint) error
}
