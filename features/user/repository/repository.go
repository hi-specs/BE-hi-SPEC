package repository

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/user"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Name        string `json:"name" form:"name"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Password    string `json:"password" form:"password"`
	Avatar      string `json:"avatar" form:"avatar"`
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

func (uq *UserQuery) GetAllUser(page int, limit int) ([]user.User, int, error) {
	var Users []UserModel
	offset := (page - 1) * limit
	if err := uq.db.Offset(offset).Limit(limit).Find(&Users).Error; err != nil {
		return nil, 0, err
	}
	var result []user.User
	for _, resp := range Users {
		results := user.User{
			ID:          resp.ID,
			Email:       resp.Email,
			Name:        resp.Name,
			Address:     resp.Address,
			PhoneNumber: resp.PhoneNumber,
			Avatar:      resp.Avatar,
			Model:       resp.Model,
		}
		result = append(result, results)
	}

	var totalPage int
	tableNameUser := "user_models"
	columnNameUser := "deleted_at"
	queryuser := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NULL", tableNameUser, columnNameUser)
	err := uq.db.Raw(queryuser).Scan(&totalPage).Error
	if err != nil {
		log.Fatal(err)
	}

	if totalPage/limit == 0 {
		totalPage = totalPage / limit
	} else {
		totalPage = totalPage / limit
		totalPage++
	}

	if totalPage == 0 {
		totalPage++
	}

	return result, totalPage, err
}

func (uq *UserQuery) AddFavorite(userID, productID uint) (user.Favorite, error) {
	var Fav FavoriteModel
	Fav.UserID = userID
	Fav.ProductID = productID
	if err := uq.db.Create(&Fav).Error; err != nil {
		return user.Favorite{}, err
	}

	var User user.User
	uq.db.Table("user_models").Where("id = ?", userID).Find(&User)

	var FavList []uint
	uq.db.Table("favorite_models").Where("user_id = ?", userID).Select("product_id").Find(&FavList)

	fmt.Println(FavList)
	var Favorite []product.Product
	uq.db.Table("product_models").Where("id = ?", productID).Find(&Favorite)

	var result user.Favorite

	result.FavID = FavList
	result.Favorite = Favorite
	result.User = User

	return result, nil
}

func (uq *UserQuery) GetUser(userID uint) (user.Favorite, error) {
	var User user.User
	uq.db.Table("user_models").Where("id = ?", userID).Find(&User)

	var FavList []uint
	uq.db.Table("favorite_models").Where("user_id = ?", userID).Select("product_id").Find(&FavList)
	var Favorite []product.Product

	for _, fav := range FavList {
		tmp := new(product.Product)
		uq.db.Table("product_models").Where("id = ?", fav).Find(&tmp)
		Favorite = append(Favorite, *tmp)
	}

	var FavID []uint
	uq.db.Table("favorite_models").Where("user_id = ?", userID).Select("id").Find(&FavID)

	var result user.Favorite
	result.FavID = FavID
	result.Favorite = Favorite
	result.User = User

	return result, nil
}

func (uq *UserQuery) DelFavorite(favoriteID uint) error {
	var Fav FavoriteModel

	if err := uq.db.First(&Fav, favoriteID).Error; err != nil {
		return err
	}
	if err := uq.db.Delete(&Fav).Error; err != nil {
		return err
	}

	return nil
}

func (uq *UserQuery) SearchUser(name string, page int, limit int) ([]user.User, int, error) {
	var users []user.User
	offset := (page - 1) * limit
	qry := uq.db.Table("user_models").Offset(offset).Limit(limit)

	if name != "" {
		qry = qry.Where("name like ?", "%"+name+"%")
		qry = qry.Where("deleted_at IS NULL")
	}

	var totalPage int
	tableNameUser := "user_models"
	columnNameUser := "deleted_at"
	queryuser := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NULL", tableNameUser, columnNameUser)
	err := uq.db.Raw(queryuser).Scan(&totalPage).Error
	if err != nil {
		log.Fatal(err)
	}

	if totalPage/limit == 0 {
		totalPage = totalPage / limit
	} else {
		totalPage = totalPage / limit
		totalPage++
	}

	if totalPage == 0 {
		totalPage++
	}

	if err := qry.Find(&users).Error; err != nil {
		return nil, totalPage, err
	}

	return users, totalPage, nil
}
