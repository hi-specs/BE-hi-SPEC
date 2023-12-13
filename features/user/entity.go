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
	Register() echo.HandlerFunc
}

type Service interface {
	Login(email string, password string) (User, error)
	Register(newUser User) (User, error)
}

type Repository interface {
	Login(email string) (User, error)
	InsertUser(newUser User) (User, error)
}
