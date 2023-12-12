package main

import (
	"BE-hi-SPEC/config"
	"BE-hi-SPEC/routes"
	"BE-hi-SPEC/utils/database"

	uh "BE-hi-SPEC/features/user/handler"
	ur "BE-hi-SPEC/features/user/repository"
	us "BE-hi-SPEC/features/user/service"
	ek "BE-hi-SPEC/helper/enkrip"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	if cfg == nil {
		e.Logger.Fatal("tidak bisa start server kesalahan database")
	}
	// cld, ctx, param := cld.InitCloudnr(*cfg)

	db, err := database.InitMySql(*cfg)
	if err != nil {
		e.Logger.Fatal("tidak bisa start bro", err.Error())
	}

	db.AutoMigrate(&ur.UserModel{})

	ekrip := ek.New()
	userRepo := ur.New(db)
	userService := us.New(userRepo, ekrip)
	userHandler := uh.New(userService)
	// userHandler := uh.New(userService, cld, ctx, param)

	routes.InitRoute(e, userHandler)
	e.Logger.Fatal(e.Start(":8000"))
}
