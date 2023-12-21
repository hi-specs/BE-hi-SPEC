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
	userID, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return transaction.Transaction{}, errors.New("user does not exist")
	}
	if rolesUser != "" {
		return transaction.Transaction{}, err
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
	if err != nil {
		return transaction.TransactionList{}, err
	}

	if result == nil {
		return transaction.TransactionList{}, errors.New("transaction not found")
	}

	return *result, nil
}

func (ts *TransactionServices) MidtransCallback(transactionID string) (transaction.TransactionList, error) {
	result, err := ts.repo.MidtransCallback(transactionID)
	if err != nil {
		return transaction.TransactionList{}, err
	}

	if result == nil {
		return transaction.TransactionList{}, errors.New("transaction not found")
	}

	return *result, nil
}

func (ts *TransactionServices) UserTransaction(userID uint) (transaction.UserTransaction, error) {
	result, err := ts.repo.UserTransaction(userID)
	return result, err
}
