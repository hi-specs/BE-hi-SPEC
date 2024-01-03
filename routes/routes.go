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
	// e.GET("/users", uc.All(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/admin/users", uc.SearchUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))

	e.POST("/login", uc.Login())
	e.POST("/register", uc.Register())
	e.GET("/users", uc.GetUser(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.PATCH("/users/:id", uc.Update(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/users/:id", uc.Delete(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))

	e.POST("/fav/:id", uc.AddFavorite(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/fav/:id", uc.DelFavorite(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))

}

func RouteProduct(e *echo.Echo, ph product.Handler) {
	e.POST("/products", ph.Add(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/products", ph.SearchAll())
	e.GET("/products/:id", ph.GetProductDetail())
	e.PATCH("/products/:id", ph.UpdateProduct(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.DELETE("/products/:id", ph.DelProduct(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	// e.GET("/products", ph.GetAll())
}

func RouteTransaction(e *echo.Echo, th transaction.Handler) {
	e.GET("/admin/dashboard", th.AdminDashboard(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/admin/transactions", th.TransactionList(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))

	e.POST("/transactions", th.Checkout(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/transactions/:id", th.GetTransaction(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.GET("/transactions/:id/download", th.DownloadTransaction(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))

	// e.GET("/users/:id/transactions", th.UserTransaction(), echojwt.JWT([]byte("$!1gnK3yyy!!!")))
	e.POST("/midtrans/callback", th.MidtransCallback())
}
