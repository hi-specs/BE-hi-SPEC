package repository

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/user"
	"errors"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Name        string
	Address     string
	PhoneNumber string
	Password    string
	Avatar      string
}

type FavoriteModel struct {
	gorm.Model
	UserID    uint
	ProductID uint
}

type UserQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.Repository {
	return &UserQuery{
		db: db,
	}
}

// Login implements user.Repository.
func (uq *UserQuery) Login(email string) (user.User, error) {
	var userData = new(UserModel)

	if err := uq.db.Where("email = ?", email).First(userData).Error; err != nil {
		return user.User{}, err
	}

	var result = new(user.User)
	result.ID = userData.ID
	result.Email = userData.Email
	result.Password = userData.Password

	return *result, nil
}

// InsertUser implements user.Repository.
func (uq *UserQuery) InsertUser(newUser user.User) (user.User, error) {
	var inputDB = new(UserModel)
	inputDB.Email = newUser.Email
	inputDB.Name = newUser.Name
	inputDB.Address = newUser.Address
	inputDB.PhoneNumber = newUser.PhoneNumber
	inputDB.Password = newUser.Password

	if err := uq.db.Create(&inputDB).Error; err != nil {
		return user.User{}, err
	}

	newUser.ID = inputDB.ID

	return newUser, nil
}

// UpdateUser implements user.Repository.
func (uq *UserQuery) UpdateUser(input user.User) (user.User, error) {
	var proses UserModel
	if err := uq.db.First(&proses, input.ID).Error; err != nil {
		return user.User{}, err
	}

	// Jika tidak ada buku ditemukan
	if proses.ID == 0 {
		err := errors.New("user tidak ditemukan")
		return user.User{}, err
	}

	if input.Name != "" {
		proses.Name = input.Name
	}
	if input.Email != "" {
		proses.Email = input.Email
	}
	if input.PhoneNumber != "" {
		proses.PhoneNumber = input.PhoneNumber
	}
	if input.Avatar != "" {
		proses.Avatar = input.Avatar
	}

	if input.Address != "" {
		proses.Address = input.Address
	}

	if input.NewPassword != "" {
		proses.Password = input.NewPassword
	}
	if err := uq.db.Save(&proses).Error; err != nil {

		return user.User{}, err
	}
	result := user.User{
		ID:          proses.ID,
		Name:        proses.Name,
		Email:       proses.Email,
		Address:     proses.Address,
		Avatar:      proses.Avatar,
		Password:    proses.Password,
		PhoneNumber: proses.PhoneNumber,
	}

	return result, nil
}

// GetUserByID implements user.Repository.
func (uq *UserQuery) GetUserByID(userID uint) (*user.User, error) {
	var userModel UserModel
	if err := uq.db.First(&userModel, userID).Error; err != nil {
		return nil, err
	}

	// Jika tidak ada buku ditemukan
	if userModel.ID == 0 {
		err := errors.New("user tidak ditemukan")
		return nil, err
	}

	result := &user.User{
		ID:          userModel.ID,
		Name:        userModel.Name,
		Email:       userModel.Email,
		Address:     userModel.Address,
		Avatar:      userModel.Avatar,
		PhoneNumber: userModel.PhoneNumber,
		Password:    userModel.Password,
	}

	return result, nil
}

// DeleteUser implements user.Repository.
func (uq *UserQuery) DeleteUser(userID uint) error {
	var exitingUser UserModel

	if err := uq.db.First(&exitingUser, userID).Error; err != nil {
		return err
	}

	if err := uq.db.Delete(&exitingUser).Error; err != nil {
		return err
	}

	return nil
}

func (uq *UserQuery) GetAllUser() ([]user.User, error) {
	var Users []UserModel

	err := uq.db.Find(&Users).Error

	var result []user.User
	for _, resp := range Users {
		results := user.User{
			ID:          resp.ID,
			Email:       resp.Email,
			Name:        resp.Name,
			Address:     resp.Address,
			PhoneNumber: resp.PhoneNumber,
			Avatar:      resp.Avatar,
		}
		result = append(result, results)
	}
	return result, err
}

func (uq *UserQuery) AddFavorite(userID, productID uint) (user.Favorite, error) {
	var Fav FavoriteModel
	Fav.UserID = userID
	Fav.ProductID = productID
	uq.db.Create(&Fav)

	var User user.User
	uq.db.Table("user_models").Where("id = ?", userID).Find(&User)

	var FavList []uint
	uq.db.Table("favorite_models").Where("user_id = ?", userID).Select("product_id").Find(&FavList)

	var Favorite []product.Product
	uq.db.Table("product_models").Where("id = ?", productID).Find(&Favorite)

	var result user.Favorite

	result.Favorite = Favorite
	result.User = User

	return result, nil
}
