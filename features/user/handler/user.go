package handler

import (
	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/helper/jwt"
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
)

type UserController struct {
	srv user.Service
	// cl     *cloudinary.Cloudinary
	// ct     context.Context
	// folder string
}

func New(s user.Service) user.Handler {
	return &UserController{
		srv: s,
		// cl:     cld,
		// ct:     ctx,
		// folder: uploadparam,
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
		response.Email = result.Email
		response.ID = result.ID
		response.Password = result.Password
		response.Token = strToken

		return c.JSON(http.StatusOK, map[string]any{
			"message": "success login data",
			"data":    response,
		})
	}
}
