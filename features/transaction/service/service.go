package service

import "BE-hi-SPEC/features/transaction"

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
