package handler

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/transaction"
	"context"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	p      product.Handler
	s      transaction.Service
	cl     *cloudinary.Cloudinary
	ct     context.Context
	folder string
}

func New(s transaction.Service, cld *cloudinary.Cloudinary, ctx context.Context, uploadparam string) transaction.Handler {
	return &TransactionHandler{
		s:      s,
		cl:     cld,
		ct:     ctx,
		folder: uploadparam,
	}
}

func (th *TransactionHandler) AdminDashboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := th.s.AdminDashboard()

		if err != nil {
			c.Logger().Error("Error fetching product: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve product data",
			})
		}

		var response AdminDashboard
		response.TotalProduct = result.TotalProduct
		response.TotalTransaction = result.TotalTransaction
		response.TotalUser = result.TotalUser

		var responses []AllProduct
		for _, result2 := range result.Product {
			responses = append(responses, AllProduct{
				ID:       result2.ID,
				Name:     result2.Name,
				Price:    result2.Price,
				Picture:  result2.Picture,
				Category: result2.Category,
			})
		}
		response.Product = responses

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Success fetching all data for transaction dashboard",
			"data":    response,
		})
	}
}
