package service_test

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/features/transaction"
	"BE-hi-SPEC/features/transaction/mocks"
	"BE-hi-SPEC/features/transaction/service"
	"BE-hi-SPEC/features/user"
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

func TestAdminDashboard(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)
	t.Run("Success Case", func(t *testing.T) {
		mockDashboard := transaction.TransactionDashboard{
			TotalProduct:     100,
			TotalUser:        5000.0,
			TotalTransaction: 90,
		}
		repo.On("AdminDashboard").Return(mockDashboard, nil).Once()

		result, err := transactionService.AdminDashboard()

		assert.Nil(t, err)
		assert.Equal(t, mockDashboard, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		repo.On("AdminDashboard").Return(transaction.TransactionDashboard{}, errors.New("failed to fetch admin dashboard data")).Once()

		_, err := transactionService.AdminDashboard()

		assert.Error(t, err)
		assert.Equal(t, "failed to fetch admin dashboard data", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestCheckout(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		mockUserID := uint(1)
		mockProductID := 123
		mockTotalPrice := 1000
		mockTransaction := transaction.Transaction{
			ID:         1,
			ProductID:  mockProductID,
			TotalPrice: mockTotalPrice,
			Status:     "success",
		}

		repo.On("Checkout", mockUserID, mockProductID, mockTotalPrice).Return(mockTransaction, nil).Once()

		result, err := transactionService.Checkout(token, mockProductID, mockTotalPrice)

		assert.Nil(t, err)
		assert.Equal(t, mockTransaction, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		var userID = uint(1)
		var str, _ = jwt.GenerateJWT(userID)
		var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
			return []byte("$!1gnK3yyy!!!"), nil
		})
		mockUserID := uint(1)
		mockProductID := 123
		mockTotalPrice := 1000
		mockErr := errors.New("failed to process checkout")

		repo.On("Checkout", mockUserID, mockProductID, mockTotalPrice).Return(transaction.Transaction{}, mockErr).Once()

		_, err := transactionService.Checkout(token, mockProductID, mockTotalPrice)

		assert.Error(t, err)
		assert.Equal(t, "failed to process checkout", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestTransactionList(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockPage := 1
		mockLimit := 10
		mockTransactionList := []transaction.TransactionList{
			{Nota: "HI-9210371", ProductID: 123, TotalPrice: 1000, Status: "success"},
			// Add more transaction list items as needed
		}
		mockTotalPage := 2

		repo.On("TransactionList", mockPage, mockLimit).Return(mockTransactionList, mockTotalPage, nil).Once()

		result, totalPage, err := transactionService.TransactionList(mockPage, mockLimit)

		assert.Nil(t, err)
		assert.Equal(t, mockTransactionList, result)
		assert.Equal(t, mockTotalPage, totalPage)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockPage := 1
		mockLimit := 10
		mockErr := errors.New("failed to fetch transaction list")

		repo.On("TransactionList", mockPage, mockLimit).Return(nil, 0, mockErr).Once()

		_, _, err := transactionService.TransactionList(mockPage, mockLimit)

		assert.Error(t, err)
		assert.Equal(t, "failed to fetch transaction list", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestGetTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockTransactionID := uint(1)
		mockTransaction := transaction.TransactionList{
			Nota:       "HI-28916398",
			ProductID:  123,
			TotalPrice: 25000000,
			Status:     "success",
		}

		repo.On("GetTransaction", mockTransactionID).Return(&mockTransaction, nil).Once()

		result, err := transactionService.GetTransaction(mockTransactionID)

		assert.Nil(t, err)
		assert.Equal(t, mockTransaction, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Transaction Not Found", func(t *testing.T) {
		mockTransactionID := uint(1)
		repo.On("GetTransaction", mockTransactionID).Return(nil, errors.New("transaction not found")).Once()

		_, err := transactionService.GetTransaction(mockTransactionID)

		assert.Error(t, err)
		assert.Equal(t, "transaction not found", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockTransactionID := uint(1)
		mockErr := errors.New("failed to fetch transaction")

		repo.On("GetTransaction", mockTransactionID).Return(nil, mockErr).Once()

		_, err := transactionService.GetTransaction(mockTransactionID)

		assert.Error(t, err)
		assert.Equal(t, "failed to fetch transaction", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestMidtransCallback(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockTransactionID := "midtrans-123"
		mockTransaction := transaction.TransactionList{
			Nota:       "HI-28916398",
			ProductID:  123,
			TotalPrice: 25000000,
			Status:     "success",
		}

		repo.On("MidtransCallback", mockTransactionID).Return(&mockTransaction, nil).Once()

		result, err := transactionService.MidtransCallback(mockTransactionID)

		assert.Nil(t, err)
		assert.Equal(t, mockTransaction, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Transaction Not Found", func(t *testing.T) {
		mockTransactionID := "midtrans-123"
		repo.On("MidtransCallback", mockTransactionID).Return(nil, errors.New("transaction not found")).Once()

		_, err := transactionService.MidtransCallback(mockTransactionID)

		assert.Error(t, err)
		assert.Equal(t, "transaction not found", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockTransactionID := "midtrans-123"
		mockErr := errors.New("repository error")

		repo.On("MidtransCallback", mockTransactionID).Return(nil, mockErr).Once()

		_, err := transactionService.MidtransCallback(mockTransactionID)

		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())

		repo.AssertExpectations(t)
	})
}

func TestUserTransaction(t *testing.T) {
	repo := mocks.NewRepository(t)
	transactionService := service.New(repo)

	t.Run("Success Case", func(t *testing.T) {
		mockUserID := uint(1)
		mockUserTransaction := transaction.UserTransaction{
			User:        user.User{},
			Product:     []product.Product{},
			Transaction: []transaction.Transaction{},
		}

		repo.On("UserTransaction", mockUserID).Return(mockUserTransaction, nil).Once()

		result, err := transactionService.UserTransaction(mockUserID)

		assert.Nil(t, err)
		assert.Equal(t, mockUserTransaction, result)

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - User Not Found", func(t *testing.T) {
		mockUserID := uint(1)
		repo.On("UserTransaction", mockUserID).Return(transaction.UserTransaction{}, errors.New("user not found")).Once()

		_, err := transactionService.UserTransaction(mockUserID)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())

		repo.AssertExpectations(t)
	})

	t.Run("Error Case - Repository Failure", func(t *testing.T) {
		mockUserID := uint(1)
		mockErr := errors.New("repository error")

		repo.On("UserTransaction", mockUserID).Return(transaction.UserTransaction{}, mockErr).Once()

		_, err := transactionService.UserTransaction(mockUserID)

		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())

		repo.AssertExpectations(t)
	})
}
