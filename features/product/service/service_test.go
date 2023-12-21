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
)

var userID = uint(1)
var str, _ = jwt.GenerateJWT(userID)
var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3yyy!!!"), nil
})

func TestUpdateProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockProductID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Price: 1000}
		mockUpdatedProduct := product.Product{ID: mockProductID, Name: "UpdatedProduct", Price: 1000}
		repo.On("UpdateProduct", mockProductID, mockInput).Return(mockUpdatedProduct, nil).Once()

		result, err := productService.UpdateProduct(mockProductID, mockInput)

		assert.Nil(t, err)
		assert.Equal(t, mockUpdatedProduct, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Product Not Found", func(t *testing.T) {
		mockProductID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Price: 1000}
		repo.On("UpdateProduct", mockProductID, mockInput).Return(product.Product{}, errors.New("not found")).Once()

		_, err := productService.UpdateProduct(mockProductID, mockInput)

		assert.Error(t, err)
		assert.Equal(t, "failed to update the product", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockProductID := uint(1)
		mockInput := product.Product{Name: "UpdatedProduct", Price: 1000}
		repo.On("UpdateProduct", mockProductID, mockInput).Return(product.Product{}, errors.New("database error")).Once()

		_, err := productService.UpdateProduct(mockProductID, mockInput)

		assert.Error(t, err)
		assert.Equal(t, "failed to update the product", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestDelProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockProductID := uint(1)
		repo.On("DelProduct", mockProductID).Return(nil).Once()

		err := productService.DelProduct(mockProductID)

		assert.Nil(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Product Not Found", func(t *testing.T) {
		mockProductID := uint(1)
		repo.On("DelProduct", mockProductID).Return(errors.New("product not found")).Once()

		err := productService.DelProduct(mockProductID)

		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockProductID := uint(1)
		repo.On("DelProduct", mockProductID).Return(errors.New("database error")).Once()

		err := productService.DelProduct(mockProductID)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestSemuaProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockProducts := []product.Product{
			{ID: 1, Name: "Product1"},
			{ID: 2, Name: "Product2"},
		}
		mockTotalPage := 2
		repo.On("GetAllProduct", 1, 10).Return(mockProducts, mockTotalPage, nil).Once()

		result, totalPage, err := productService.SemuaProduct(1, 10)

		assert.Nil(t, err)
		assert.Equal(t, mockProducts, result)
		assert.Equal(t, mockTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		repo.On("GetAllProduct", 1, 10).Return(nil, 0, errors.New("failed get all product")).Once()

		_, _, err := productService.SemuaProduct(1, 10)

		assert.Error(t, err)
		assert.Equal(t, "failed get all product", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestSatuProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockProductID := uint(1)
		mockProduct := &product.Product{ID: mockProductID, Name: "MockProduct"}
		repo.On("GetProductID", mockProductID).Return(mockProduct, nil).Once()

		result, err := productService.SatuProduct(mockProductID)

		assert.Nil(t, err)
		assert.Equal(t, *mockProduct, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Product Not Found", func(t *testing.T) {
		mockProductID := uint(1)
		repo.On("GetProductID", mockProductID).Return(nil, errors.New("failed get product")).Once()

		_, err := productService.SatuProduct(mockProductID)

		assert.Error(t, err)
		assert.Equal(t, "failed get product", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockProductID := uint(1)
		repo.On("GetProductID", mockProductID).Return(nil, errors.New("failed get product")).Once()

		_, err := productService.SatuProduct(mockProductID)

		assert.Error(t, err)
		assert.Equal(t, "failed get product", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestCariProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockName := "MockProduct"
		mockCategory := "Multimedia"
		mockMinPrice := "100"
		mockMaxPrice := "500"
		mockPage := 1
		mockLimit := 10
		mockProducts := []product.Product{
			{ID: 1, Name: "Product1", Category: "Multimedia", Price: 20000000},
			{ID: 2, Name: "Product2", Category: "Multimedia", Price: 29000000},
		}
		mockTotalPage := 2
		repo.On("SearchProduct", mockName, mockCategory, mockMinPrice, mockMaxPrice, mockPage, mockLimit).
			Return(mockProducts, mockTotalPage, nil).Once()

		result, totalPage, err := productService.CariProduct(mockName, mockCategory, mockMinPrice, mockMaxPrice, mockPage, mockLimit)

		assert.Nil(t, err)
		assert.Equal(t, mockProducts, result)
		assert.Equal(t, mockTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockName := "MockProduct"
		mockCategory := "Multimedia"
		mockMinPrice := "100"
		mockMaxPrice := "500"
		mockPage := 1
		mockLimit := 10
		repo.On("SearchProduct", mockName, mockCategory, mockMinPrice, mockMaxPrice, mockPage, mockLimit).
			Return(nil, 0, errors.New("database error")).Once()

		_, _, err := productService.CariProduct(mockName, mockCategory, mockMinPrice, mockMaxPrice, mockPage, mockLimit)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Products Not Found", func(t *testing.T) {
		mockName := "NonExistentProduct"
		mockCategory := "Multimedia"
		mockMinPrice := "100"
		mockMaxPrice := "500"
		mockPage := 1
		mockLimit := 10
		repo.On("SearchProduct", mockName, mockCategory, mockMinPrice, mockMaxPrice, mockPage, mockLimit).
			Return(nil, 0, errors.New("products not found")).Once()

		_, _, err := productService.CariProduct(mockName, mockCategory, mockMinPrice, mockMaxPrice, mockPage, mockLimit)

		assert.Error(t, err)
		assert.Equal(t, "products not found", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestTalkToGpt(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		mockUserID := uint(1)
		mockNewProduct := product.Product{
			Name:      "TestProduct",
			Category:  "Test",
			CPU:       "Test",
			RAM:       "Test",
			Display:   "Test Display",
			Storage:   "Test Storage",
			Thickness: "Test Thickness",
			Weight:    "Test",
			Bluetooth: "Test",
			HDMI:      "Test",
			Price:     10000000,
			Picture:   "Picture.jpg",
		}

		repo.On("InsertProduct", mockUserID, mockNewProduct).Return(mockNewProduct, nil).Once()

		result, err := productService.TalkToGpt(token, mockNewProduct)

		assert.Nil(t, err)
		assert.Equal(t, mockNewProduct, result)

		repo.AssertExpectations(t)
	})

}
