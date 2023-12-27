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

func TestSemuaProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		expectedProduct := []product.Product{{ID: 1, Name: "User1", Price: 50000000}, {ID: 2, Name: "User2", Price: 60000000}}
		expectedTotalPage := 3

		repo.On("GetAllProduct", 1, 10).Return(expectedProduct, expectedTotalPage, nil).Once()
		products, totalPage, err := productService.SemuaProduct(1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, products)
		assert.Equal(t, expectedTotalPage, totalPage)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Error Repository", func(t *testing.T) {

		repo.On("GetAllProduct", 1, 10).Return([]product.Product{}, 0, errors.New("failed get all product")).Once()

		products, totalPage, err := productService.SemuaProduct(1, 10)
		assert.Error(t, err)
		assert.Equal(t, 0, totalPage)
		assert.Equal(t, []product.Product([]product.Product{}), products)
		repo.AssertExpectations(t)

	})
}

func TestSatuProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		productID := uint(1)
		expectedProduct := &product.Product{
			ID:        productID,
			Name:      "Asus TUF",
			Category:  "Gaming",
			CPU:       "AMD Ryzen 5",
			Display:   "14-inch Retina display",
			Storage:   "1 TB SSD",
			Thickness: "14.9 mm",
			Weight:    "1.37 kg",
			RAM:       "8GB",
			Bluetooth: "yes",
			HDMI:      "yes",
			Price:     50000000,
			Picture:   "test.jpg",
		}
		repo.On("GetProductID", uint(1)).Return(expectedProduct, nil).Once()
		products, err := productService.SatuProduct(productID)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, &products)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Error Repository", func(t *testing.T) {
		productID := uint(1)
		repo.On("GetProductID", uint(1)).Return(nil, errors.New("failed get product")).Once()

		products, err := productService.SatuProduct(productID)
		assert.Error(t, err)
		assert.Equal(t, product.Product{}, products)
		repo.AssertExpectations(t)

	})
}

func TestUpdateProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		input := product.Product{
			ID:        uint(1),
			Name:      "Asus TUF",
			Category:  "Gaming",
			CPU:       "AMD Ryzen 5",
			Display:   "14-inch Retina display",
			Storage:   "1 TB SSD",
			Thickness: "14.9 mm",
			Weight:    "1.37 kg",
			RAM:       "8GB",
			Bluetooth: "yes",
			HDMI:      "yes",
			Price:     50000000,
			Picture:   "test.jpg",
		}
		repo.On("UpdateProduct", uint(1), uint(1), input).Return(input, nil).Once()

		updatedProduct, err := productService.UpdateProduct(token, uint(1), input)
		assert.NoError(t, err)
		assert.Equal(t, input, updatedProduct)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Token Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		input := product.Product{
			ID:        uint(1),
			Name:      "Asus TUF",
			Category:  "Gaming",
			CPU:       "AMD Ryzen 5",
			Display:   "14-inch Retina display",
			Storage:   "1 TB SSD",
			Thickness: "14.9 mm",
			Weight:    "1.37 kg",
			RAM:       "8GB",
			Bluetooth: "yes",
			HDMI:      "yes",
			Price:     50000000,
			Picture:   "test.jpg",
		}
		updatedProduct, err := productService.UpdateProduct(token, uint(1), input)
		assert.Error(t, err)
		assert.Equal(t, product.Product{}, updatedProduct)
	})

	t.Run("Failed Case Empty Role", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = ""
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		input := product.Product{
			ID:        uint(1),
			Name:      "Asus TUF",
			Category:  "Gaming",
			CPU:       "AMD Ryzen 5",
			Display:   "14-inch Retina display",
			Storage:   "1 TB SSD",
			Thickness: "14.9 mm",
			Weight:    "1.37 kg",
			RAM:       "8GB",
			Bluetooth: "yes",
			HDMI:      "yes",
			Price:     50000000,
			Picture:   "test.jpg",
		}
		updatedProduct, err := productService.UpdateProduct(token, uint(1), input)
		assert.Error(t, err)
		assert.Equal(t, product.Product{}, updatedProduct)
	})

	t.Run("Failed Case Admin Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		input := product.Product{
			ID:        uint(1),
			Name:      "Asus TUF",
			Category:  "Gaming",
			CPU:       "AMD Ryzen 5",
			Display:   "14-inch Retina display",
			Storage:   "1 TB SSD",
			Thickness: "14.9 mm",
			Weight:    "1.37 kg",
			RAM:       "8GB",
			Bluetooth: "yes",
			HDMI:      "yes",
			Price:     50000000,
			Picture:   "test.jpg",
		}
		updatedProduct, err := productService.UpdateProduct(token, uint(1), input)
		assert.Error(t, err)
		assert.Equal(t, product.Product{}, updatedProduct)
	})

	t.Run("Failed Case Empty Input", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		input := product.Product{
			ID:        uint(1),
			Name:      "",
			Category:  "",
			CPU:       "",
			Display:   "",
			Storage:   "",
			Thickness: "",
			Weight:    "",
			RAM:       "",
			Bluetooth: "",
			HDMI:      "",
			Price:     0,
			Picture:   "",
		}
		repo.On("UpdateProduct", uint(1), uint(1), input).Return(product.Product{}, errors.New("failed to update the product")).Once()

		updatedProduct, err := productService.UpdateProduct(token, uint(1), input)
		assert.Error(t, err)
		assert.Equal(t, product.Product{}, updatedProduct)
	})
}

func TestDelProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		productID := uint(1)
		repo.On("DelProduct", userID, productID).Return(nil).Once()
		err := productService.DelProduct(token, productID)

		assert.Nil(t, err)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
	t.Run("Failed Case Token Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		productID := uint(1)
		err := productService.DelProduct(token, productID)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Admin Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		productID := uint(1)
		err := productService.DelProduct(token, productID)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestCariProduct(t *testing.T) {
	repo := mocks.NewRepository(t)
	productService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		expectedProduct := []product.Product{{ID: uint(1), Name: "Asus TUF", Price: 2000000, Picture: "tes1.jpg"}, {ID: uint(2), Name: "Macbook Pro M1", Price: 15000000, Picture: "tes2.jpg"}}
		expectedTotalPage := 3
		repo.On("SearchProduct", "Asus TUF", "Gaming", uint(10000000), uint(25000000), 1, 10).Return(expectedProduct, expectedTotalPage, nil).Once()
		product, totalPage, err := productService.CariProduct("Asus TUF", "Gaming", "10000000", "25000000", 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Empty minimal Price", func(t *testing.T) {
		var name = "Asus TUF"
		var category = "Gaming"
		repo.On("SearchProduct", name, category, uint(1), uint(25000000), 1, 10).Return(nil, 0, errors.New("tidak boleh kosong")).Once()
		products, totalPage, err := productService.CariProduct(name, category, "", "25000000", 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []product.Product([]product.Product(nil)), products)
		assert.Equal(t, 0, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Minimal Price must be Integer", func(t *testing.T) {
		var name = "Asus TUF"
		var category = "Gaming"
		products, totalPage, err := productService.CariProduct(name, category, "invalidMinPrice", "25000000", 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []product.Product{}, products)
		assert.Equal(t, 0, totalPage)

	})

	t.Run("Failed Case Maximal Price must be Integer", func(t *testing.T) {
		var name = "Asus TUF"
		var category = "Gaming"
		products, totalPage, err := productService.CariProduct(name, category, "10000000", "invalidMinPrice", 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []product.Product{}, products)
		assert.Equal(t, 0, totalPage)

	})

	t.Run("Failed Case Empty maximal Price", func(t *testing.T) {
		var name = "Asus TUF"
		var category = "Gaming"
		repo.On("SearchProduct", name, category, uint(10000000), uint(0), 1, 10).Return(nil, 0, errors.New("tidak boleh kosong")).Once()
		products, totalPage, err := productService.CariProduct(name, category, "10000000", "", 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []product.Product([]product.Product(nil)), products)
		assert.Equal(t, 0, totalPage)

		repo.AssertExpectations(t)
	})
}
