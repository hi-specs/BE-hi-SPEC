package service_test

import (
	"errors"
	"testing"

	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/features/user/service"

	"BE-hi-SPEC/features/user/mocks"
	eMock "BE-hi-SPEC/helper/enkrip/mocks"
	"BE-hi-SPEC/helper/jwt"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := eMock.NewHashInterface(t)
	m := service.New(repo, enkrip)

	var newUser = user.User{
		Email:       "john.doe@example.com",
		Name:        "John Doe",
		Address:     "123 Main St",
		PhoneNumber: "123456789",
		Password:    "password123",
	}

	t.Run("Success Case", func(t *testing.T) {
		// Mock hash password
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		// Mock repository untuk mengembalikan hasil tanpa error
		repo.On("InsertUser", newUser).Return(user.User{ID: 1, Email: newUser.Email}, nil).Once()

		// Panggil fungsi Register
		result, err := m.Register(newUser)

		// Periksa hasil
		assert.Nil(t, err)
		assert.Equal(t, user.User{ID: 1, Email: newUser.Email}, result)

		// Periksa ekspektasi panggilan fungsi
		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Duplicate Data", func(t *testing.T) {
		// Mock hash password
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		// Mock repository untuk mengembalikan error duplicate
		repo.On("InsertUser", newUser).Return(user.User{}, errors.New("duplicate key")).Once()

		// Panggil fungsi Register
		_, err := m.Register(newUser)

		// Periksa hasil
		assert.Error(t, err)
		assert.Equal(t, "data telah terdaftar pada sistem", err.Error())

		// Periksa ekspektasi panggilan fungsi
		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - General Error", func(t *testing.T) {
		// Mock hash password
		enkrip.On("HashPassword", newUser.Password).Return("password123", nil).Once()

		// Mock repository untuk mengembalikan error umum
		repo.On("InsertUser", newUser).Return(user.User{}, errors.New("general error")).Once()

		// Panggil fungsi Register
		_, err := m.Register(newUser)

		// Periksa hasil
		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())

		// Periksa ekspektasi panggilan fungsi
		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
	t.Run("Error Case - Empty Email", func(t *testing.T) {
		// Atur data pengguna baru dengan email kosong
		newUser := user.User{
			Email:       "",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "password123",
		}

		// Panggil fungsi Register
		_, err := m.Register(newUser)

		// Periksa hasil
		assert.Error(t, err)
		assert.Equal(t, "email cannot be empty", err.Error())

		// Tidak perlu memanggil ekspektasi panggilan fungsi karena tidak ada panggilan fungsi yang terjadi
	})
	t.Run("Error Case - Empty Name", func(t *testing.T) {
		// Atur data pengguna baru dengan email kosong
		newUser := user.User{
			Email:       "john.doe@example.com",
			Name:        "",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "password123",
		}

		// Panggil fungsi Register
		_, err := m.Register(newUser)

		// Periksa hasil
		assert.Error(t, err)
		assert.Equal(t, "name cannot be empty", err.Error())

		// Tidak perlu memanggil ekspektasi panggilan fungsi karena tidak ada panggilan fungsi yang terjadi
	})

	t.Run("Error Case - Empty Address", func(t *testing.T) {
		// Atur data pengguna baru dengan email kosong
		newUser := user.User{
			Email:       "john.doe@example.com",
			Name:        "John Doe",
			Address:     "",
			PhoneNumber: "123456789",
			Password:    "password123",
		}

		// Panggil fungsi Register
		_, err := m.Register(newUser)

		// Periksa hasil
		assert.Error(t, err)
		assert.Equal(t, "address cannot be empty", err.Error())

		// Tidak perlu memanggil ekspektasi panggilan fungsi karena tidak ada panggilan fungsi yang terjadi
	})

	t.Run("Error Case - Empty Phone Number", func(t *testing.T) {
		// Atur data pengguna baru dengan email kosong
		newUser := user.User{
			Email:       "john.doe@example.com",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "",
			Password:    "password123",
		}

		// Panggil fungsi Register
		_, err := m.Register(newUser)

		// Periksa hasil
		assert.Error(t, err)
		assert.Equal(t, "Phone Number cannot be empty", err.Error())

		// Tidak perlu memanggil ekspektasi panggilan fungsi karena tidak ada panggilan fungsi yang terjadi
	})

	t.Run("Error Case - Empty password", func(t *testing.T) {
		// Atur data pengguna baru dengan email kosong
		newUser := user.User{
			Email:       "john.doe@example.com",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "",
		}

		// Panggil fungsi Register
		_, err := m.Register(newUser)

		// Periksa hasil
		assert.Error(t, err)
		assert.Equal(t, "password cannot be empty", err.Error())

		// Tidak perlu memanggil ekspektasi panggilan fungsi karena tidak ada panggilan fungsi yang terjadi
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	hashMock := eMock.NewHashInterface(t)
	userService := service.New(repo, hashMock)

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
		repo.On("Login", email).Return(user.User{}, errors.New("data tidak ditemukan")).Once()

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
	mockRepo := new(mocks.Repository)
	userService := service.New(mockRepo, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		mockRepo.On("GetUserByID", userID).Return(&user.User{ID: userID}, nil).Once()
		mockRepo.On("DeleteUser", userID).Return(nil).Once()

		err := userService.HapusUser(token, userID)

		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Invalid Token", func(t *testing.T) {
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		token.Valid = false

		err := userService.HapusUser(token, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "token tidak valid", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - User Not Found", func(t *testing.T) {
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		token.Valid = true
		mockRepo.On("GetUserByID", userID).Return(&user.User{}, errors.New("not found")).Once()

		err := userService.HapusUser(token, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "failed to retrieve the user for deletion", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Delete User Failure", func(t *testing.T) {
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		token.Valid = true
		mockRepo.On("GetUserByID", userID).Return(&user.User{ID: userID}, nil).Once()
		mockRepo.On("DeleteUser", userID).Return(errors.New("delete user failed")).Once()

		err := userService.HapusUser(token, userID)

		assert.NotNil(t, err)
		assert.Equal(t, "failed to delete the user", err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestAddFavorite(t *testing.T) {
	mockRepo := new(mocks.Repository)
	userService := service.New(mockRepo, nil)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		mockProductID := uint(1)
		mockFavorites := user.Favorite{
			User:  user.User{ID: userID},
			FavID: []uint{mockProductID},
		}
		mockRepo.On("AddFavorite", userID, mockProductID).Return(mockFavorites, nil).Once()

		// Call the method being tested.
		result, err := userService.AddFavorite(token, mockProductID)

		// Assert that the result and error match the expectations.
		assert.Nil(t, err)
		assert.Equal(t, mockFavorites, result)

		// Assert that the expected repository methods were called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Duplicate", func(t *testing.T) {
		// Mock the expected behavior of the repository for the error case.
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		mockProductID := uint(1)
		mockRepo.On("AddFavorite", userID, mockProductID).Return(user.Favorite{}, errors.New("duplicate")).Once()

		// Call the method being tested.
		_, err := userService.AddFavorite(token, mockProductID)

		// Assert that the error matches the expectations.
		assert.NotNil(t, err)
		assert.Equal(t, "favorite telah ditambahkan", err.Error())

		// Assert that the expected repository methods were called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - System Error", func(t *testing.T) {
		// Mock the expected behavior of the repository for the error case.
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		mockProductID := uint(1)
		mockRepo.On("AddFavorite", userID, mockProductID).Return(user.Favorite{}, errors.New("system error")).Once()

		// Call the method being tested.
		_, err := userService.AddFavorite(token, mockProductID)

		// Assert that the error matches the expectations.
		assert.NotNil(t, err)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())

		// Assert that the expected repository methods were called.
		mockRepo.AssertExpectations(t)
	})
}

func TestDelFavorite(t *testing.T) {
	mockRepo := new(mocks.Repository)
	userService := service.New(mockRepo, nil)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		mockFavoriteID := uint(1)
		mockRepo.On("DelFavorite", mockFavoriteID).Return(nil).Once()

		// Call the method being tested.
		err := userService.DelFavorite(nil, mockFavoriteID)

		// Assert that there is no error.
		assert.Nil(t, err)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the expected behavior of the repository for the error case.
		mockFavoriteID := uint(1)
		mockRepo.On("DelFavorite", mockFavoriteID).Return(errors.New("database error")).Once()

		// Call the method being tested.
		err := userService.DelFavorite(nil, mockFavoriteID)

		// Assert that the error matches the expectations.
		assert.NotNil(t, err)
		assert.Equal(t, "failed to delete the favorite", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	userService := service.New(mockRepo, nil)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		mockUserID := uint(1)
		mockFavorites := user.Favorite{
			User: user.User{
				ID:          mockUserID,
				Email:       "john@example.com",
				Name:        "John Doe",
				Address:     "123 Main St",
				PhoneNumber: "123456789",
				Password:    "hashedpassword",
				NewPassword: "",
				Avatar:      "avatar.jpg",
			},
			FavID:    []uint{1, 2, 3},
			Favorite: []product.Product{{ID: 1, Name: "Product 1"}, {ID: 2, Name: "Product 2"}},
		}
		mockRepo.On("GetUser", mockUserID).Return(mockFavorites, nil).Once()

		// Call the method being tested.
		result, err := userService.GetUser(mockUserID)

		// Assert that there is no error.
		assert.Nil(t, err)

		// Assert that the result matches the expectations.
		assert.Equal(t, mockFavorites, result)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the expected behavior of the repository for the error case.
		mockUserID := uint(1)
		mockRepo.On("GetUser", mockUserID).Return(user.Favorite{}, errors.New("database error")).Once()

		// Call the method being tested.
		result, err := userService.GetUser(mockUserID)

		// Assert that the error matches the expectations.
		assert.NotNil(t, err)
		assert.Equal(t, "database error", err.Error())

		// Assert that the result is empty.
		assert.Equal(t, user.Favorite{}, result)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockHash := new(eMock.HashInterface)
	userService := service.New(mockRepo, mockHash)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		mockUserInput := user.User{
			ID:          userID,
			Email:       "john@example.com",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "oldpassword",
			NewPassword: "newpassword",
			Avatar:      "avatar.jpg",
		}
		mockBaseUser := user.User{
			ID:          userID,
			Email:       "john@example.com",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "oldpassword",
			NewPassword: "newpassword",
			Avatar:      "avatar.jpg",
		}
		mockUpdatedUser := user.User{
			ID:          userID,
			Email:       "john@example.com",
			Name:        "John Doe",
			Address:     "123 Main St",
			PhoneNumber: "123456789",
			Password:    "oldpassword",
			NewPassword: "newpassword",
			Avatar:      "avatar.jpg",
		}

		mockRepo.On("GetUserByID", userID).Return(&mockBaseUser, nil).Once()
		mockHash.On("Compare", mockBaseUser.Password, mockUserInput.Password).Return(nil).Once()
		mockHash.On("HashPassword", mockUserInput.NewPassword).Return("newpassword", nil).Once()
		mockRepo.On("UpdateUser", mockUserInput).Return(mockUpdatedUser, nil).Once()

		// Call the method being tested.
		result, err := userService.UpdateUser(token, mockUserInput)

		// Assert that there is no error.
		assert.Nil(t, err)

		// Assert that the result matches the expectations.
		assert.Equal(t, mockUpdatedUser, result)

		// Assert that the expected repository methods were called.
		mockRepo.AssertExpectations(t)
		mockHash.AssertExpectations(t)
	})

}

func TestSearchUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockHash := new(eMock.HashInterface)
	userService := service.New(mockRepo, mockHash)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		mockUsers := []user.User{
			{ID: 1, Name: "John Doe"},
			{ID: 2, Name: "Jane Doe"},
		}
		mockTotalPage := 2
		mockRepo.On("SearchUser", "John", 1, 10).Return(mockUsers, mockTotalPage, nil).Once()

		// Call the method being tested.
		result, totalPage, err := userService.SearchUser("John", 1, 10)

		// Assert that the result and error match the expectations.
		assert.Nil(t, err)
		assert.Equal(t, mockUsers, result)
		assert.Equal(t, mockTotalPage, totalPage)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the repository to simulate an error.
		mockRepo.On("SearchUser", "John", 1, 10).Return(nil, 0, errors.New("database error")).Once()

		// Call the method being tested.
		_, _, err := userService.SearchUser("John", 1, 10)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed get all user", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	mockHash := new(eMock.HashInterface)
	userService := service.New(mockRepo, mockHash)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		mockUsers := []user.User{
			{ID: 1, Name: "John Doe"},
			{ID: 2, Name: "Jane Doe"},
		}
		mockTotalPage := 2
		mockRepo.On("GetAllUser", 1, 10).Return(mockUsers, mockTotalPage, nil).Once()

		// Call the method being tested.
		result, totalPage, err := userService.GetAllUser(1, 10)

		// Assert that the result and error match the expectations.
		assert.Nil(t, err)
		assert.Equal(t, mockUsers, result)
		assert.Equal(t, mockTotalPage, totalPage)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - User Not Found", func(t *testing.T) {
		// Mock the repository to simulate "not found" error.
		mockRepo.On("GetAllUser", 1, 10).Return(nil, 0, errors.New("not found")).Once()

		// Call the method being tested.
		_, _, err := userService.GetAllUser(1, 10)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "username tidak ditemukan", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the repository to simulate a general error.
		mockRepo.On("GetAllUser", 1, 10).Return(nil, 0, errors.New("database error")).Once()

		// Call the method being tested.
		_, _, err := userService.GetAllUser(1, 10)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}
