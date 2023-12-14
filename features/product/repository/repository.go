package respository

import (
	"BE-hi-SPEC/features/product"

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
			ID:        s.ID,
			Category:  s.Category,
			Name:      s.Name,
			CPU:       s.CPU,
			RAM:       s.RAM,
			Display:   s.Display,
			Storage:   s.Storage,
			Thickness: s.Thickness,
			Weight:    s.Weight,
			Bluetooth: s.Bluetooth,
			HDMI:      s.HDMI,
			Price:     s.Price,
			Picture:   s.Picture,
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

// SearchProductByName implements product.Repository.
func (pq *ProductQuery) SearchProductByName(name string) ([]product.Product, error) {
	var products []product.Product

	if err := pq.db.Table("product_models").Where("name LIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
