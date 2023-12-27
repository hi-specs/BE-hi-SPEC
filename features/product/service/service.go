package service

import (
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/helper/jwt"
	"errors"
	"strconv"

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
func (ps *ProductServices) UpdateProduct(token *golangjwt.Token, productID uint, input product.Product) (product.Product, error) {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return product.Product{}, err
	}
	if rolesUser == "" {
		return product.Product{}, err
	}
	if rolesUser != "admin" {
		return product.Product{}, errors.New("unauthorized access: admin role required")
	}
	result, err := ps.repo.UpdateProduct(userId, productID, input)
	if err != nil {
		return product.Product{}, errors.New("failed to update the product")
	}

	return result, nil
}

// CariProduct implements product.Service.
func (ps *ProductServices) CariProduct(name string, category string, minPrice string, maxPrice string, page, limit int) ([]product.Product, int, error) {
	// checking value of minPrice
	if minPrice == "" {
		minPrice = "0"
	}
	minp, err := strconv.Atoi(minPrice)
	if minp == 0 {
		minp = 1
	}
	if err != nil {
		return []product.Product{}, 0, errors.New("MinPrice Value Must Integer")
	}

	// checking value of maxPrice
	if maxPrice == "" {
		maxPrice = "0"
	}
	maxp, err := strconv.Atoi(maxPrice)
	if err != nil {
		return []product.Product{}, 0, errors.New("MaxPrice Value Must Integer")
	}
	products, totalPage, err := ps.repo.SearchProduct(name, category, uint(minp), uint(maxp), page, limit)
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
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return product.Product{}, errors.New("Token Error")
	}
	if rolesUser != "admin" {
		return product.Product{}, errors.New("unauthorized access: admin role required")
	}

	result, err := ps.repo.InsertProduct(userId, newProduct)
	if err != nil {
		return product.Product{}, errors.New("Inputan tidak boleh kosong")
	}

	return result, err
}

// SemuaProduct implements product.Service.
func (ps *ProductServices) SemuaProduct(page int, limit int) ([]product.Product, int, error) {
	result, totalPage, err := ps.repo.GetAllProduct(page, limit)
	if err != nil {
		return nil, 0, errors.New("failed get all product")
	}
	return result, totalPage, err
}

func (ps *ProductServices) DelProduct(token *golangjwt.Token, productID uint) error {
	userId, rolesUser, err := jwt.ExtractToken(token)
	if err != nil {
		return err
	}
	if rolesUser != "admin" {
		return errors.New("unauthorized access: admin role required")
	}
	err = ps.repo.DelProduct(userId, productID)
	return err
}
