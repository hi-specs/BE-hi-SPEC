package handler

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/transaction"
	"BE-hi-SPEC/features/user/handler"
	"BE-hi-SPEC/features/user/repository"
	"BE-hi-SPEC/helper/responses"
	"BE-hi-SPEC/utils/cld"
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	golangjwt "github.com/golang-jwt/jwt/v5"
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
		page, err := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 {
			limit = 10
		}
		result, totalPage, err := th.s.AdminDashboard(c.Get("user").(*golangjwt.Token), page, limit)

		if err != nil {
			c.Logger().Error("ERROR Fetch Transaction Dashboard, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "admin role required") {
				statusCode = http.StatusUnauthorized
				message = "Anda tidak memiliki izin untuk mengakses halaman ini"
			}

			return c.JSON(statusCode, map[string]interface{}{
				"message": message,
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
			"message":    "Success fetching all data for transaction dashboard",
			"data":       response,
			"pagination": map[string]interface{}{"page": page, "limit": limit, "total_page": totalPage},
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

		result, err := th.s.Checkout(c.Get("user").(*golangjwt.Token), input.ProductID, input.TotalPrice)

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
		page, err := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 {
			limit = 10
		}
		result, totalPage, err := th.s.TransactionList(c.Get("user").(*golangjwt.Token), page, limit)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"message": err.Error(),
			})
		}

		// slicing data user
		var User []repository.UserModel
		var UserName []string
		var UserPicture []string
		for _, result := range result {
			User = append(User, result.Users...)
		}
		for _, result := range User {
			UserName = append(UserName, result.Name)
			UserPicture = append(UserPicture, result.Avatar)
		}

		// slicing data product
		var Product []product.Product
		var ProductName []string
		var ProductPicture []string

		for _, result := range result {
			Product = append(Product, result.Products...)
		}

		for _, result := range Product {
			ProductName = append(ProductName, result.Name)
			ProductPicture = append(ProductPicture, result.Picture)
		}

		// slicing data transaction
		var Transaction []TransactionList
		for _, result := range result {
			Transaction = append(Transaction, TransactionList{
				TransactionID: result.TransactionID,
				Nota:          result.Nota,
				ProductID:     result.ProductID,
				TotalPrice:    result.TotalPrice,
				Status:        result.Status,
				Timestamp:     result.Timestamp,
				Token:         result.Token,
				Url:           result.Url,
			})
		}

		// slicing data to response
		var responses []TransactionsResponse
		for x, result := range result {
			responses = append(responses, TransactionsResponse{
				UserPicture:    UserPicture[x],
				UserName:       UserName[x],
				NameProduct:    ProductName[x],
				PictureProduct: ProductPicture[x],
				Nota:           result.Nota,
				TotalPrice:     uint(result.TotalPrice),
				Timestamp:      result.Timestamp,
				Status:         result.Status,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":      "Get All Transaction Successful",
			"transactions": responses,
			"pagination":   map[string]interface{}{"page": page, "limit": limit, "total_page": totalPage},
		})

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
		result, err := th.s.GetTransaction(c.Get("user").(*golangjwt.Token), uint(transactionID))

		if err != nil {
			c.Logger().Error("ERROR Get User Transaction, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika melihat data ini"

			if strings.Contains(err.Error(), "tidak ditemukan") {
				statusCode = http.StatusNotFound
				message = "user tidak ditemukan"
			} else if strings.Contains(err.Error(), "admin role required") {
				statusCode = http.StatusForbidden
				message = "Anda tidak memiliki izin untuk melihat data ini"
			}

			return c.JSON(statusCode, map[string]interface{}{
				"message": message,
			})
		}

		return responses.PrintResponse(c, http.StatusOK, "Detail Of Transaction", result)

	}
}

func (th *TransactionHandler) MidtransCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(MidtransCallBack)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}
		result, err := th.s.MidtransCallback(input.OrderID)
		if err != nil {
			c.Logger().Error("Error fetching product: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve product data",
			})
		}
		var response TransactionList
		response.TransactionID = result.TransactionID
		response.Nota = result.Nota
		response.ProductID = result.ProductID
		response.TotalPrice = result.TotalPrice
		response.Status = result.Status
		response.Timestamp = result.Timestamp
		response.Token = result.Token
		response.Url = result.Url

		return responses.PrintResponse(c, http.StatusOK, "Detail Of Midtrans Callback", response)

	}
}

func (th *TransactionHandler) UserTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
			})
		}
		result, err := th.s.UserTransaction(c.Get("user").(*golangjwt.Token), uint(userID))

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "terdaftar") {
				statusCode = http.StatusBadRequest
				message = "data yang diinputkan sudah terdaftar ada sistem"
			} else if strings.Contains(err.Error(), "admin role required") {
				statusCode = http.StatusForbidden
				message = "Anda tidak memiliki izin untuk melihat data ini"
			}

			return c.JSON(statusCode, map[string]any{
				"message": message,
			})
		}

		// mengambil data nama product
		listProd := result.Product
		var products []string
		for _, prod := range listProd {
			products = append(products, prod.Name)
		}

		// slicing data transaksi
		var nota []UserNota
		for x, trans := range result.Transaction {
			nota = append(nota, UserNota{
				ID:         trans.ID,
				Nota:       trans.Nota,
				TotalPrice: trans.TotalPrice,
				Status:     trans.Status,
				Token:      trans.Token,
				Url:        trans.Url,
				Product:    products[x],
			},
			)
		}

		// slicing data user
		dataUser := result.User
		var user handler.GetUserResponse
		user.Address = dataUser.Address
		user.Avatar = dataUser.Avatar
		user.Email = dataUser.Email
		user.ID = dataUser.ID
		user.PhoneNumber = dataUser.PhoneNumber
		user.Time = dataUser.CreatedAt
		user.Name = dataUser.Name

		// slicing responses
		var responses UserTransactionResponse
		responses.User = user
		responses.Nota = nota

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Data Transaction By User",
			"data":    responses,
		})
	}
}

func (th *TransactionHandler) DownloadTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
				"data":    err,
			})
		}

		// go to service
		err2 := th.s.DownloadTransaction(c.Get("user").(*golangjwt.Token), uint(transactionID))
		if err2 != nil {

			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": err2.Error(),
			})
		}

		// open file
		pdf := "helper/gofpdf/invoice.pdf"
		File, _ := os.Open(pdf)
		defer File.Close()

		// upload to cloudinary
		link, err := cld.UploadImage(th.cl, th.ct, File, th.folder)

		// update database
		th.s.UpdatePdfTransaction(link, uint(transactionID))

		// c.Attachment("helper/gofpdf/invoice.pdf", "invoice.pdf")
		return c.JSON(http.StatusOK, map[string]any{
			"message": "Download Transaction Successful",
			"url":     link,
		})
	}
}
