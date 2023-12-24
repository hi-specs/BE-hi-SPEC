package service_test

import (
	"errors"
	"testing"

	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/features/user/service"

	"BE-hi-SPEC/features/user/mocks"
	eMock "BE-hi-SPEC/helper/enkrip/mocks"
	jMock "BE-hi-SPEC/helper/jwt/mocks"

	gojwt "github.com/golang-jwt/jwt/v5"
	golangjwt "github.com/golang-jwt/jwt/v5"

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

	t.Run("empty role", func(t *testing.T) {
		token := &gojwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(1), "", nil).Once()

		err := userService.HapusUser(token, 1)

		assert.Nil(t, err)
		jwtMock.AssertExpectations(t)
	})

	t.Run("Success Case - Admin Deletes User", func(t *testing.T) {
		token := &gojwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(1), "admin", nil).Once()
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: userID}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(nil).Once()

		err := userService.HapusUser(token, userID)

		assert.Nil(t, err)
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Success Case - User Deletes Own Account", func(t *testing.T) {
		token := &gojwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(1), "user", nil).Once()
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: userID}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(nil).Once()

		err := userService.HapusUser(token, userID)

		assert.Nil(t, err)
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - User Delete Another User", func(t *testing.T) {
		token := &gojwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "user", nil).Once()
		err := userService.HapusUser(token, userID)
		assert.Error(t, err)
		assert.Equal(t, "you don't have permission to delete this user", err.Error())
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Retrieve User for Deletion", func(t *testing.T) {
		token := &gojwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "admin", nil).Once()
		repoMock.On("GetUserByID", userID).Return(nil, errors.New("failed to retrieve user")).Once()
		err := userService.HapusUser(token, userID)
		assert.Error(t, err)
		assert.Equal(t, "failed to retrieve the user for deletion", err.Error())
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Extract User ID and Roles from Token", func(t *testing.T) {
		token := &gojwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(0), "", errors.New("failed to extract user ID and roles")).Once()
		err := userService.HapusUser(token, userID)
		assert.Error(t, err)
		assert.Equal(t, "failed to extract user ID and roles", err.Error())
		jwtMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Delete User", func(t *testing.T) {
		token := &gojwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "admin", nil).Once()
		repoMock.On("GetUserByID", userID).Return(&user.User{ID: uint(123), Role: "admin"}, nil).Once()
		repoMock.On("DeleteUser", userID).Return(errors.New("failed to delete user")).Once()
		err := userService.HapusUser(token, userID)
		assert.Error(t, err)
		assert.Equal(t, "failed to delete the user", err.Error())
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

}

func TestGetAllUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	jwtMock := jMock.NewJWTService(t)
	userService := service.New(repoMock, nil, jwtMock)

	t.Run("Admin User - Successful Get All Users", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "admin", nil).Once()
		expectedUsers := []user.User{{ID: 1, Name: "User1"}, {ID: 2, Name: "User2"}}
		expectedTotalPage := 3
		repoMock.On("GetAllUser", uint(990), 1, 10).Return(expectedUsers, expectedTotalPage, nil).Once()
		users, totalPage, err := userService.GetAllUser(token, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, expectedTotalPage, totalPage)
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Non-Admin User - Unauthorized Access", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "user", nil).Once()
		users, totalPage, err := userService.GetAllUser(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, "unauthorized access: admin role required", err.Error())
		assert.Nil(t, users)
		assert.Equal(t, 0, totalPage)
		jwtMock.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Error", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "admin", nil).Once()
		repoMock.On("GetAllUser", uint(990), 1, 10).Return(nil, 0, errors.New("system error")).Once()
		users, totalPage, err := userService.GetAllUser(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())
		assert.Nil(t, users)
		assert.Equal(t, 0, totalPage)
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - User Not Found", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "admin", nil).Once()
		repoMock.On("GetAllUser", uint(990), 1, 10).Return(nil, 0, errors.New("user not found")).Once()
		users, totalPage, err := userService.GetAllUser(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, "username tidak ditemukan", err.Error())
		assert.Nil(t, users)
		assert.Equal(t, 0, totalPage)
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Extract Token", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(0), "", errors.New("failed to extract token")).Once()

		users, totalPage, err := userService.GetAllUser(token, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, 0, totalPage)
		assert.Equal(t, "failed to extract token", err.Error())

		jwtMock.AssertExpectations(t)
	})
}

