package handler

import (
	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/helper/jwt"
	cld "BE-hi-SPEC/utils/cld"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	gojwt "github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type UserController struct {
	srv    user.Service
	cl     *cloudinary.Cloudinary
	ct     context.Context
	folder string
}

func New(s user.Service, cld *cloudinary.Cloudinary, ctx context.Context, uploadparam string) user.Handler {
	return &UserController{
		srv:    s,
		cl:     cld,
		ct:     ctx,
		folder: uploadparam,
	}
}

// Login implements user.Handler.
func (uc *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(LoginRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		result, err := uc.srv.Login(input.Email, input.Password)

		if err != nil {
			c.Logger().Error("ERROR Login, explain:", err.Error())
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]any{
					"message": "email tidak ditemukan",
				})
			}
			if strings.Contains(err.Error(), "password salah") {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "password salah",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "email tidak di temukan",
			})
		}

		strToken, err := jwt.GenerateJWT(result.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika mengenkripsi data",
			})
		}

		var response = new(LoginResponse)
		response.ID = result.ID
		response.Email = result.Email
		response.Password = result.Password
		response.Token = strToken

		return c.JSON(http.StatusOK, map[string]any{
			"message": "success login data",
			"data":    response,
		})
	}
}

// Register implements user.Handler.
func (uc *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(RegisterRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input tidak sesuai",
			})
		}
		var inputProses = new(user.User)
		inputProses.Email = input.Email
		inputProses.Name = input.Name
		inputProses.Address = input.Address
		inputProses.PhoneNumber = input.PhoneNumber
		inputProses.Password = input.Password

		result, err := uc.srv.Register(*inputProses)
		if err != nil {
			c.Logger().Error("terjadi kesalahan", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "dobel input nama",
				})
			}
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "email telah terdaftar",
			})
		}
		var response = new(RegisterResponse)
		response.ID = result.ID
		response.Name = result.Name
		response.Email = result.Email
		response.Address = result.Address
		response.PhoneNumber = result.PhoneNumber

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "success",
			"data":    response,
		})
	}
}

// Update implements user.Handler.
func (uc *UserController) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(PutRequest)
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
				"data":    nil,
			})
		}
		if userID == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Harap Login dulu",
				"data":    nil,
			})
		}
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
				"data":    nil,
			})
		}

		formHeader, err := c.FormFile("avatar")
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"message": "ID user tidak valid",
					})
				}
				var inputProcess = new(user.User)
				inputProcess.Avatar = ""
				inputProcess.ID = uint(userID)
				inputProcess.Address = input.Address
				inputProcess.Password = input.Password
				inputProcess.NewPassword = input.NewPassword
				inputProcess.PhoneNumber = input.PhoneNumber
				inputProcess.Email = input.Email
				inputProcess.Name = input.Name

				result, err := uc.srv.UpdateUser(c.Get("user").(*gojwt.Token), *inputProcess)

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

				var response = new(PutResponse)
				response.ID = result.ID
				response.Name = result.Name
				response.Email = result.Email
				response.PhoneNumber = result.PhoneNumber
				response.Address = result.Address
				response.Avatar = result.Avatar

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

		link, err := cld.UploadImage(uc.cl, uc.ct, formFile, uc.folder)
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

		var inputProcess = new(user.User)
		inputProcess.Avatar = link
		inputProcess.ID = uint(userID)
		inputProcess.Address = input.Address
		inputProcess.Password = input.Password
		inputProcess.NewPassword = input.NewPassword
		inputProcess.PhoneNumber = input.PhoneNumber
		inputProcess.Email = input.Email
		inputProcess.Name = input.Name

		result, err := uc.srv.UpdateUser(c.Get("user").(*gojwt.Token), *inputProcess)

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

		var response = new(PutResponse)
		response.ID = result.ID
		response.Name = result.Name
		response.Email = result.Email
		response.PhoneNumber = result.PhoneNumber
		response.Address = result.Address
		response.Avatar = result.Avatar

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "success create data",
			"data":    response,
		})
	}
}

