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

func (ts *TransactionServices) AdminDashboard(token *golangjwt.Token, page int, limit int) (transaction.TransactionDashboard, int, error) {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return transaction.TransactionDashboard{}, 0, errors.New("token error")
	}
	if rolesUser != "admin" {
		return transaction.TransactionDashboard{}, 0, errors.New("unauthorized access: admin role required")
	}
	result, totalPage, err := ts.repo.AdminDashboard(userId, page, limit)
	if err != nil {
		return transaction.TransactionDashboard{}, 0, errors.New("Failed get All data Dashboard")
	}
	return result, totalPage, err
}

func (ts *TransactionServices) Checkout(token *golangjwt.Token, ProductID int, TotalPrice int) (transaction.Transaction, error) {
	userID, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return transaction.Transaction{}, errors.New("user does not exist")
	}
	if rolesUser == "" {
		return transaction.Transaction{}, errors.New("roles user can't empty")
	}
	result, err := ts.repo.Checkout(userID, int(ProductID), TotalPrice)
	if err != nil {
		return transaction.Transaction{}, errors.New("Repository Error")
	}
	return result, err
}

func (ts *TransactionServices) TransactionList(token *golangjwt.Token, page, limit int) ([]transaction.TransactionList, int, error) {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return []transaction.TransactionList{}, 0, errors.New("user does not exist")
	}
	if rolesUser != "admin" {
		return []transaction.TransactionList{}, 0, errors.New("you are not authorized")
	}
	result, totalPage, err := ts.repo.TransactionList(uint(userId), page, limit)
	if err != nil {
		return []transaction.TransactionList{}, 0, errors.New("Repository error")
	}
	return result, totalPage, err
}

func (ts *TransactionServices) GetTransaction(token *golangjwt.Token, transactionID uint) (transaction.TransactionList, error) {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return transaction.TransactionList{}, errors.New("user does not exist")
	}
	if rolesUser != "admin" {
		return transaction.TransactionList{}, errors.New("unauthorized access: admin role required")
	}
	result, err := ts.repo.GetTransaction(userId, transactionID)
	if err != nil {
		return transaction.TransactionList{}, errors.New("Repository error")
	}

	if result == nil {
		return transaction.TransactionList{}, errors.New("transaction not found")
	}

	return *result, nil
}

func (ts *TransactionServices) MidtransCallback(transactionID string) (transaction.TransactionList, error) {
	result, err := ts.repo.MidtransCallback(transactionID)
	if err != nil {
		return transaction.TransactionList{}, errors.New("Errors")
	}

	if result == nil {
		return transaction.TransactionList{}, errors.New("transaction not found")
	}

	return *result, nil
}

func (ts *TransactionServices) UserTransaction(token *golangjwt.Token, userID uint) (transaction.UserTransaction, error) {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return transaction.UserTransaction{}, errors.New("user does not exist")
	}
	if rolesUser != "admin" {
		return transaction.UserTransaction{}, errors.New("unauthorized access: admin role required")
	}
	result, err := ts.repo.UserTransaction(int(userId), userID)
	if err != nil {
		return transaction.UserTransaction{}, errors.New("Repository error")
	}
	return result, err
}

func (ts *TransactionServices) DownloadTransaction(token *golangjwt.Token, transactionID uint) error {
	userID, _, err := jwt.ExtractToken(token)
	if err != nil {
		return errors.New("Token error")
	}
	err2 := ts.repo.DownloadTransaction(userID, transactionID)
	if err2 != nil {
		return errors.New("Repository error")
	}
	return nil
}

func (ts *TransactionServices) UpdatePdfTransaction(link string, transactionID uint) error {
	err := ts.repo.UpdatePdfTransaction(link, transactionID)
	if err != nil {
		return errors.New("Repository error")
	}
	return nil
}
