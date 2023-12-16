package respository

import (
	"BE-hi-SPEC/features/product"
	"errors"

	"gorm.io/gorm"
)

type ProductModel struct {
	gorm.Model
	Category  string
	Name      string
	CPU       string
	RAM       string
	Display   string
	Storage   string
	Thickness string
	Weight    string
	Bluetooth string
	HDMI      string
	Price     int
	Picture   string
}

type ProductQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) product.Repository {
	return &ProductQuery{
		db: db,
	}
}

// UpdateProduct implements product.Repository.
func (pq *ProductQuery) UpdateProduct(productID uint, input product.Product) (product.Product, error) {
	var proses ProductModel
	if err := pq.db.First(&proses, input.ID).Error; err != nil {
		return product.Product{}, err
	}

	// Jika tidak ada buku ditemukan
	if proses.ID == 0 {
		err := errors.New("user tidak ditemukan")
		return product.Product{}, err
	}
	if input.Category != "" {
		proses.Category = input.Category
	}
	if input.Name != "" {
		proses.Name = input.Name
	}
	if input.CPU != "" {
		proses.CPU = input.CPU
	}
	if input.RAM != "" {
		proses.RAM = input.RAM
	}
	if input.Display != "" {
		proses.Display = input.Display
	}
	if input.Name != "" {
		proses.Name = input.Name
	}
	if input.Storage != "" {
		proses.Storage = input.Storage
	}
	if input.Thickness != "" {
		proses.Thickness = input.Thickness
	}
	if input.Weight != "" {
		proses.Weight = input.Weight
	}

	if input.Bluetooth != "" {
		proses.Bluetooth = input.Bluetooth
	}

	if input.HDMI != "" {
		proses.HDMI = input.HDMI
	}

	if input.Price != 0 {
		proses.Price = input.Price
	}

	if input.Picture != "" {
		proses.Picture = input.Picture
	}
	if err := pq.db.Save(&proses).Error; err != nil {

		return product.Product{}, err
	}
	result := product.Product{
		ID:        proses.ID,
		Category:  proses.Category,
		Name:      proses.Name,
		CPU:       proses.CPU,
		RAM:       proses.RAM,
		Display:   proses.Display,
		Storage:   proses.Storage,
		Thickness: proses.Thickness,
		Weight:    proses.Weight,
		Bluetooth: proses.Bluetooth,
		HDMI:      proses.HDMI,
		Price:     proses.Price,
		Picture:   proses.Picture,
	}
	return result, nil
}

// SearchProduct implements product.Repository.
func (pq *ProductQuery) SearchProduct(name string, category string, minPrice uint, maxPrice uint) ([]product.Product, error) {
	var products []product.Product

	qry := pq.db.Table("product_models")

	if name != "" {
		qry = qry.Where("name like ?", "%"+name+"%")
	}

	if category != "" {
		qry = qry.Where("category like ?", "%"+category+"%")
	}

	if minPrice != 0 {
		qry = qry.Where("price >= ?", minPrice)
	}

	if maxPrice != 0 {
		qry = qry.Where("price <= ?", maxPrice)
	}

	if err := qry.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetAllProduct implements product.Repository.
func (pq *ProductQuery) GetAllProduct(page int, limit int) ([]product.Product, error) {
	var products []ProductModel
	offset := (page - 1) * limit
	if err := pq.db.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}
	var result []product.Product
	for _, s := range products {
		result = append(result, product.Product{
			ID:       s.ID,
			Name:     s.Name,
			Price:    s.Price,
			Category: s.Category,
			Picture:  s.Picture,
		})
	}
	return result, nil
}

func (gq ProductQuery) InsertProduct(UserID uint, newProduct product.Product) (product.Product, error) {
	var inputDB = new(ProductModel)
	inputDB.Name = newProduct.Name
	inputDB.CPU = newProduct.CPU
	inputDB.RAM = newProduct.RAM
	inputDB.Display = newProduct.Display
	inputDB.Storage = newProduct.Storage
	inputDB.Thickness = newProduct.Thickness
	inputDB.Weight = newProduct.Weight
	inputDB.Bluetooth = newProduct.Bluetooth
	inputDB.HDMI = newProduct.HDMI
	inputDB.Price = newProduct.Price
	inputDB.Picture = newProduct.Picture
	inputDB.Category = newProduct.Category

	gq.db.Create(&inputDB)
	newProduct.ID = inputDB.ID
	return newProduct, nil
}

// GetProductID implements product.Repository.
func (pq *ProductQuery) GetProductID(productID uint) (*product.Product, error) {
	var productModel ProductModel
	if err := pq.db.First(&productModel, productID).Error; err != nil {
		return nil, err
	}
	result := &product.Product{
		ID:        productModel.ID,
		Category:  productModel.Category,
		Name:      productModel.Name,
		CPU:       productModel.CPU,
		RAM:       productModel.RAM,
		Display:   productModel.Display,
		Storage:   productModel.Storage,
		Thickness: productModel.Thickness,
		Weight:    productModel.Weight,
		Bluetooth: productModel.Bluetooth,
		HDMI:      productModel.HDMI,
		Price:     productModel.Price,
		Picture:   productModel.Picture,
	}
	return result, nil
}

func (pq *ProductQuery) DelProduct(productID uint) error {
	var prod = new(ProductModel)
	if err := pq.db.Where("id", productID).Find(&prod).Error; err != nil {
		return err
	}

	pq.db.Where("id", productID).Delete(&prod)
	return nil
}
