package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID          uint
	Email       string
	Name        string
	Address     string
	PhoneNumber string
	Password    string
	NewPassword string
	Avatar      string
}

type Handler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type Service interface {
	Login(email string, password string) (User, error)
	Register(newUser User) (User, error)
	UpdateUser(token *jwt.Token, input User) (User, error)
	HapusUser(token *jwt.Token, userID uint) error
}

type Repository interface {
	Login(email string) (User, error)
	InsertUser(newUser User) (User, error)
	UpdateUser(input User) (User, error)
	GetUserByID(userID uint) (*User, error)
	DeleteUser(userID uint) error
}
