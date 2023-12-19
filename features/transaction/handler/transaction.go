package handler

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/transaction"
	"BE-hi-SPEC/helper/responses"
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	gojwt "github.com/golang-jwt/jwt/v5"
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

func (th *TransactionHandler) Checkout() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(TransactionRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		result, err := th.s.Checkout(c.Get("user").(*gojwt.Token), input.ProductID, input.TotalPrice)

		if err != nil {
			c.Logger().Error("terjadi kesalahan", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "dobel input nama",
				})
			}
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "transaction duplicate",
			})
		}

		var response = new(TransactionResponse)
		response.ID = result.ID
		response.ProductID = result.ProductID
		response.TotalPrice = result.TotalPrice
		response.Status = result.Status
		response.Url = result.Url
		response.Token = result.Token
		response.Nota = result.Nota

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "Transaction created successfully",
			"data":    response,
		})
	}
}

func (th *TransactionHandler) TransactionList() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := th.s.TransactionList()
		if err != nil {
			c.Logger().Error("Error fetching transaction: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve product data",
			})
		}
		return responses.PrintResponse(c, http.StatusOK, "list of transaction", result)
	}
}

func (th *TransactionHandler) GetTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
			})
		}
		result, err := th.s.GetTransaction(uint(transactionID))

		if err != nil {
			c.Logger().Error("Error fetching product: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve product data",
			})
		}

		return responses.PrintResponse(c, http.StatusOK, "Detail Of Transaction", result)

	}
}

func (th *TransactionHandler) MidtransCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		transactionID := c.QueryParam("order_id")

		result, err := th.s.MidtransCallback(transactionID)

		if err != nil {
			c.Logger().Error("Error fetching product: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve product data",
			})
		}

		return responses.PrintResponse(c, http.StatusOK, "Detail Of Midtrans Callback", result)

	}
}
