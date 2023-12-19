package repository

import (
	p "BE-hi-SPEC/features/product"
	pr "BE-hi-SPEC/features/product/repository"
	"BE-hi-SPEC/helper/midtrans"
	"errors"
	"strconv"

	"BE-hi-SPEC/features/transaction"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type TransactionModel struct {
	gorm.Model
	Nota       string
	ProductID  uint
	UserID     uint
	TotalPrice uint
	Status     string
	Token      string
	Url        string
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
	columnNameTransaction := "status"
	Status := "Success"
	querytransaction := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s = '%s'", tableNameTransaction, columnNameTransaction, Status)
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

	var id = strconv.Itoa(int(inputDB.ID))
	inputDB.Nota = "HI-SPEC-" + id

	midtrans := midtrans.MidtransCreateToken(int(inputDB.ID), TotalPrice)

	inputDB.Url = midtrans.RedirectURL
	inputDB.Token = midtrans.Token
	if err := tq.db.Save(&inputDB).Error; err != nil {
		return transaction.Transaction{}, err
	}

	var result transaction.Transaction
	result.ID = int(inputDB.ID)
	result.ProductID = int(inputDB.ProductID)
	result.TotalPrice = int(inputDB.TotalPrice)
	result.Status = inputDB.Status
	result.Nota = inputDB.Nota
	result.Token = midtrans.Token
	result.Url = midtrans.RedirectURL

	return result, nil
}

func (tq *TransactionQuery) TransactionList() ([]transaction.TransactionList, error) {
	var tm []TransactionModel

	err := tq.db.Find(&tm).Error

	var result []transaction.TransactionList
	for _, resp := range tm {
		results := transaction.TransactionList{
			TransactionID: int(resp.ID),
			ProductID:     int(resp.ProductID),
			TotalPrice:    int(resp.TotalPrice),
			Status:        resp.Status,
			Timestamp:     resp.CreatedAt,
			Token:         resp.Token,
			Url:           resp.Url,
			Nota:          resp.Nota,
		}
		result = append(result, results)
	}
	return result, err
}

func (tq *TransactionQuery) GetTransaction(transactionID uint) (*transaction.TransactionList, error) {
	var tm TransactionModel
	if err := tq.db.First(&tm, transactionID).Error; err != nil {
		return nil, err
	}

	// transaksi tidak ditemukan
	if tm.ID == 0 {
		err := errors.New("transaction doesnt exist")
		return nil, err
	}

	result := &transaction.TransactionList{
		TransactionID: int(tm.ID),
		ProductID:     int(tm.ProductID),
		TotalPrice:    int(tm.TotalPrice),
		Status:        tm.Status,
		Timestamp:     tm.CreatedAt,
		Token:         tm.Token,
		Url:           tm.Url,
		Nota:          tm.Nota,
	}

	return result, nil
}

func (tq *TransactionQuery) MidtransCallback(transactionID string) (*transaction.TransactionList, error) {
	var tm TransactionModel
	if err := tq.db.Table("transaction_models").Where("Nota = ?", transactionID).Find(&tm).Error; err != nil {

		fmt.Println(tm)
		return nil, err
	}

	// transaksi tidak ditemukan
	if tm.ID == 0 {
		err := errors.New("transaction doesnt exist")
		return nil, err
	}

	ms := midtrans.MidtransStatus(transactionID)
	tm.Status = ms

	if err := tq.db.Save(&tm).Error; err != nil {
		return nil, err
	}

	result := &transaction.TransactionList{
		TransactionID: int(tm.ID),
		ProductID:     int(tm.ProductID),
		TotalPrice:    int(tm.TotalPrice),
		Status:        ms,
		Timestamp:     tm.CreatedAt,
		Token:         tm.Token,
		Url:           tm.Url,
		Nota:          tm.Nota,
	}

	return result, nil
}
