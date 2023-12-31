// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	jwt "github.com/golang-jwt/jwt/v5"
	mock "github.com/stretchr/testify/mock"

	product "BE-hi-SPEC/features/product"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CariProduct provides a mock function with given fields: name, category, minPrice, maxPrice, page, limit
func (_m *Service) CariProduct(name string, category string, minPrice string, maxPrice string, page int, limit int) ([]product.Product, int, error) {
	ret := _m.Called(name, category, minPrice, maxPrice, page, limit)

	var r0 []product.Product
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, int, int) ([]product.Product, int, error)); ok {
		return rf(name, category, minPrice, maxPrice, page, limit)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, string, int, int) []product.Product); ok {
		r0 = rf(name, category, minPrice, maxPrice, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]product.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, string, int, int) int); ok {
		r1 = rf(name, category, minPrice, maxPrice, page, limit)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(string, string, string, string, int, int) error); ok {
		r2 = rf(name, category, minPrice, maxPrice, page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DelProduct provides a mock function with given fields: token, productID
func (_m *Service) DelProduct(token *jwt.Token, productID uint) error {
	ret := _m.Called(token, productID)

	var r0 error
	if rf, ok := ret.Get(0).(func(*jwt.Token, uint) error); ok {
		r0 = rf(token, productID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SatuProduct provides a mock function with given fields: productID
func (_m *Service) SatuProduct(productID uint) (product.Product, error) {
	ret := _m.Called(productID)

	var r0 product.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (product.Product, error)); ok {
		return rf(productID)
	}
	if rf, ok := ret.Get(0).(func(uint) product.Product); ok {
		r0 = rf(productID)
	} else {
		r0 = ret.Get(0).(product.Product)
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(productID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SemuaProduct provides a mock function with given fields: page, limit
func (_m *Service) SemuaProduct(page int, limit int) ([]product.Product, int, error) {
	ret := _m.Called(page, limit)

	var r0 []product.Product
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(int, int) ([]product.Product, int, error)); ok {
		return rf(page, limit)
	}
	if rf, ok := ret.Get(0).(func(int, int) []product.Product); ok {
		r0 = rf(page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]product.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int) int); ok {
		r1 = rf(page, limit)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(int, int) error); ok {
		r2 = rf(page, limit)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// TalkToGpt provides a mock function with given fields: token, newProduct
func (_m *Service) TalkToGpt(token *jwt.Token, newProduct product.Product) (product.Product, error) {
	ret := _m.Called(token, newProduct)

	var r0 product.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(*jwt.Token, product.Product) (product.Product, error)); ok {
		return rf(token, newProduct)
	}
	if rf, ok := ret.Get(0).(func(*jwt.Token, product.Product) product.Product); ok {
		r0 = rf(token, newProduct)
	} else {
		r0 = ret.Get(0).(product.Product)
	}

	if rf, ok := ret.Get(1).(func(*jwt.Token, product.Product) error); ok {
		r1 = rf(token, newProduct)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProduct provides a mock function with given fields: token, productID, input
func (_m *Service) UpdateProduct(token *jwt.Token, productID uint, input product.Product) (product.Product, error) {
	ret := _m.Called(token, productID, input)

	var r0 product.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(*jwt.Token, uint, product.Product) (product.Product, error)); ok {
		return rf(token, productID, input)
	}
	if rf, ok := ret.Get(0).(func(*jwt.Token, uint, product.Product) product.Product); ok {
		r0 = rf(token, productID, input)
	} else {
		r0 = ret.Get(0).(product.Product)
	}

	if rf, ok := ret.Get(1).(func(*jwt.Token, uint, product.Product) error); ok {
		r1 = rf(token, productID, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
