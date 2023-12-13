package handler

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/helper/responses"
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
)

type ProductHandler struct {
	s product.Service
}

func New(s product.Service) product.Handler {
	return &ProductHandler{
		s: s,
	}
}

func (ph *ProductHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(ProductRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		var inputProcess = new(product.Product)
		inputProcess.Name = input.Laptop

		result, err := ph.s.TalkToGpt(*inputProcess)
		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "terdaftar") {
				statusCode = http.StatusBadRequest
				message = "data yang diinputkan sudah terdaftar ada sistem"
			}

			return responses.PrintResponse(c, statusCode, message, nil)
		}

		var response = new(ProductResponse)
		response.Name = result.Name
		response.CPU = result.CPU
		response.RAM = result.RAM
		response.Display = result.Display
		response.Storage = result.Storage
		response.Thickness = result.Thickness
		response.Weight = result.Weight
		response.Bluetooth = result.Bluetooth
		response.HDMI = result.HDMI
		response.Price = result.Price

		return responses.PrintResponse(c, http.StatusCreated, "success create data", response)
	}

}
