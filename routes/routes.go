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
	e.GET("/users", uc.All())
	e.POST("/login", uc.Login())
	e.POST("/register", uc.Register())
	e.PATCH("/user/:id", uc.Update(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/user/:id", uc.Delete(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.POST("/user/fav/add/:id", uc.AddFavorite(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/user/fav/list/:id", uc.GetAllFavorite())
	e.DELETE("/user/fav/del/:id", uc.DelFavorite(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func RouteProduct(e *echo.Echo, ph product.Handler) {
	e.POST("/product", ph.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/product", ph.GetAll())
	e.GET("/product/:id", ph.GetProductDetail())
	e.GET("/product/search", ph.SearchProductByName())
	e.GET("/product/search/category", ph.SearchProductByCategory())
}
