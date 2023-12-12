package repository

import (
	"BE-hi-SPEC/features/user"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Email       string
	Name        string
	Address     string
	PhoneNumber uint
	Password    string
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
