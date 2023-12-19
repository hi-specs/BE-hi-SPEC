package service

import (
	"BE-hi-SPEC/features/transaction"
	"BE-hi-SPEC/helper/jwt"
	"errors"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type TransactionServices struct {
	repo transaction.Repository
}

func New(r transaction.Repository) transaction.Service {
	return &TransactionServices{
		repo: r,
	}
}

func (ts *TransactionServices) AdminDashboard() (transaction.TransactionDashboard, error) {
	result, err := ts.repo.AdminDashboard()
	return result, err
}

func (ts *TransactionServices) Checkout(token *golangjwt.Token, ProductID int, TotalPrice int) (transaction.Transaction, error) {
	userID, err := jwt.ExtractToken(token)
	if err != nil {
		return transaction.Transaction{}, errors.New("user does not exist")
	}

	result, err := ts.repo.Checkout(userID, int(ProductID), TotalPrice)
	return result, err
}

func (ts *TransactionServices) TransactionList(page, limit int) ([]transaction.TransactionList, int, error) {
	result, totalPage, err := ts.repo.TransactionList(page, limit)
	return result, totalPage, err
}

func (ts *TransactionServices) GetTransaction(transactionID uint) (transaction.TransactionList, error) {
	result, err := ts.repo.GetTransaction(transactionID)
	return *result, err
}

func (ts *TransactionServices) MidtransCallback(transactionID string) (transaction.TransactionList, error) {
	result, err := ts.repo.MidtransCallback(transactionID)
	return *result, err
}
