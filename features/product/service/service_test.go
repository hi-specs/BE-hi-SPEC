package service_test

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/product/mocks"
	"BE-hi-SPEC/features/product/service"
	"BE-hi-SPEC/helper/jwt"
	"errors"
	"testing"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var userID = uint(1)
var str, _ = jwt.GenerateJWT(userID)
var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3yyy!!!"), nil
})

func TestSemuaProduct(t *testing.T) {
	mockRepo := new(mocks.Repository)
	productService := service.New(mockRepo)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		mockProducts := []product.Product{
			{ID: 1, Name: "product1"},
			{ID: 2, Name: "product2"},
		}
		mockRepo.On("GetAllProduct", mock.Anything, mock.Anything).Return(mockProducts, nil).Once()

		// Call the method being tested.
		result, err := productService.SemuaProduct(1, 10)

		// Assert that the result and error match the expectations.
		assert.Nil(t, err)
		assert.Equal(t, mockProducts, result)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the repository to simulate an error.
		mockRepo.On("GetAllProduct", mock.Anything, mock.Anything).Return(nil, errors.New("database error")).Once()

		// Call the method being tested.
		_, err := productService.SemuaProduct(1, 10)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed get all product", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Products Not Found", func(t *testing.T) {
		// Mock the repository to simulate "not found" error.
		mockRepo.On("GetAllProduct", mock.Anything, mock.Anything).Return([]product.Product{}, errors.New("not found")).Once()

		// Call the method being tested.
		_, err := productService.SemuaProduct(1, 10)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed get all product", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}

func TestDelProduct(t *testing.T) {
	mockRepo := new(mocks.Repository)
	productService := service.New(mockRepo)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		productID := uint(1)
		mockRepo.On("DelProduct", productID).Return(nil).Once()

		// Call the method being tested.
		err := productService.DelProduct(productID)

		// Assert that the error is nil, indicating success.
		assert.Nil(t, err)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the repository to simulate an error.
		productID := uint(1)
		mockRepo.On("DelProduct", productID).Return(errors.New("database error")).Once()

		// Call the method being tested.
		err := productService.DelProduct(productID)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}

func TestSatuProduct(t *testing.T) {
	mockRepo := new(mocks.Repository)
	productService := service.New(mockRepo)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		productID := uint(1)
		mockProduct := &product.Product{ID: productID, Name: "TestProduct"}
		mockRepo.On("GetProductID", productID).Return(mockProduct, nil).Once()

		// Call the method being tested.
		result, err := productService.SatuProduct(productID)

		// Assert that the result and error match the expectations.
		assert.Nil(t, err)
		assert.Equal(t, *mockProduct, result)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the repository to simulate an error.
		productID := uint(1)
		mockRepo.On("GetProductID", productID).Return(nil, errors.New("database error")).Once()

		// Call the method being tested.
		_, err := productService.SatuProduct(productID)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed get product", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Product Not Found", func(t *testing.T) {
		// Mock the repository to simulate "not found" error.
		productID := uint(1)
		mockRepo.On("GetProductID", productID).Return(nil, errors.New("not found")).Once()

		// Call the method being tested.
		_, err := productService.SatuProduct(productID)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed get product", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}

func TestCariProduct(t *testing.T) {
	mockRepo := new(mocks.Repository)
	productService := service.New(mockRepo)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		mockProducts := []product.Product{
			{ID: 1, Name: "Product1", Category: "Multimedia", Price: 25000000},
			{ID: 2, Name: "Product2", Category: "Gaming", Price: 10000000},
		}
		mockRepo.On("SearchProduct", "queryName", "queryCategory", uint(100), uint(1000), 1, 10).Return(mockProducts, nil).Once()

		// Call the method being tested.
		result, err := productService.CariProduct("queryName", "queryCategory", 100, 1000, 1, 10)

		// Assert that the result and error match the expectations.
		assert.Nil(t, err)
		assert.Equal(t, mockProducts, result)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the repository to simulate an error.
		mockRepo.On("SearchProduct", "queryName", "queryCategory", uint(100), uint(1000), 1, 10).Return(nil, errors.New("repository error")).Once()

		// Call the method being tested.
		_, err := productService.CariProduct("queryName", "queryCategory", 100, 1000, 1, 10)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateProduct(t *testing.T) {
	mockRepo := new(mocks.Repository)
	productService := service.New(mockRepo)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		productID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Category: "Multimedia", Price: 20000000}
		mockRepo.On("UpdateProduct", productID, mockInput).Return(mockInput, nil).Once()

		// Call the method being tested.
		result, err := productService.UpdateProduct(productID, mockInput)

		// Assert that the result and error match the expectations.
		assert.Nil(t, err)
		assert.Equal(t, mockInput, result)

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the repository to simulate an error.
		productID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Category: "Multimedia", Price: 20000000}
		mockRepo.On("UpdateProduct", productID, mockInput).Return(product.Product{}, errors.New("repository error")).Once()

		// Call the method being tested.
		_, err := productService.UpdateProduct(productID, mockInput)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed to update the product", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Case - Product Not Found", func(t *testing.T) {
		// Mock the repository to simulate "not found" error.
		productID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Category: "Electronics", Price: 1500}
		mockRepo.On("UpdateProduct", productID, mockInput).Return(product.Product{}, errors.New("not found")).Once()

		// Call the method being tested.
		_, err := productService.UpdateProduct(productID, mockInput)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed to update the product", err.Error())

		// Assert that the expected repository method was called.
		mockRepo.AssertExpectations(t)
	})
}
