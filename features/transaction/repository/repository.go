package repository

import (
	"BE-hi-SPEC/features/product"
	p "BE-hi-SPEC/features/product"
	pr "BE-hi-SPEC/features/product/repository"
	"BE-hi-SPEC/features/transaction"
	"BE-hi-SPEC/features/user"
	"BE-hi-SPEC/features/user/repository"
	ur "BE-hi-SPEC/features/user/repository"
	"BE-hi-SPEC/helper/gofpdf"
	"BE-hi-SPEC/helper/midtrans"
	"errors"
	"fmt"
	"log"
	"strconv"

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

func (tq *TransactionQuery) AdminDashboard(userID uint, page int, limit int) (transaction.TransactionDashboard, int, error) {
	// mendapatkan nilai total product
	var productCount int
	tableName := "product_models"
	columnName := "deleted_at"
	query := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NULL", tableName, columnName)
	err := tq.db.Raw(query).Scan(&productCount).Error
	if err != nil {
		log.Fatal(err)
	}

	// mendapatkan nilai total user
	var userCount int
	tableNameUser := "user_models"
	columnNameUser := "deleted_at"
	queryuser := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NULL", tableNameUser, columnNameUser)
	err2 := tq.db.Raw(queryuser).Scan(&userCount).Error
	if err2 != nil {
		log.Fatal(err)
	}

	// mendapatkan nilai total transaksi yang sukses
	var transactionCount int
	tableNameTransaction := "transaction_models"
	columnNameTransaction := "created_at"
	querytransaction := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NOT NULL", tableNameTransaction, columnNameTransaction)
	err3 := tq.db.Raw(querytransaction).Scan(&transactionCount).Error
	if err3 != nil {
		log.Fatal(err)
	}

	var products []pr.ProductModel
	offset := (page - 1) * limit
	if err := tq.db.Table("product_models").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return transaction.TransactionDashboard{}, 0, err
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

	var totalPage int
	tableName2 := "transaction_models"
	columnName2 := "deleted_at"
	queryuser2 := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NULL", tableName2, columnName2)
	err5 := tq.db.Raw(queryuser2).Scan(&totalPage).Error
	if err != nil {
		log.Fatal(err5)
	}

	if totalPage%limit == 0 {
		totalPage = totalPage / limit
	} else {
		totalPage = totalPage / limit
		totalPage++
	}

	return *result, totalPage, err
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

func (tq *TransactionQuery) TransactionList(page, limit int) ([]transaction.TransactionList, int, error) {
	var tm []TransactionModel
	var result []transaction.TransactionList
	offset := (page - 1) * limit
	if err := tq.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tm).Error; err != nil {
		return nil, 0, err
	}

	// get list of user ID
	var userID []int
	for _, result := range tm {
		userID = append(userID, int(result.UserID))
	}

	// get list of user
	var User []ur.UserModel
	for _, result := range userID {
		tmp := new(repository.UserModel)
		tq.db.Table("user_models").Where("id = ?", result).Find(&tmp)
		User = append(User, *tmp)
	}

	// get list of product ID
	var productID []int
	for _, result := range tm {
		productID = append(productID, int(result.ProductID))
	}

	// get list of product
	var Product []pr.ProductModel
	for _, result := range productID {
		tmp := new(pr.ProductModel)
		tq.db.Table("product_models").Where("id = ?", result).Find(&tmp)
		Product = append(Product, *tmp)
	}

	// slicing data that we have, into the struct for return

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
			Users:         User,
			Products:      Product,
		}
		result = append(result, results)
	}

	// fmt.Println(len(tm), len(User))
	// mendapatkan nilai total transaction
	var totalPage int
	tableNameUser := "transaction_models"
	columnNameUser := "deleted_at"
	queryuser := fmt.Sprintf("SELECT COUNT(*) AS null_count FROM %s WHERE %s IS NULL", tableNameUser, columnNameUser)
	err := tq.db.Raw(queryuser).Scan(&totalPage).Error
	if err != nil {
		log.Fatal(err)
	}

	if totalPage%limit == 0 {
		totalPage = totalPage / limit
	} else {
		totalPage = totalPage / limit
		totalPage++
	}

	return result, totalPage, err
}

func (tq *TransactionQuery) GetTransaction(userID uint, transactionID uint) (*transaction.TransactionList, error) {
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

func (tq *TransactionQuery) UserTransaction(userId int, userID uint) (transaction.UserTransaction, error) {
	var tl []transaction.Transaction
	var pl []product.Product
	var user user.User

	// mendapatkan detail user
	tq.db.Table("user_models").Where("id = ?", userID).Find(&user)

	// mendapatkan list product
	var productID []int
	tq.db.Table("transaction_models").Where("user_id = ?", userID).Select("product_id").Find(&productID)
	for _, prod := range productID {
		tmp := new(product.Product)
		tq.db.Table("product_models").Where("id = ?", prod).Find(&tmp)
		pl = append(pl, *tmp)
	}

	// mendapatkan list transaksi
	var transID []int
	tq.db.Table("transaction_models").Where("user_id = ?", userID).Select("id").Find(&transID)
	for _, trans := range transID {
		tmp := new(transaction.Transaction)
		tq.db.Table("transaction_models").Where("id = ?", trans).Find(&tmp)
		tl = append(tl, *tmp)
	}
	var result transaction.UserTransaction
	result.User = user
	result.Product = pl
	result.Transaction = tl

	return result, nil
}

func (tq *TransactionQuery) DownloadTransaction(userID uint, transactionID uint) error {
	var trans TransactionModel
	var user ur.UserModel
	var prod pr.ProductModel
	// get transaction detail
	if err := tq.db.Table("transaction_models").Where("id = ?", transactionID).Find(&trans).Error; err != nil {
		return err
	}

	// get product detail
	if err := tq.db.Table("product_models").Where("id = ?", trans.ProductID).Find(&prod).Error; err != nil {
		return err
	}

	// get user detail
	if err := tq.db.Table("user_models").Where("id = ?", trans.UserID).Find(&user).Error; err != nil {
		return err
	}

	if userID != trans.UserID {
		return errors.New("no authorized")
	}
	fmt.Println(trans, user, prod)
	// Convert TransactionModel to gofpdf.TM
	// tm := gofpdf.TM(trans)
	// um := gofpdf.UM(user)
	// pm := gofpdf.PM(prod)
	gofpdf.GeneratePDF(gofpdf.TM(trans), gofpdf.PM(prod), gofpdf.UM(user))
	fmt.Println(trans)

	return nil
}