func TestAddFavorite(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	jwtMock := jMock.NewJWTService(t)
	userService := service.New(repoMock, nil, jwtMock)

	t.Run("Successful Add Favorite", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "user", nil).Once()
		expectedFavorite := user.Favorite{User: user.User{ID: uint(123)}, FavID: []uint{456}}
		repoMock.On("AddFavorite", uint(123), uint(456)).Return(expectedFavorite, nil).Once()
		favorite, err := userService.AddFavorite(token, 456)
		assert.NoError(t, err)
		assert.Equal(t, expectedFavorite, favorite)
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Extract Token", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(0), "", errors.New("failed to extract token")).Once()
		favorite, err := userService.AddFavorite(token, 456)
		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorite)
		assert.Equal(t, "failed to extract token", err.Error())
		jwtMock.AssertExpectations(t)
		// Assuming repoMock.AssertExpectations(t) is not relevant for this specific test case
	})

	t.Run("Error Case - Duplicate Favorite", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "user", nil).Once()
		repoMock.On("AddFavorite", uint(123), uint(456)).Return(user.Favorite{}, errors.New("duplicate entry")).Once()
		favorite, err := userService.AddFavorite(token, 456)
		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorite)
		assert.Equal(t, "favorite telah ditambahkan", err.Error())
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Error", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "user", nil).Once()
		repoMock.On("AddFavorite", uint(123), uint(456)).Return(user.Favorite{}, errors.New("system error")).Once()
		favorite, err := userService.AddFavorite(token, 456)
		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorite)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())
		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

}

func TestSearchUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	jwtMock := jMock.NewJWTService(t)
	userService := service.New(repoMock, nil, jwtMock)

	t.Run("Admin User - Successful Search User", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "admin", nil).Once()
		expectedUsers := []user.User{{ID: 1, Name: "User1"}, {ID: 2, Name: "User2"}}
		expectedTotalPage := 3
		repoMock.On("SearchUser", uint(990), "John Doe", 1, 10).Return(expectedUsers, expectedTotalPage, nil).Once()

		users, totalPage, err := userService.SearchUser(token, "John Doe", 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, expectedTotalPage, totalPage)

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Non-Admin User - Unauthorized Access", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "user", nil).Once()

		users, totalPage, err := userService.SearchUser(token, "John Doe", 1, 10)

		assert.Error(t, err)
		assert.Equal(t, "unauthorized access: admin role required", err.Error())
		assert.Nil(t, users)
		assert.Equal(t, 0, totalPage)

		jwtMock.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Error", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(990), "admin", nil).Once()
		repoMock.On("SearchUser", uint(990), "John Doe", 1, 10).Return(nil, 0, errors.New("system error")).Once()

		users, totalPage, err := userService.SearchUser(token, "John Doe", 1, 10)

		assert.Error(t, err)
		assert.Equal(t, "failed get all user", err.Error())
		assert.Nil(t, users)
		assert.Equal(t, 0, totalPage)

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	jwtMock := jMock.NewJWTService(t)
	userService := service.New(repoMock, nil, jwtMock)

	t.Run("Successful Get User", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "user", nil).Once()
		expectedFavorites := user.Favorite{User: user.User{
			ID: uint(123),
		}}
		repoMock.On("GetUser", uint(123)).Return(expectedFavorites, nil).Once()

		favorites, err := userService.GetUser(token)

		assert.NoError(t, err)
		assert.Equal(t, expectedFavorites, favorites)

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Extract Token", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(0), "", errors.New("failed to extract token")).Once()

		favorites, err := userService.GetUser(token)

		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorites)
		assert.Equal(t, "failed to extract token", err.Error())

		jwtMock.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Error", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "user", nil).Once()
		repoMock.On("GetUser", uint(123)).Return(user.Favorite{}, errors.New("system error")).Once()

		favorites, err := userService.GetUser(token)

		assert.Error(t, err)
		assert.Equal(t, user.Favorite{}, favorites)
		assert.Equal(t, "system error", err.Error())

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
	})
}

