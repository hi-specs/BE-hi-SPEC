package respository

import (
	"BE-hi-SPEC/features/product"

	"gorm.io/gorm"
)

type ProductModel struct {
	gorm.Model
	Name      string
	CPU       string
	RAM       string
	Display   string
	Storage   string
	Thickness string
	Weight    string
	Bluetooth string
	HDMI      string
	Price     string
}

type ProductQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) product.Repository {
	return &ProductQuery{
		db: db,
	}
}

func (gq ProductQuery) InsertProduct(newProduct product.Product) (product.Product, error) {
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

	gq.db.Create(&inputDB)
	return newProduct, nil
}