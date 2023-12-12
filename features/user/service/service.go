package service

import (
	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/helper/enkrip"
	"errors"
	"strings"
)

type UserService struct {
	repo user.Repository
	hash enkrip.HashInterface
}

func New(r user.Repository, h enkrip.HashInterface) user.Service {
	return &UserService{
		repo: r,
		hash: h,
	}
}

// Login implements user.Service.
func (us *UserService) Login(email string, password string) (user.User, error) {
	if email == "" || password == "" {
		return user.User{}, errors.New("email and password are required")
	}
	result, err := us.repo.Login(email)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return user.User{}, errors.New("data tidak ditemukan")
		}
		return user.User{}, errors.New("terjadi kesalahan pada sistem")
	}

	err = us.hash.Compare(result.Password, password)

	if err != nil {
		return user.User{}, errors.New("password salah")
	}

	return result, nil
}
