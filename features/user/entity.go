package user

import "github.com/labstack/echo/v4"

type User struct {
	ID          uint
	Email       string
	Name        string
	Address     string
	PhoneNumber string
	Password    string
}

type Handler interface {
	Login() echo.HandlerFunc
}

type Service interface {
	Login(email string, password string) (User, error)
}

type Repository interface {
	Login(email string) (User, error)
}