// Delete implements user.Handler.
func (uc *UserController) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
			})
		}

		err = uc.srv.HapusUser(c.Get("user").(*gojwt.Token), uint(userID))
		if err != nil {
			c.Logger().Error("ERROR Delete User, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika menghapus user"

			if strings.Contains(err.Error(), "tidak ditemukan") {
				statusCode = http.StatusNotFound
				message = "user tidak ditemukan"
			} else if strings.Contains(err.Error(), "tidak memiliki izin") {
				statusCode = http.StatusForbidden
				message = "Anda tidak memiliki izin untuk menghapus user ini"
			}

			return c.JSON(statusCode, map[string]interface{}{
				"message": message,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success delete user",
		})
	}
}

func (uc *UserController) All() echo.HandlerFunc {
	return func(c echo.Context) error {
		AllUser, err := uc.srv.GetAllUser()
		if err != nil {
			c.Logger().Error("ERROR GetAll, explain:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response []GetUserResponse

		for _, result := range AllUser {
			responses := GetUserResponse{
				ID:          result.ID,
				Email:       result.Email,
				Name:        result.Name,
				Address:     result.Address,
				PhoneNumber: result.PhoneNumber,
				Avatar:      result.Avatar,
			}
			response = append(response, responses)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get all data",
			"data":    response,
		})
	}
}

func (uc *UserController) AddFavorite() echo.HandlerFunc {
	return func(c echo.Context) error {

		productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Product ID tidak valid",
			})
		}

		result, err := uc.srv.AddFavorite(c.Get("user").(*gojwt.Token), uint(productID))

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

		var responses GetAllFavoriteResponse
		responses.User.Email = result.User.Email
		responses.User.Name = result.User.Name
		responses.User.ID = result.User.ID
		responses.User.Avatar = result.User.Avatar

		var prod []GetAllFavoriteProduct
		for _, result := range result.Favorite {
			prod = append(prod, GetAllFavoriteProduct{
				ID:      result.ID,
				Name:    result.Name,
				Price:   result.Price,
				Picture: result.Picture,
			})
		}
		responses.Product = prod

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get all data",
			"data":    responses,
		})

	}
}

func (uc *UserController) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID user tidak valid",
			})
		}
		result, err := uc.srv.GetUser(uint(userID))

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

		var favID = result.FavID
		var responses GetAllFavoriteResponse
		responses.User.Email = result.User.Email
		responses.User.Name = result.User.Name
		responses.User.ID = result.User.ID
		responses.User.Avatar = result.User.Avatar

		var prod []GetAllFavoriteProduct
		for x, result := range result.Favorite {
			prod = append(prod, GetAllFavoriteProduct{
				FavID:   favID[x],
				ID:      result.ID,
				Name:    result.Name,
				Price:   result.Price,
				Picture: result.Picture,
			})
		}

		responses.Product = prod

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get all data",
			"data":    responses,
		})

	}
}

func (uc *UserController) DelFavorite() echo.HandlerFunc {
	return func(c echo.Context) error {
		favoriteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "ID favorite tidak valid",
			})
		}

		fmt.Println(favoriteID)
		err = uc.srv.DelFavorite(c.Get("user").(*gojwt.Token), uint(favoriteID))
		if err != nil {
			c.Logger().Error("ERROR Delete User, explain:", err.Error())
			var statusCode = http.StatusInternalServerError
			var message = "terjadi permasalahan ketika menghapus favorite"

			if strings.Contains(err.Error(), "tidak ditemukan") {
				statusCode = http.StatusNotFound
				message = "favorite tidak ditemukan"
			} else if strings.Contains(err.Error(), "tidak memiliki izin") {
				statusCode = http.StatusForbidden
				message = "Anda tidak memiliki izin untuk menghapus favorite ini"
			}

			return c.JSON(statusCode, map[string]interface{}{
				"message": message,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success delete favorite",
		})
	}
}
