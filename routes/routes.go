package routes

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo, uc user.Handler, ph product.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	RouteUser(e, uc)
	RouteProduct(e, ph)
}

func RouteUser(e *echo.Echo, uc user.Handler) {
	e.POST("/login", uc.Login())
	e.POST("/register", uc.Register())
	e.PATCH("/user/:id", uc.Update(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func RouteProduct(e *echo.Echo, ph product.Handler) {
	e.POST("/product", ph.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}
