package user

import (
	"BE-hi-SPEC/features/product"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
	Avatar      string `json:"avatar"`
	Role        string `json:"role"`
}

type Transaction struct {
	gorm.Model
	Nota       string
	ProductID  uint
	UserID     uint
	TotalPrice uint
	Status     string
	Token      string
	Url        string
}

type Favorite struct {
	User          User
	FavID         []uint
	Favorite      []product.Product
	Transaction   []Transaction
	TransProducts []product.Product
}

// Error implements error.
func (Favorite) Error() string {
	panic("unimplemented")
}

type Handler interface {
	Login() echo.HandlerFunc
	Register() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	All() echo.HandlerFunc
	AddFavorite() echo.HandlerFunc
	GetUser() echo.HandlerFunc
	DelFavorite() echo.HandlerFunc
	SearchUser() echo.HandlerFunc
}

type Service interface {
	Login(email string, password string) (User, error)
	Register(newUser User) (User, error)
	UpdateUser(token *jwt.Token, input User) (User, error)
	HapusUser(token *jwt.Token, userID uint) error
	GetAllUser(token *jwt.Token, page int, limit int) ([]User, int, error)
	AddFavorite(token *jwt.Token, productID uint) (Favorite, error)
	GetUser(token *jwt.Token) (Favorite, error)
	DelFavorite(token *jwt.Token, favoriteID uint) error
	SearchUser(token *jwt.Token, name string, page int, limit int) ([]User, int, error)
}

type Repository interface {
	Login(email string) (User, error)
	InsertUser(newUser User) (User, error)
	UpdateUser(input User) (User, error)
	GetUserByID(userID uint) (*User, error)
	DeleteUser(userID uint) error
	GetAllUser(userID uint, page int, limit int) ([]User, int, error)
	AddFavorite(userID, productID uint) (Favorite, error)
	GetUser(userID uint) (Favorite, error)
	DelFavorite(favoriteID uint, userID uint) error
	SearchUser(userID uint, name string, page int, limit int) ([]User, int, error)
}
