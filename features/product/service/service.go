package service

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/helper/jwt"
	"errors"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type ProductServices struct {
	repo product.Repository
}

func New(r product.Repository) product.Service {
	return &ProductServices{
		repo: r,
	}
}

// UpdateProduct implements product.Service.
func (ps *ProductServices) UpdateProduct(productID uint, input product.Product) (product.Product, error) {
	result, err := ps.repo.UpdateProduct(productID, input)
	if err != nil {
		return product.Product{}, errors.New("failed to update the product")
	}

	return result, nil
}

// CariProduct implements product.Service.
func (ps *ProductServices) CariProduct(name string, category string, minPrice uint, maxPrice uint, page, limit int) ([]product.Product, int, error) {
	products, totalPage, err := ps.repo.SearchProduct(name, category, minPrice, maxPrice, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return products, totalPage, err
}

// SatuProduct implements product.Service.
func (ps *ProductServices) SatuProduct(productID uint) (product.Product, error) {
	result, err := ps.repo.GetProductID(productID)
	if err != nil {
		return product.Product{}, errors.New("failed get product")
	}
	return *result, nil
}

func (ps *ProductServices) TalkToGpt(token *golangjwt.Token, newProduct product.Product) (product.Product, error) {
	userId, err := jwt.ExtractToken(token)
	if err != nil {
		return product.Product{}, err
	}

	if err != nil {
		return product.Product{}, err
	}

	result, err := ps.repo.InsertProduct(userId, newProduct)

	return result, nil
}

// SemuaProduct implements product.Service.
func (ps *ProductServices) SemuaProduct(page int, limit int) ([]product.Product, int, error) {
	result, totalPage, err := ps.repo.GetAllProduct(page, limit)
	if err != nil {
		return nil, 0, errors.New("failed get all product")
	}
	return result, totalPage, err
}

func (ps *ProductServices) DelProduct(productID uint) error {
	err := ps.repo.DelProduct(productID)
	return err
}
