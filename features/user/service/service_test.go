package service_test

import (
	"errors"
	"testing"

	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/features/user/service"

	"BE-hi-SPEC/features/user/mocks"
	eMock "BE-hi-SPEC/helper/enkrip/mocks"
	"BE-hi-SPEC/helper/jwt"
	jMock "BE-hi-SPEC/helper/jwt/mocks"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := eMock.NewHashInterface(t)
	m := service.New(repo, enkrip, nil)

	var newUser = user.User{
		Email:       "john.doe@example.com",
		Name:        "John Doe",
		Address:     "123 Main St",
		PhoneNumber: "123456789",
		Password:    "password123",
	}

	t.Run("Success Case", func(t *testing.T) {
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		repo.On("InsertUser", newUser).Return(user.User{ID: 1, Email: newUser.Email}, nil).Once()

		result, err := m.Register(newUser)

		assert.Nil(t, err)
		assert.Equal(t, user.User{ID: 1, Email: newUser.Email}, result)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Duplicate Data", func(t *testing.T) {
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		repo.On("InsertUser", newUser).Return(user.User{}, errors.New("duplicate key")).Once()

		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "data telah terdaftar pada sistem", err.Error())

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - General Error", func(t *testing.T) {
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		repo.On("InsertUser", newUser).Return(user.User{}, errors.New("general error")).Once()

		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
	t.Run("Error Case - Empty Email", func(t *testing.T) {
		newUser := user.User{
			Email:       "",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "password123",
		}

		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "email cannot be empty", err.Error())

	})
	t.Run("Error Case - Empty Name", func(t *testing.T) {
		newUser := user.User{
			Email:       "john.doe@example.com",
			Name:        "",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "password123",
		}

		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "name cannot be empty", err.Error())

	})

	t.Run("Error Case - Empty Address", func(t *testing.T) {
		newUser := user.User{
			Email:       "john.doe@example.com",
			Name:        "John Doe",
			Address:     "",
			PhoneNumber: "123456789",
			Password:    "password123",
		}

		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "address cannot be empty", err.Error())

	})

	t.Run("Error Case - Empty Phone Number", func(t *testing.T) {
		newUser := user.User{
			Email:       "john.doe@example.com",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "",
			Password:    "password123",
		}

		_, err := m.Register(newUser)

		// Periksa hasil
		assert.Error(t, err)
		assert.Equal(t, "Phone Number cannot be empty", err.Error())

	})

	t.Run("Error Case - Empty password", func(t *testing.T) {
		newUser := user.User{
			Email:       "john.doe@example.com",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "",
		}

		_, err := m.Register(newUser)

		assert.Error(t, err)
		assert.Equal(t, "password cannot be empty", err.Error())

	})
	t.Run("Error Case - Encryption Failure", func(t *testing.T) {
		enkrip.On("HashPassword", newUser.Password).Return("", errors.New("terdapat masalah saat memproses enkripsi password")).Once()
		_, err := m.Register(newUser)
		assert.Error(t, err)
		assert.Equal(t, "terdapat masalah saat memproses enkripsi password", err.Error())
		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	hashMock := eMock.NewHashInterface(t)
	userService := service.New(repo, hashMock, nil)

	email := "test@example.com"
	password := "test_password"
	hashedPassword := "hashed_test_password"

	t.Run("Success Case", func(t *testing.T) {
		repo.On("Login", email).Return(user.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
		}, nil).Once()

		hashMock.On("Compare", hashedPassword, password).Return(nil).Once()

		result, err := userService.Login(email, password)

		assert.Nil(t, err)
		assert.Equal(t, user.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
		}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - Empty Email and Password", func(t *testing.T) {
		result, err := userService.Login("", "")

		assert.Error(t, err)
		assert.Equal(t, "email and password are required", err.Error())
		assert.Equal(t, user.User{}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - User Not Found", func(t *testing.T) {
		repo.On("Login", email).Return(user.User{}, errors.New("not found")).Once()

		result, err := userService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "data tidak ditemukan", err.Error())
		assert.Equal(t, user.User{}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - Incorrect Password", func(t *testing.T) {
		repo.On("Login", email).Return(user.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
		}, nil).Once()

		hashMock.On("Compare", hashedPassword, password).Return(errors.New("password salah")).Once()

		result, err := userService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "password salah", err.Error())
		assert.Equal(t, user.User{}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - System Error", func(t *testing.T) {
		repo.On("Login", email).Return(user.User{}, errors.New("data tidak ditemukan")).Once()

		result, err := userService.Login(email, password)

		assert.Error(t, err)
		assert.Equal(t, "data tidak ditemukan", err.Error())
		assert.Equal(t, user.User{}, result)

		repo.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}

func TestHapusUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	jwtMock := jMock.NewJWTService(t)
	userService := service.New(repoMock, nil, jwtMock)

	var userID = uint(1)
	var rolesUser = "admin"
	var str, _ = jwt.GenerateJWT(userID, rolesUser)
	var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
		return []byte("$!1gnK3yyy!!!"), nil
	})

	t.Run("Success Case - Admin Deletes User", func(t *testing.T) {
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: userID}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(nil).Once()

		err := userService.HapusUser(token, userID)

		assert.Nil(t, err)
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Success Case - User Deletes Own Account", func(t *testing.T) {
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: userID}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(nil).Once()

		err := userService.HapusUser(token, userID)

		assert.Nil(t, err)
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - User Not Found", func(t *testing.T) {
		repoMock.On("GetUserByID", userID).Return(nil, errors.New("user not found")).Once()

		err := userService.HapusUser(token, userID)

		assert.Error(t, err)
		assert.Equal(t, "failed to retrieve the user for deletion", err.Error())

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Delete User", func(t *testing.T) {
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: userID}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(errors.New("delete user error")).Once()

		err := userService.HapusUser(token, userID)

		assert.Error(t, err)
		assert.Equal(t, "failed to delete the user", err.Error())

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})
}
