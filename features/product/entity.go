package product

import (
	"github.com/labstack/echo/v4"
)

type Product struct {
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
}

type Handler interface {
	Add() echo.HandlerFunc
}
type Service interface {
	TalkToGpt(newProduct Product) (Product, error)
}

type Repository interface {
	InsertProduct(newProduct Product) (Product, error)
}
