package product

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Product struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CPU       string `json:"cpu"`
	RAM       string `json:"ram"`
	Display   string `json:"display"`
	Storage   string `json:"storage"`
	Thickness string `json:"thickness"`
	Weight    string `json:"weight"`
	Bluetooth string `json:"bluetooth"`
	HDMI      string `json:"hdmi"`
	Price     string `json:"price"`
	Picture   string `json:"picture"`
}

type Handler interface {
	Add() echo.HandlerFunc
}
type Service interface {
	TalkToGpt(token *jwt.Token, newProduct Product) (Product, error)
}

type Repository interface {
	InsertProduct(UserID uint, newProduct Product) (Product, error)
}
