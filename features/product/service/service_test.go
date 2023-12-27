package service_test

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/product/mocks"
	"BE-hi-SPEC/features/product/service"
	"errors"
	"testing"

	golangjwt "BE-hi-SPEC/helper/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestTalkToGPT(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		input := product.Product{ID: uint(1), Name: "Asus TUF", Category: "Gaming"}

		repo.On("InsertProduct", uint(1), input).Return(input, nil).Once()

		products, err := productService.TalkToGpt(token, input)

		assert.NoError(t, err, products)
		assert.Equal(t, product.Product{ID: uint(1), Name: "Asus TUF", Category: "Gaming"}, input)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Token Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		input := product.Product{ID: uint(1), Name: "Asus TUF", Category: "Gaming"}

		products, err := productService.TalkToGpt(token, input)

		assert.Error(t, err)
		assert.Equal(t, product.Product{}, products)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Empty Input", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		input := product.Product{ID: uint(1), Name: "Asus TUF", Category: "Gaming"}
		repo.On("InsertProduct", uint(1), input).Return(product.Product{}, errors.New("Inputan tidak boleh kosong")).Once()

		products, err := productService.TalkToGpt(token, input)

		assert.Error(t, err)
		assert.Equal(t, product.Product{}, products)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Wrong Roles", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		input := product.Product{ID: uint(1), Name: "Asus TUF", Category: "Gaming"}

		products, err := productService.TalkToGpt(token, input)

		assert.Error(t, err)
		assert.Equal(t, product.Product{}, products)

	})
}
