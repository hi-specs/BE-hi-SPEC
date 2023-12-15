package handler

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/helper/responses"
	"BE-hi-SPEC/utils/cld"
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	golangjwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type ProductHandler struct {
	s      product.Service
	cl     *cloudinary.Cloudinary
	ct     context.Context
	folder string
}

func New(s product.Service, cld *cloudinary.Cloudinary, ctx context.Context, uploadparam string) product.Handler {
	return &ProductHandler{
		s:      s,
		cl:     cld,
		ct:     ctx,
		folder: uploadparam,
	}
}

// SearchProductByRangePrice implements product.Handler.
func (ph *ProductHandler) SearchProductByRangePrice() echo.HandlerFunc {
	return func(c echo.Context) error {
		minPrice := c.QueryParam("minPrice")
		maxPrice := c.QueryParam("maxPrice")

		minPriceUint, err := strconv.ParseUint(minPrice, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid minPrice"})
		}

		maxPriceUint, err := strconv.ParseUint(maxPrice, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid maxPrice"})
		}

		products, err := ph.s.CariProductPrice(uint(minPriceUint), uint(maxPriceUint))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		var response []SearchResponse
		for _, result := range products {
			response = append(response, SearchResponse{
				ID:      result.ID,
				Name:    result.Name,
				Price:   result.Price,
				Picture: result.Picture,
			})
		}

		return c.JSON(http.StatusOK, response)
	}
}

// SearchProductByCategory implements product.Handler.
func (ph *ProductHandler) SearchProductByCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		category := c.QueryParam("category")

		products, err := ph.s.CariProductCategory(category)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		var response []SearchResponse
		for _, result := range products {
			response = append(response, SearchResponse{
				ID:      result.ID,
				Name:    result.Name,
				Price:   result.Price,
				Picture: result.Picture,
			})
		}

		return c.JSON(http.StatusOK, response)
	}
}

// SearchProductByName implements product.Handler.
func (ph *ProductHandler) SearchProductByName() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.QueryParam("name")

		products, err := ph.s.CariProduct(name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		var response []SearchResponse
		for _, result := range products {
			response = append(response, SearchResponse{
				ID:      result.ID,
				Name:    result.Name,
				Price:   result.Price,
				Picture: result.Picture,
			})
		}

		return c.JSON(http.StatusOK, response)
	}
}

// GetProductDetail implements product.Handler.
func (ph *ProductHandler) GetProductDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
				"data":    nil,
			})
		}
		result, err := ph.s.SatuProduct(uint(productID))
		if err != nil {
			c.Logger().Error("Error fetching product: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve product data",
			})
		}
		var response = new(ProductResponse)
		response.ID = result.ID
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
		response.Picture = result.Picture
		response.Category = result.Category

		return responses.PrintResponse(c, http.StatusCreated, "success create data", response)
	}
}

// GetAll implements product.Handler.
func (ph *ProductHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, err := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 {
			limit = 5
		}
		results, err := ph.s.SemuaProduct(page, limit)
		if err != nil {
			c.Logger().Error("Error fetching product: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to retrieve product data",
			})
		}
		var response []AllResponse
		for _, result := range results {
			response = append(response, AllResponse{
				ID:       result.ID,
				Name:     result.Name,
				Category: result.Category,
				Price:    result.Price,
				Picture:  result.Picture,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    "Success fetching all Posts data",
			"data":       response,
			"pagination": map[string]interface{}{"page": page, "limit": limit},
		})
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
		formHeader, err := c.FormFile("picture")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				inputProcess := &product.Product{
					Name:     input.Laptop,
					Category: input.Category,
					Picture:  "",
				}

				result, err := ph.s.TalkToGpt(c.Get("user").(*golangjwt.Token), *inputProcess)

				if err != nil {
					c.Logger().Error("ERROR Register, explain:", err.Error())
					var statusCode = http.StatusInternalServerError
					var message = "terjadi permasalahan ketika memproses data"

					if strings.Contains(err.Error(), "terdaftar") {
						statusCode = http.StatusBadRequest
						message = "data yang diinputkan sudah terdaftar ada sistem"
					}

					return c.JSON(statusCode, map[string]any{
						"message": message,
					})
				}

				var response = new(ProductResponse)
				response.ID = result.ID
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
				response.Picture = result.Picture
				response.Category = result.Category

				return c.JSON(http.StatusCreated, map[string]any{
					"message": "success create data",
					"data":    response,
				})

			}
			return c.JSON(
				http.StatusInternalServerError, map[string]any{
					"message": "formheader error",
				})

		}
		formFile, err := formHeader.Open()
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError, map[string]any{
					"message": "formfile error",
				})
		}

		link, err := cld.UploadImage(ph.cl, ph.ct, formFile, ph.folder)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "harap pilih gambar",
					"data":    nil,
				})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]any{
					"message": "kesalahan pada server",
					"data":    nil,
				})
			}
		}

		var inputProcess = new(product.Product)
		inputProcess.Name = input.Laptop
		inputProcess.Category = input.Category
		inputProcess.Picture = link

		result, err := ph.s.TalkToGpt(c.Get("user").(*golangjwt.Token), *inputProcess)
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
		response.ID = result.ID
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
		response.Picture = result.Picture
		response.Category = result.Category

		return responses.PrintResponse(c, http.StatusCreated, "success create data", response)
	}

}
