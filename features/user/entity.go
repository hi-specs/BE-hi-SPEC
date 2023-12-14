package user

import (
	"BE-hi-SPEC/features/product"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
	Avatar      string `json:"avatar"`
}

type Favorite struct {
	User     User              `json:"user" form:"user"`
	Favorite []product.Product `json:"favorite" form:"favorite"`
}

type Handler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	All() echo.HandlerFunc
	AddFavorite() echo.HandlerFunc
	GetAllFavorite() echo.HandlerFunc
	DelFavorite() echo.HandlerFunc
}

type Service interface {
	Login(email string, password string) (User, error)
	Register(newUser User) (User, error)
	UpdateUser(token *jwt.Token, input User) (User, error)
	HapusUser(token *jwt.Token, userID uint) error
	GetAllUser() ([]User, error)
	AddFavorite(token *jwt.Token, productID uint) (Favorite, error)
	GetAllFavorite(userID uint) (Favorite, error)
	DelFavorite(token *jwt.Token, favoriteID uint) error
}

type Repository interface {
	Login(email string) (User, error)
	InsertUser(newUser User) (User, error)
	UpdateUser(input User) (User, error)
	GetUserByID(userID uint) (*User, error)
	DeleteUser(userID uint) error
	GetAllUser() ([]User, error)
	AddFavorite(userID, productID uint) (Favorite, error)
	GetAllFavorite(userID uint) (Favorite, error)
	DelFavorite(favoriteID uint) error
}
