package main

import (
	"BE-hi-SPEC/config"
	"BE-hi-SPEC/routes"
	"BE-hi-SPEC/utils/cld"
	"BE-hi-SPEC/utils/database"

	// "BE-hi-SPEC/features/product"
	uh "BE-hi-SPEC/features/user/handler"
	ur "BE-hi-SPEC/features/user/repository"
	us "BE-hi-SPEC/features/user/service"

	ph "BE-hi-SPEC/features/product/handler"
	pr "BE-hi-SPEC/features/product/repository"
	ps "BE-hi-SPEC/features/product/service"

	ek "BE-hi-SPEC/helper/enkrip"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	if cfg == nil {
		e.Logger.Fatal("tidak bisa start server kesalahan database")
	}
	cld, ctx, param := cld.InitCloudnr(*cfg)

	db, err := database.InitMySql(*cfg)
	if err != nil {
		e.Logger.Fatal("tidak bisa start bro", err.Error())
	}

	db.AutoMigrate(&ur.UserModel{}, &pr.ProductModel{}, &ur.FavoriteModel{})

	ekrip := ek.New()
	userRepo := ur.New(db)
	userService := us.New(userRepo, ekrip)
	userHandler := uh.New(userService, cld, ctx, param)

	productRepo := pr.New(db)
	productService := ps.New(productRepo)
	productHandler := ph.New(productService, cld, ctx, param)

	routes.InitRoute(e, userHandler, productHandler)
	e.Logger.Fatal(e.Start(":8000"))
}
