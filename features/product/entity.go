package product

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Product struct {
	ID        uint   `json:"id"`
	Category  string `json:"category"`
	Name      string `json:"name"`
	CPU       string `json:"cpu"`
	RAM       string `json:"ram"`
	Display   string `json:"display"`
	Storage   string `json:"storage"`
	Thickness string `json:"thickness"`
	Weight    string `json:"weight"`
	Bluetooth string `json:"bluetooth"`
	HDMI      string `json:"hdmi"`
	Price     int    `json:"price"`
	Picture   string `json:"picture"`
}

type Handler interface {
	Add() echo.HandlerFunc
	GetAll() echo.HandlerFunc
	GetProductDetail() echo.HandlerFunc
	SearchAll() echo.HandlerFunc
	DelProduct() echo.HandlerFunc
	UpdateProduct() echo.HandlerFunc
}
type Service interface {
	TalkToGpt(token *jwt.Token, newProduct Product) (Product, error)
	SemuaProduct(page, limit int) ([]Product, int, error)
	SatuProduct(productID uint) (Product, error)
	CariProduct(name string, category string, minPrice string, maxPrice string, page int, limit int) ([]Product, int, error)
	UpdateProduct(token *jwt.Token, productID uint, input Product) (Product, error)
	DelProduct(token *jwt.Token, productID uint) error
}

type Repository interface {
	InsertProduct(UserID uint, newProduct Product) (Product, error)
	GetAllProduct(page, limit int) ([]Product, int, error)
	GetProductID(productID uint) (*Product, error)
	UpdateProduct(UserID uint, productID uint, input Product) (Product, error)
	SearchProduct(name string, category string, minPrice uint, maxPrice uint, page int, limit int) ([]Product, int, error)
	DelProduct(UserID uint, productID uint) error
}
