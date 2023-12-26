package service_test

import (
	"errors"
	"testing"

	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/features/user/service"

	"BE-hi-SPEC/features/user/mocks"
	eMock "BE-hi-SPEC/helper/enkrip/mocks"
	golangjwt "BE-hi-SPEC/helper/jwt"

	"github.com/golang-jwt/jwt/v5"
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
	userService := service.New(repoMock, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		repoMock.On("GetUserByID", userID).Return(&user.User{ID: uint(1)}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(nil).Once()
		err := userService.HapusUser(token, userID)

		assert.Nil(t, err)
		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case Empty Role", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = ""
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		err := userService.HapusUser(token, userID)

		assert.Error(t, err)
	})

	t.Run("Failed Case No Permission", func(t *testing.T) {
		var userID = uint(0)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		err := userService.HapusUser(token, uint(1))

		assert.Error(t, err)
	})

	t.Run("Failed Case Empty User ID", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		repoMock.On("GetUserByID", userID).Return(nil, errors.New("failed to retrieve the user for deletion")).Once()

		err := userService.HapusUser(token, userID)

		assert.Error(t, err)
	})

	t.Run("Failed Case UserID not equal ID user", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: uint(2), Role: "user"}, nil).Once()

		err := userService.HapusUser(token, userID)

		assert.Error(t, err)
	})

	t.Run("Failed Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: uint(2), Role: "user"}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(errors.New("failed to delete the user")).Once()

		err := userService.HapusUser(token, userID)

		assert.Error(t, err)
	})

}

func TestGetAllUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	userService := service.New(repoMock, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		expectedUsers := []user.User{{ID: 1, Name: "User1"}, {ID: 2, Name: "User2"}}
		expectedTotalPage := 3

		repoMock.On("GetAllUser", uint(1), 1, 10).Return(expectedUsers, expectedTotalPage, nil).Once()
		users, totalPage, err := userService.GetAllUser(token, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, expectedTotalPage, totalPage)
	})

	t.Run("Failed Case Roles Not Admin", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		users, totalPage, err := userService.GetAllUser(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []user.User([]user.User(nil)), users)
		assert.Equal(t, 0, totalPage)
	})

	t.Run("Failed Case User Not Found", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		repoMock.On("GetAllUser", userID, 1, 10).Return(nil, 0, errors.New("terjadi kesalahan pada sistem")).Once()
		users, totalPage, err := userService.GetAllUser(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []user.User([]user.User(nil)), users)
		assert.Equal(t, 0, totalPage)
		repoMock.AssertExpectations(t)
	})
}

func TestSearchUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	userService := service.New(repoMock, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		expectedUsers := []user.User{{ID: 1, Name: "User1"}, {ID: 2, Name: "User2"}}
		expectedTotalPage := 3
		repoMock.On("SearchUser", uint(1), "John Doe", 1, 10).Return(expectedUsers, expectedTotalPage, nil).Once()

		users, totalPage, err := userService.SearchUser(token, "John Doe", 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, expectedTotalPage, totalPage)

		repoMock.AssertExpectations(t)

	})

	t.Run("Failed Case Roles Not Admin", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		users, totalPage, err := userService.SearchUser(token, "budi", 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []user.User([]user.User(nil)), users)
		assert.Equal(t, 0, totalPage)
	})

	t.Run("Failed Case User Not Found", func(t *testing.T) {
		var userID = uint(1)
		var name = "budi"
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		repoMock.On("SearchUser", userID, name, 1, 10).Return(nil, 0, errors.New("failed get user")).Once()
		users, totalPage, err := userService.SearchUser(token, name, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []user.User([]user.User(nil)), users)
		assert.Equal(t, 0, totalPage)
		repoMock.AssertExpectations(t)
	})
}

func TestDelFavorite(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	userService := service.New(repoMock, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		favoriteID := uint(1)
		repoMock.On("DelFavorite", favoriteID).Return(nil).Once()

		err := userService.DelFavorite(token, favoriteID)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case Favorite nil", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		favoriteID := uint(456)

		repoMock.On("DelFavorite", favoriteID).Return(errors.New("failed to delete favorite")).Once()

		err := userService.DelFavorite(token, favoriteID)

		assert.Error(t, err)
		assert.Equal(t, "failed to delete the favorite", err.Error())
		repoMock.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	userService := service.New(repoMock, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		expectedFavorites := user.Favorite{User: user.User{
			ID: uint(1),
		}}
		repoMock.On("GetUser", uint(1)).Return(expectedFavorites, nil).Once()

		favorites, err := userService.GetUser(token)

		assert.NoError(t, err)
		assert.Equal(t, expectedFavorites, favorites)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case Roles Nil", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = ""
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		favorites, err := userService.GetUser(token)

		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorites)

		repoMock.AssertExpectations(t)
	})
}

func TestAddFavorite(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	userService := service.New(repoMock, nil)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		expectedFavorite := user.Favorite{User: user.User{ID: uint(1)}, FavID: []uint{456}}
		repoMock.On("AddFavorite", uint(1), uint(456)).Return(expectedFavorite, nil).Once()
		favorite, err := userService.AddFavorite(token, 456)
		assert.NoError(t, err)
		assert.Equal(t, expectedFavorite, favorite)
		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case Roles Nil", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = ""
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		favorite, err := userService.AddFavorite(token, 456)

		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorite)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case Duplicate Favorite", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		repoMock.On("AddFavorite", uint(1), uint(456)).Return(user.Favorite{}, errors.New("duplicate entry")).Once()
		favorite, err := userService.AddFavorite(token, 456)
		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorite)
		assert.Equal(t, "favorite telah ditambahkan", err.Error())
		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case Favorite nil", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		repoMock.On("AddFavorite", uint(1), uint(456)).Return(user.Favorite{}, errors.New("system error")).Once()
		favorite, err := userService.AddFavorite(token, 456)
		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorite)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())
		repoMock.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	enkrip := eMock.NewHashInterface(t)

	userService := service.New(repoMock, enkrip)
	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		input := user.User{ID: uint(1), Name: "UpdatedUser", NewPassword: "newpass"}
		baseUser := user.User{ID: uint(1), Name: "User123", Password: "oldpass"}
		repoMock.On("GetUserByID", uint(1)).Return(&baseUser, nil).Once()
		enkrip.On("HashPassword", "newpass").Return("hashednewpass", nil).Once()
		input.NewPassword = "hashednewpass"
		repoMock.On("UpdateUser", input).Return(user.User{ID: uint(1), Name: "UpdatedUser", Password: "hashednewpass"}, nil).Once()

		input.NewPassword = "newpass"
		updatedUser, err := userService.UpdateUser(token, input)

		assert.NoError(t, err)
		assert.Equal(t, user.User{ID: uint(1), Name: "UpdatedUser", Password: "hashednewpass"}, updatedUser)

		repoMock.AssertExpectations(t)
		enkrip.AssertExpectations(t)
	})

	t.Run("Failed Case Roles Nil", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = ""
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		input := user.User{ID: uint(1), Name: "UpdatedUser", NewPassword: "newpass"}

		users, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})

	t.Run("Failed Case User tidak cocok", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		input := user.User{ID: uint(9), Name: "UpdatedUser", NewPassword: "newpass"}

		users, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, users)

		repoMock.AssertExpectations(t)
	})
}
