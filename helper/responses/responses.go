package responses

import "github.com/labstack/echo/v4"

func PrintResponse(c echo.Context, code int, message string, data any) error {
	var result = map[string]any{}
	if data != nil {
		result["data"] = data
	}
	result["message"] = message
	return c.JSON(code, result)
}
