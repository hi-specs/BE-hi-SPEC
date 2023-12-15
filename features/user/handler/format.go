package handler

import (
	"mime/multipart"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type RegisterRequest struct {
	Name        string `json:"name" form:"name"`
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

type RegisterResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email" form:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type PutRequest struct {
	ID          uint           `json:"id" form:"id"`
	Name        string         `json:"name" form:"name"`
	Email       string         `json:"email" form:"email"`
	Avatar      multipart.File `json:"avatar" form:"avatar"`
	Address     string         `json:"address" form:"address"`
	PhoneNumber string         `json:"phone_number" form:"phone_number"`
	Password    string         `json:"password" form:"password"`
	NewPassword string         `json:"newpassword" form:"newpassword"`
}

type PutResponse struct {
	ID          uint   `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Email       string `json:"email" form:"email"`
	Avatar      string `json:"avatar" form:"avatar"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

type GetUserResponse struct {
	ID          uint   `json:"id" form:"id"`
	Email       string `json:"email" form:"email"`
	Name        string `json:"name" form:"name"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Avatar      string `json:"avatar" form:"avatar"`
}

type GetAllUserResponse struct {
	Users []GetUserResponse
}

type FavoriteRequest struct {
	ProductID uint `json:"product_id" form:"product_id"`
}

type GetAllFavoriteUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type GetAllFavoriteProduct struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
type GetAllFavoriteResponse struct {
	User    GetAllFavoriteUser
	Product []GetAllFavoriteProduct
}
