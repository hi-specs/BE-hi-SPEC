package service_test

import (
	"errors"
	"testing"

	"BE-hi-SPEC/features/product"

	"github.com/stretchr/testify/assert"

	"BE-hi-SPEC/features/product/mocks"
	"BE-hi-SPEC/features/product/service"
)

func TestUpdateProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		// Mock the expected behavior of the repository.
		mockProductID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Price: 1000}
		mockUpdatedProduct := product.Product{ID: mockProductID, Name: "UpdatedProduct", Price: 1000}
		repo.On("UpdateProduct", mockProductID, mockInput).Return(mockUpdatedProduct, nil).Once()

		// Call the method being tested.
		result, err := productService.UpdateProduct(mockProductID, mockInput)

		// Assert that the result and error match the expectations.
		assert.Nil(t, err)
		assert.Equal(t, mockUpdatedProduct, result)

		// Assert that the expected repository method was called.
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Product Not Found", func(t *testing.T) {
		// Mock the repository to simulate "not found" error.
		mockProductID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Price: 1000}
		repo.On("UpdateProduct", mockProductID, mockInput).Return(product.Product{}, errors.New("not found")).Once()

		// Call the method being tested.
		_, err := productService.UpdateProduct(mockProductID, mockInput)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed to update the product", err.Error())

		// Assert that the expected repository method was called.
		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		// Mock the repository to simulate a general error.
		mockProductID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Price: 1000}
		repo.On("UpdateProduct", mockProductID, mockInput).Return(product.Product{}, errors.New("database error")).Once()

		// Call the method being tested.
		_, err := productService.UpdateProduct(mockProductID, mockInput)

		// Assert that the error matches the expectations.
		assert.Error(t, err)
		assert.Equal(t, "failed to update the product", err.Error())

		// Assert that the expected repository method was called.
		repo.AssertExpectations(t)
	})
}