func TestDelFavorite(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	jwtMock := jMock.NewJWTService(t)
	userService := service.New(repoMock, nil, jwtMock)

	t.Run("Successful Delete Favorite", func(t *testing.T) {
		favoriteID := uint(456)
		repoMock.On("DelFavorite", favoriteID).Return(nil).Once()

		err := userService.DelFavorite(&golangjwt.Token{}, favoriteID)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Delete Favorite", func(t *testing.T) {
		favoriteID := uint(456)
		repoMock.On("DelFavorite", favoriteID).Return(errors.New("failed to delete favorite")).Once()

		err := userService.DelFavorite(&golangjwt.Token{}, favoriteID)

		assert.Error(t, err)
		assert.Equal(t, "failed to delete the favorite", err.Error())
		repoMock.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	jwtMock := jMock.NewJWTService(t)
	hashMock := eMock.NewHashInterface(t)

	userService := service.New(repoMock, hashMock, jwtMock)

	t.Run("Admin User - Successful Update User", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "admin", nil).Once()

		input := user.User{ID: uint(123), Name: "UpdatedUser", NewPassword: "newpass"}
		baseUser := user.User{ID: uint(123), Name: "User123", Password: "oldpass"}
		repoMock.On("GetUserByID", uint(123)).Return(&baseUser, nil).Once()
		hashMock.On("HashPassword", "newpass").Return("hashednewpass", nil).Once()
		input.NewPassword = "hashednewpass"
		repoMock.On("UpdateUser", input).Return(user.User{ID: uint(123), Name: "UpdatedUser", Password: "hashednewpass"}, nil).Once()

		input.NewPassword = "newpass"
		updatedUser, err := userService.UpdateUser(token, input)

		assert.NoError(t, err)
		assert.Equal(t, user.User{ID: uint(123), Name: "UpdatedUser", Password: "hashednewpass"}, updatedUser)

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Admin User - Error Updating User in Database", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "admin", nil).Once()
		input := user.User{ID: uint(123), Name: "UpdatedUser", NewPassword: "newpass"}
		baseUser := user.User{ID: uint(123), Name: "User123", Password: "oldpass"}
		repoMock.On("GetUserByID", uint(123)).Return(&baseUser, nil).Once()
		hashMock.On("HashPassword", "newpass").Return("hashednewpass", nil).Once()
		input.NewPassword = "hashednewpass"
		repoMock.On("UpdateUser", input).Return(user.User{}, errors.New("failed to update user in database")).Once()

		input.NewPassword = "newpass"
		updatedUser, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, updatedUser)
		assert.Equal(t, "kesalahan pada database", err.Error())

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Mismatched User ID", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(456), "user", nil).Once()
		input := user.User{ID: uint(789), Name: "UpdatedUser", NewPassword: "newpass"}

		updatedUser, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, updatedUser)
		assert.Equal(t, "id tidak cocok", err.Error())

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Admin User - User Not Found", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "admin", nil).Once()
		input := user.User{ID: uint(123), Name: "UpdatedUser", NewPassword: "newpass"}
		repoMock.On("GetUserByID", uint(123)).Return(&user.User{}, errors.New("user not found")).Once()

		updatedUser, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, updatedUser)
		assert.Equal(t, "user tidak ditemukan", err.Error())

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	// password := "test_password"
	// hashedPassword := "hashed_test_password"

	// t.Run("Admin User - Successful Update User with New Password", func(t *testing.T) {
	// 	token := &golangjwt.Token{}
	// 	jwtMock.On("ExtractToken", token).Return(uint(123), "admin", nil).Once()
	// 	input := user.User{ID: uint(123), Name: "UpdatedUser", NewPassword: "newpass"}
	// 	baseUser := user.User{ID: uint(123), Name: "User123", Password: "oldpass"}
	// 	repoMock.On("GetUserByID", uint(123)).Return(baseUser, nil).Once()
	// 	hashMock.On("HashPassword", "newpass").Return("hashednewpass", nil).Once()
	// 	repoMock.On("UpdateUser", input).Return(user.User{ID: uint(123), Name: "UpdatedUser", Password: "hashednewpass"}, nil).Once()

	// 	updatedUser, err := userService.UpdateUser(token, input)

	// 	assert.NoError(t, err)
	// 	assert.Equal(t, user.User{ID: uint(123), Name: "UpdatedUser", Password: "hashednewpass"}, updatedUser)

	// 	jwtMock.AssertExpectations(t)
	// 	repoMock.AssertExpectations(t)
	// 	hashMock.AssertExpectations(t)
	// })

	// t.Run("Admin User - Error Hashing New Password", func(t *testing.T) {
	// 	token := &golangjwt.Token{}
	// 	jwtMock.On("ExtractToken", token).Return(uint(123), "admin", nil).Once()
	// 	input := user.User{ID: uint(123), Name: "UpdatedUser", NewPassword: "newpass"}
	// 	baseUser := user.User{ID: uint(123), Name: "User123", Password: "oldpass"}
	// 	repoMock.On("GetUserByID", uint(123)).Return(baseUser, nil).Once()
	// 	hashMock.On("HashPassword", "newpass").Return("", errors.New("failed to hash password")).Once()

	// 	updatedUser, err := userService.UpdateUser(token, input)

	// 	assert.Error(t, err)
	// 	assert.Equal(t, user.User{}, updatedUser)
	// 	assert.Equal(t, "masukkan password baru dengan benar", err.Error())

	// 	jwtMock.AssertExpectations(t)
	// 	repoMock.AssertExpectations(t)
	// 	hashMock.AssertExpectations(t)
	// })

	t.Run("Admin User - User Not Found", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(123), "admin", nil).Once()
		input := user.User{ID: 123, Name: "UpdatedUser", NewPassword: "newpass"}
		repoMock.On("GetUserByID", uint(123)).Return(&user.User{}, errors.New("user not found")).Once()

		updatedUser, err := userService.UpdateUser(token, input)

		assert.Error(t, err)
		assert.Equal(t, user.User{}, updatedUser)
		assert.Equal(t, "user tidak ditemukan", err.Error())

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})

	t.Run("Error Case - Failed to Extract Token", func(t *testing.T) {
		token := &golangjwt.Token{}
		jwtMock.On("ExtractToken", token).Return(uint(0), "", errors.New("failed to extract token")).Once()

		updatedUser, err := userService.UpdateUser(token, user.User{})

		assert.Error(t, err)
		assert.Equal(t, user.User{}, updatedUser)
		assert.Equal(t, "harap login", err.Error())

		jwtMock.AssertExpectations(t)
		repoMock.AssertExpectations(t)
		hashMock.AssertExpectations(t)
	})
}
