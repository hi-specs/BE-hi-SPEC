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

// UpdateProduct implements product.Handler.
func (ph *ProductHandler) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(PutProductRequest)
		productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
				"data":    nil,
			})
		}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
				"data":    nil,
			})
		}

		formHeader, err := c.FormFile("picture")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": "ID user tidak valid",
					})
				}
				var inputProcess = new(product.Product)
				inputProcess.Picture = ""
				inputProcess.ID = uint(productID)
				inputProcess.Category = input.Category
				inputProcess.Name = input.Name
				inputProcess.CPU = input.CPU
				inputProcess.RAM = input.RAM
				inputProcess.Display = input.Display
				inputProcess.Storage = input.Storage
				inputProcess.Thickness = input.Thickness
				inputProcess.Weight = input.Weight
				inputProcess.Bluetooth = input.Bluetooth
				inputProcess.HDMI = input.HDMI
				inputProcess.Price = input.Price

				result, err := ph.s.UpdateProduct(uint(productID), *inputProcess)

				if err != nil {
					c.Logger().Error("ERROR Register, explain:", err.Error())
					var statusCode = http.StatusInternalServerError
					var message = "terjadi permasalahan ketika memproses data"

					if strings.Contains(err.Error(), "terdaftar") {
						statusCode = http.StatusBadRequest
						message = "data yang diinputkan sudah terdaftar ada sistem"
					}
					if strings.Contains(err.Error(), "yang lama") {
						statusCode = http.StatusBadRequest
						message = "harap masukkan password yang lama jika ingin mengganti password"
					}

					return c.JSON(statusCode, map[string]any{
						"message": message,
					})
				}

				var response = new(ProductResponse)
				response.ID = result.ID
				response.Category = result.Category
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

				return c.JSON(http.StatusCreated, map[string]any{
					"message": "success create data",
					"data":    response,
				})
			}
			return c.JSON(
				http.StatusBadRequest, map[string]any{
					"message": "formheader error",
					"data":    nil,
				})

		}

		formFile, err := formHeader.Open()
		if err != nil {
			return c.JSON(
				http.StatusBadRequest, map[string]any{
					"message": "formfile error",
					"data":    nil,
				})
		}
		defer formFile.Close()

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
		inputProcess.Picture = link
		inputProcess.ID = uint(productID)
		inputProcess.Category = input.Category
		inputProcess.Name = input.Name
		inputProcess.CPU = input.CPU
		inputProcess.RAM = input.RAM
		inputProcess.Display = input.Display
		inputProcess.Storage = input.Storage
		inputProcess.Thickness = input.Thickness
		inputProcess.Weight = input.Weight
		inputProcess.Bluetooth = input.Bluetooth
		inputProcess.HDMI = input.HDMI
		inputProcess.Price = input.Price

		result, err := ph.s.UpdateProduct(uint(productID), *inputProcess)

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika memproses data"

			if strings.Contains(err.Error(), "terdaftar") {
				statusCode = http.StatusBadRequest
				message = "data yang diinputkan sudah terdaftar ada sistem"
			}
			if strings.Contains(err.Error(), "yang lama") {
				statusCode = http.StatusBadRequest
				message = "harap masukkan password yang lama jika ingin mengganti password"
			}

			return c.JSON(statusCode, map[string]any{
				"message": message,
			})
		}

		var response = new(ProductResponse)
		response.ID = result.ID
		response.Category = result.Category
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

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "success updated data",
			"data":    response,
		})
	}
}

// SearchAll implements product.Handler.
func (ph *ProductHandler) SearchAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		page, err := strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit <= 0 {
			limit = 5
		}
		var minPrice int
		var maxPrice int
		name := c.QueryParam("name")
		category := c.QueryParam("category")

		minPrice, _ = strconv.Atoi(c.QueryParam("minprice"))
		maxPrice, _ = strconv.Atoi(c.QueryParam("maxprice"))

		products, totalPage, err := ph.s.CariProduct(name, category, uint(minPrice), uint(maxPrice), page, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		var response []SearchResponse
		for _, result := range products {
			response = append(response, SearchResponse{
				ID:       result.ID,
				Category: result.Category,
				Name:     result.Name,
				Price:    result.Price,
				Picture:  result.Picture,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":    "Success fetching all Search data",
			"data":       response,
			"pagination": map[string]interface{}{"page": page, "limit": limit, "total_page": totalPage},
		})
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

		return responses.PrintResponse(c, http.StatusOK, "Success Get Product Detail", response)
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
			limit = 10
		}
		results, totalPage, err := ph.s.SemuaProduct(page, limit)
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
			"pagination": map[string]interface{}{"page": page, "limit": limit, "total_page": totalPage},
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
					"message": "Success Create Product Data",
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

		return responses.PrintResponse(c, http.StatusCreated, "Success Create Product Data", response)
	}

}

func (ph *ProductHandler) DelProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID product tidak valid",
				"data":    nil,
			})
		}

		errDel := ph.s.DelProduct(uint(productID))

		if errDel != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "product tidak ditemukan",
			})
		}
		return c.JSON(http.StatusOK, map[string]any{
			"message": "delete product successful",
		})
	}
}
