package routes

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/transaction"
	"BE-hi-SPEC/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoute(e *echo.Echo, uc user.Handler, ph product.Handler, th transaction.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	RouteUser(e, uc)
	RouteProduct(e, ph)
	RouteTransaction(e, th)
}

func RouteUser(e *echo.Echo, uc user.Handler) {
	e.GET("/users", uc.All(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.POST("/login", uc.Login())
	e.POST("/register", uc.Register())
	e.GET("/user", uc.GetUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.PATCH("/user/:id", uc.Update(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/user/:id", uc.Delete(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.POST("/user/fav/add/:id", uc.AddFavorite(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/user/fav/:id", uc.DelFavorite(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/user/search", uc.SearchUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))

}

func RouteProduct(e *echo.Echo, ph product.Handler) {
	e.POST("/product", ph.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/products", ph.GetAll())
	e.GET("/product/:id", ph.GetProductDetail())
	e.GET("/product/search", ph.SearchAll())
	e.PATCH("/product/:id", ph.UpdateProduct(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/product/:id", ph.DelProduct(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
}

func RouteTransaction(e *echo.Echo, th transaction.Handler) {
	e.GET("/dashboard", th.AdminDashboard())
	e.POST("/transaction", th.Checkout(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/transactions", th.TransactionList())
	e.GET("/transaction/:id", th.GetTransaction())
	e.GET("/transaction/user/:id", th.UserTransaction())
	// e.POST("/transaction/download/:id", th.DownloadTransaction(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))

	e.POST("/midtrans/callback", th.MidtransCallback())
}
