package service_test

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/transaction"
	"BE-hi-SPEC/features/transaction/mocks"
	"BE-hi-SPEC/features/transaction/service"
	"BE-hi-SPEC/features/user"
	"errors"
	"strconv"
	"testing"

	golangjwt "BE-hi-SPEC/helper/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAdminDashboard(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		expectedDashboard := transaction.TransactionDashboard{TotalProduct: 5, TotalUser: 2, TotalTransaction: 5}
		expectedTotalPage := 3
		repo.On("AdminDashboard", uint(1), 1, 10).Return(expectedDashboard, expectedTotalPage, nil).Once()
		transactions, totalPage, err := transactionService.AdminDashboard(token, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedDashboard, transactions)
		assert.Equal(t, expectedTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Token Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		transactions, totalPage, err := transactionService.AdminDashboard(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionDashboard{}, transactions)
		assert.Equal(t, 0, totalPage)
	})

	t.Run("Failed Case Admin Role Required", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		transactions, totalPage, err := transactionService.AdminDashboard(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionDashboard{}, transactions)
		assert.Equal(t, 0, totalPage)
	})

	t.Run("Failed Case Repository Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		repo.On("AdminDashboard", uint(1), 1, 10).Return(transaction.TransactionDashboard{}, 0, errors.New("Failed get All data Dashboard")).Once()

		transactions, totalPage, err := transactionService.AdminDashboard(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionDashboard{}, transactions)
		assert.Equal(t, 0, totalPage)
	})
}

func TestCheckout(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		input := transaction.Transaction{ProductID: 3, TotalPrice: 19000000}
		repo.On("Checkout", uint(2), 3, 19000000).Return(input, nil).Once()

		transactions, err := transactionService.Checkout(token, 3, 19000000)
		assert.NoError(t, err)
		assert.Equal(t, input, transactions)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Token Error", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		transactions, err := transactionService.Checkout(token, 3, 19000000)
		assert.Error(t, err)
		assert.Equal(t, transaction.Transaction{}, transactions)
	})

	t.Run("Failed Case Roles Empty", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = ""
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		transactions, err := transactionService.Checkout(token, 3, 19000000)
		assert.Error(t, err)
		assert.Equal(t, transaction.Transaction{}, transactions)
	})

	t.Run("Failed Case Repository Error", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		repo.On("Checkout", uint(2), 3, 19000000).Return(transaction.Transaction{}, errors.New("Repository Error")).Once()

		transactions, err := transactionService.Checkout(token, 3, 19000000)
		assert.Error(t, err)
		assert.Equal(t, transaction.Transaction{}, transactions)
	})
}

func TestTransactionList(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		expectedTransaction := []transaction.TransactionList{{TransactionID: 1, ProductID: 3, TotalPrice: 2000000, Status: "Payment", Token: "wdawhdi", Url: "www.sandbox.com", Nota: "Hi-Spec 001"}, {TransactionID: 2, ProductID: 4, TotalPrice: 15000000, Status: "Payment", Token: "wdawhdi", Url: "www.sandbox.com", Nota: "Hi-Spec 002"}}
		expectedTotalPage := 3
		repo.On("TransactionList", uint(1), 1, 10).Return(expectedTransaction, expectedTotalPage, nil).Once()
		transactions, totalPage, err := transactionService.TransactionList(token, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction, transactions)
		assert.Equal(t, expectedTotalPage, totalPage)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Token error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		transactions, totalPage, err := transactionService.TransactionList(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []transaction.TransactionList{}, transactions)
		assert.Equal(t, 0, totalPage)
	})

	t.Run("Failed Case Admin Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		transactions, totalPage, err := transactionService.TransactionList(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []transaction.TransactionList{}, transactions)
		assert.Equal(t, 0, totalPage)
	})

	t.Run("Failed Case Repository Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		repo.On("TransactionList", uint(1), 1, 10).Return([]transaction.TransactionList{}, 0, errors.New("Repository error")).Once()

		transactions, totalPage, err := transactionService.TransactionList(token, 1, 10)
		assert.Error(t, err)
		assert.Equal(t, []transaction.TransactionList{}, transactions)
		assert.Equal(t, 0, totalPage)
	})
}

func TestGetTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		transactionID := uint(1)
		expectedTransaction := &transaction.TransactionList{
			TransactionID: int(transactionID),
			ProductID:     3,
			TotalPrice:    20000000,
			Status:        "Success",
			Token:         "dwamhlodho",
			Url:           "www.sandbox.com",
			Nota:          "Hi-SPEC 001",
		}
		repo.On("GetTransaction", uint(1), uint(1)).Return(expectedTransaction, nil).Once()
		transactions, err := transactionService.GetTransaction(token, transactionID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction, &transactions)

		repo.AssertExpectations(t)

	})

	t.Run("Failed Case Token Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		transactionID := uint(1)

		transactions, err := transactionService.GetTransaction(token, transactionID)
		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionList{}, transactions)

	})

	t.Run("Failed Case Admin Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		transactionID := uint(1)

		transactions, err := transactionService.GetTransaction(token, transactionID)
		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionList{}, transactions)
	})

	t.Run("Failed Case Repository Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		transactionID := uint(1)
		repo.On("GetTransaction", uint(1), uint(1)).Return(&transaction.TransactionList{}, errors.New("Repository error")).Once()

		transactions, err := transactionService.GetTransaction(token, transactionID)
		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionList{}, transactions)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Transaction not Found", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		transactionID := uint(1)
		repo.On("GetTransaction", userID, transactionID).Return(&transaction.TransactionList{}, errors.New("transaction not found")).Once()

		transactions, err := transactionService.GetTransaction(token, transactionID)
		assert.Error(t, err, errors.New("transaction not found"))
		assert.Equal(t, transactions, transaction.TransactionList{})
		repo.AssertExpectations(t)
	})
}

func TestMidTransCallback(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		transactionID := 1
		expectedTransaction := &transaction.TransactionList{
			TransactionID: 1,
			ProductID:     3,
			TotalPrice:    20000000,
			Status:        "Success",
			Token:         "dwamhlodho",
			Url:           "www.sandbox.com",
			Nota:          "Hi-SPEC 001",
		}
		repo.On("MidtransCallback", strconv.Itoa(transactionID)).Return(expectedTransaction, nil).Once()
		transactions, err := transactionService.MidtransCallback(strconv.Itoa(transactionID))
		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction, &transactions)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Repository Error", func(t *testing.T) {
		transactionID := 1
		repo.On("MidtransCallback", strconv.Itoa(transactionID)).Return(&transaction.TransactionList{}, errors.New("Errors")).Once()
		transactions, err := transactionService.MidtransCallback(strconv.Itoa(transactionID))
		assert.Error(t, err)
		assert.Equal(t, transaction.TransactionList{}, transactions)

	})
}

func TestUserTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		expectedTransaction := transaction.UserTransaction{
			User:        user.User{},
			Product:     []product.Product{},
			Transaction: []transaction.Transaction{},
		}
		repo.On("UserTransaction", int(userID), uint(1)).Return(expectedTransaction, nil).Once()
		transactions, err := transactionService.UserTransaction(token, userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction, transactions)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case User not Found", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		transactions, err := transactionService.UserTransaction(token, userID)
		assert.Error(t, err)
		assert.Equal(t, transaction.UserTransaction{}, transactions)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Admin Required", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		transactions, err := transactionService.UserTransaction(token, userID)
		assert.Error(t, err)
		assert.Equal(t, transaction.UserTransaction{}, transactions)

		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Repository Error", func(t *testing.T) {
		var userID = uint(1)
		var rolesUser = "admin"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})

		repo.On("UserTransaction", int(userID), uint(1)).Return(transaction.UserTransaction{}, errors.New("Repository error")).Once()
		transactions, err := transactionService.UserTransaction(token, userID)
		assert.Error(t, err)
		assert.Equal(t, transaction.UserTransaction{}, transactions)

		repo.AssertExpectations(t)
	})
}

func TestDownloadTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		transactionID := uint(2)
		repo.On("DownloadTransaction", userID, transactionID).Return(nil).Once()
		err := transactionService.DownloadTransaction(token, transactionID)

		assert.Nil(t, err)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Token Error", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		transactionID := uint(2)
		err := transactionService.DownloadTransaction(token, transactionID)

		assert.Error(t, err)
	})

	t.Run("Failed Case Repository Error", func(t *testing.T) {
		var userID = uint(2)
		var rolesUser = "user"
		var str, _ = golangjwt.GenerateJWT(userID, rolesUser)
		var token, _ = jwt.Parse(str, func(t *jwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		transactionID := uint(2)
		repo.On("DownloadTransaction", userID, transactionID).Return(errors.New("Repository error")).Once()
		err := transactionService.DownloadTransaction(token, transactionID)

		assert.Error(t, err)
	})

}

func TestUpdatePdfTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		transactionID := uint(2)
		link := "www.cloudinary.com"
		repo.On("UpdatePdfTransaction", link, transactionID).Return(nil).Once()
		err := transactionService.UpdatePdfTransaction(link, transactionID)

		assert.Nil(t, err)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Case Repository Error", func(t *testing.T) {
		transactionID := uint(2)
		link := "www.cloudinary.com"
		repo.On("UpdatePdfTransaction", link, transactionID).Return(errors.New("Repository Error")).Once()
		err := transactionService.UpdatePdfTransaction(link, transactionID)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

}
