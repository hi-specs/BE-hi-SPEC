package repository

import (
	p "BE-hi-SPEC/features/product"
	pr "BE-hi-SPEC/features/product/repository"

	"BE-hi-SPEC/features/transaction"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type TransactionModel struct {
	gorm.Model
	ProductID  uint
	UserID     uint
	TotalPrice uint
	Status     string
}

type TransactionQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) transaction.Repository {
	return &TransactionQuery{
		db: db,
	}
}

func (tq *TransactionQuery) AdminDashboard() (transaction.TransactionDashboard, error) {
	// mendapatkan nilai total product
	var productCount int
	tableName := "product_models"
	columnName := "created_at"
	query := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NOT NULL", tableName, columnName)
	err := tq.db.Raw(query).Scan(&productCount).Error
	if err != nil {
		log.Fatal(err)
	}

	// mendapatkan nilai total user
	var userCount int
	tableNameUser := "user_models"
	columnNameUser := "created_at"
	queryuser := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NOT NULL", tableNameUser, columnNameUser)
	err2 := tq.db.Raw(queryuser).Scan(&userCount).Error
	if err2 != nil {
		log.Fatal(err)
	}

	// mendapatkan nilai total transaksi yang sukses
	var transactionCount int
	tableNameTransaction := "transaction_models"
	columnNameTransaction := "deleted_at"
	querytransaction := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NOT NULL", tableNameTransaction, columnNameTransaction)
	err3 := tq.db.Raw(querytransaction).Scan(&transactionCount).Error
	if err3 != nil {
		log.Fatal(err)
	}

	var products []pr.ProductModel
	if err4 := tq.db.Find(&products).Error; err != nil {
		return transaction.TransactionDashboard{}, err4
	}
	var prod []p.Product
	for _, s := range products {
		prod = append(prod, p.Product{
			ID:       s.ID,
			Name:     s.Name,
			Price:    s.Price,
			Category: s.Category,
			Picture:  s.Picture,
		})
	}

	// memasukan nilai yang didapat dari DB kedalam struct
	var result = new(transaction.TransactionDashboard)
	result.TotalProduct = productCount
	result.TotalUser = userCount
	result.TotalTransaction = transactionCount
	result.Product = prod

	return *result, err
}

func (tq *TransactionQuery) Checkout(userID uint, ProductID int, TotalPrice int) (transaction.Transaction, error) {
	var inputDB TransactionModel
	inputDB.ProductID = uint(ProductID)
	inputDB.TotalPrice = uint(TotalPrice)
	inputDB.UserID = userID
	inputDB.Status = "Pending"

	if err := tq.db.Create(&inputDB).Error; err != nil {
		return transaction.Transaction{}, err
	}

	var result transaction.Transaction
	result.ID = int(inputDB.ID)
	result.ProductID = int(inputDB.ProductID)
	result.TotalPrice = int(inputDB.TotalPrice)
	result.Status = inputDB.Status

	return result, nil
}

func (tq *TransactionQuery) TransactionList() ([]transaction.TransactionList, error) {
	var tm []TransactionModel
	if err := tq.db.Find(&tm).Error; err != nil {
		return nil, err
	}

	var result []transaction.TransactionList

	for _, tl := range tm {
		result = append(result, transaction.TransactionList{
			ProductID:     int(tl.ProductID),
			TransactionID: int(tl.ID),
			TotalPrice:    int(tl.TotalPrice),
			Status:        tl.Status,
			Timestamp:     tl.CreatedAt,
		})
	}

	return result, nil
}
