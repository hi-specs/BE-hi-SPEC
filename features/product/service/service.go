package service

import (
	"BE-hi-SPEC/config"
	"BE-hi-SPEC/features/product"
	"context"
	"encoding/json"
	"fmt"

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

func (ps *ProductServices) TalkToGpt(newProduct product.Product) (product.Product, error) {
	client := openai.NewClient(config.InitConfig().OPEN_AI_KEY)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "specification of" + newProduct.Name + "written as json following this specification {Name, CPU, RAM, Display, Storage, Thickness, Weight, Bluetooth(yes/no), HDMI(yes/no), Price(in Indonesia with tax format rupiah)} without any explanation",
				},
			},
		},
	)
	// if err != nil {
	// 	fmt.Printf("ChatCompletion error: %v\n", err)
	// 	return
	// }

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &newProduct)
	fmt.Println(resp.Choices[0].Message.Content)
	if err != nil {
		return product.Product{}, err
	}
	result, err := ps.repo.InsertProduct(newProduct)

	return result, nil
}
