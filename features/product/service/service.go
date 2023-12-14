package service

import (
	"BE-hi-SPEC/config"
	"BE-hi-SPEC/features/product"
	"BE-hi-SPEC/helper/jwt"
	"context"
	"encoding/json"
	"errors"

	golangjwt "github.com/golang-jwt/jwt/v5"
	openai "github.com/sashabaranov/go-openai"
)

type ProductServices struct {
	repo product.Repository
}

func New(r product.Repository) product.Service {
	return &ProductServices{
		repo: r,
	}
}

// CariProduct implements product.Service.
func (ps *ProductServices) CariProduct(name string) ([]product.Product, error) {
	products, err := ps.repo.SearchProductByName(name)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// SatuProduct implements product.Service.
func (ps *ProductServices) SatuProduct(productID uint) (product.Product, error) {
	result, err := ps.repo.GetProductID(productID)
	if err != nil {
		return product.Product{}, errors.New("failed get all product")
	}
	return *result, nil
}

func (ps *ProductServices) TalkToGpt(token *golangjwt.Token, newProduct product.Product) (product.Product, error) {
	userId, err := jwt.ExtractToken(token)
	if err != nil {
		return product.Product{}, err
	}
	client := openai.NewClient(config.InitConfig().OPEN_AI_KEY)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "specification of" + newProduct.Name + "written as json following this specification {Name, CPU, RAM, Display, Storage, Thickness, Weight, Bluetooth(yes/no), HDMI(yes/no), Price(in Indonesia, with tax, and format Rp.)} without any explanation",
				},
			},
		},
	)
	// if err != nil {
	// 	fmt.Printf("ChatCompletion error: %v\n", err)
	// 	return
	// }

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &newProduct)
	// fmt.Println(resp.Choices[0].Message.Content)
	if err != nil {
		return product.Product{}, err
	}
	result, err := ps.repo.InsertProduct(userId, newProduct)

	return result, nil
}

// SemuaProduct implements product.Service.
func (ps *ProductServices) SemuaProduct(page int, limit int) ([]product.Product, error) {
	result, err := ps.repo.GetAllProduct(page, limit)
	if err != nil {
		return nil, errors.New("failed get all product")
	}
	return result, nil
}
