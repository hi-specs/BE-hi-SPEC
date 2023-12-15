package transaction

import (
	"BE-hi-SPEC/features/product"

	"github.com/labstack/echo/v4"
)

type TransactionDashboard struct {
	TotalProduct     int
	TotalUser        int
	TotalTransaction int
	Product          []product.Product
}

type Handler interface {
	AdminDashboard() echo.HandlerFunc
}

type Repository interface {
	AdminDashboard() (TransactionDashboard, error)
}

type Service interface {
	AdminDashboard() (TransactionDashboard, error)
}
