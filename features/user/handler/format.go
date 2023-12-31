package handler

import (
	"time"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	ID       uint   `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Role     string `json:"role"`
}

type RegisterRequest struct {
	Name        string `json:"name" form:"name"`
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

type RegisterResponse struct {
	ID          uint   `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email" form:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type PutRequest struct {
	ID          uint   `json:"user_id" form:"user_id"`
	Name        string `json:"name" form:"name"`
	Email       string `json:"email" form:"email"`
	Avatar      string `json:"avatar" form:"avatar"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Password    string `json:"password" form:"password"`
	NewPassword string `json:"newpassword" form:"newpassword"`
}

type PutResponse struct {
	ID          uint   `json:"user_id" form:"user_id"`
	Name        string `json:"name" form:"name"`
	Email       string `json:"email" form:"email"`
	Avatar      string `json:"avatar" form:"avatar"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

type GetUserResponse struct {
	ID          uint      `json:"user_id" form:"user_id"`
	Email       string    `json:"email" form:"email"`
	Name        string    `json:"name" form:"name"`
	Address     string    `json:"address" form:"address"`
	PhoneNumber string    `json:"phone_number" form:"phone_number"`
	Avatar      string    `json:"avatar" form:"avatar"`
	Time        time.Time `json:"time" form:"time"`
	Role        string    `json:"role"`
}

type GetAllUserResponse struct {
	Users []GetUserResponse
}

type FavoriteRequest struct {
	ProductID uint `json:"product_id" form:"product_id"`
}

type GetUser struct {
	ID          uint   `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Role        string `json:"role"`
}
type GetAllFavoriteProduct struct {
	FavID   uint   `json:"favorite_id"`
	ID      uint   `json:"product_id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	Picture string `json:"picture"`
}

type GetAllFavoriteResponse struct {
	User        GetUser                 `json:"user"`
	Product     []GetAllFavoriteProduct `json:"my_favorite"`
	Transaction []GetAllTransResponse   `json:"transaction"`
}

type GetAllTransResponse struct {
	ProductPicture string    `json:"product_picture" form:"product_picture"`
	ProductName    string    `json:"product_name" form:"product_name"`
	Nota           string    `json:"nota" form:"nota"`
	TotalPrice     uint      `json:"total_price" form:"total_price"`
	Timestamp      time.Time `json:"timestamp" form:"timestamp"`
	Status         string    `json:"status" form:"status"`
	TransactionID  uint      `json:"transaction_id" form:"transaction_id"`
	ProductID      uint      `json:"product_id" form:"product_id"`
	Token          string    `json:"token" form:"token"`
	Url            string    `json:"url" form:"url"`
}

type SearchUserResponse struct {
	ID          uint      `json:"user_id" form:"user_id"`
	Name        string    `json:"name" form:"name"`
	Email       string    `json:"email" form:"email"`
	Avatar      string    `json:"avatar" form:"avatar"`
	Address     string    `json:"address" form:"address"`
	Time        time.Time `json:"time" form:"time"`
	PhoneNumber string    `json:"phone_number" form:"phone_number"`
	Role        string    `json:"role" form:"role"`
}
