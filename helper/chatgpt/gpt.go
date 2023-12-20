package gpt

import (
	"BE-hi-SPEC/config"
	"BE-hi-SPEC/features/product"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func GptAPI(name string) (newProduct product.Product) {
	client := openai.NewClient(config.InitConfig().OPEN_AI_KEY)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "specification of" + name + "written as json following this specification {Name, CPU, RAM, Display, Storage, Thickness, Weight, Bluetooth(yes/no), HDMI(yes/no), Price(in Indonesia integer))} without any explanation",
				},
			},
		},
	)
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &newProduct)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return product.Product{}
	}
	return newProduct
}
