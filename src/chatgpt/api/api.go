package api

import (
	"context"
	"fmt"

	gpt3 "github.com/PullRequestInc/go-gpt3"
)

type ChatGPTApi struct {
	apiKey string
	client gpt3.Client
}

func NewEngine(apiKey string) (*ChatGPTApi, error) {
	return &ChatGPTApi{
		apiKey: apiKey,
		client: gpt3.NewClient(apiKey),
	}, nil
}

func (c *ChatGPTApi) AskQuestion(chatId int64, question string) ([]byte, error) {
	ctx := context.Background()
	resp, err := c.client.CompletionWithEngine(ctx, "text-davinci-002", gpt3.CompletionRequest{
		Prompt:    []string{question},
		MaxTokens: gpt3.IntPtr(4000),
		//Temperature:      gpt3.Float32Ptr(0),
		//TopP:             gpt3.Float32Ptr(1),
		//FrequencyPenalty: 0,
		//PresencePenalty:  0,
		//Stop:             []string{""},
		//Echo:             true,
	})
	if err != nil {
		return nil, fmt.Errorf("[] Failed to get response from openai API: %v", err)
	}
	return []byte(resp.Choices[0].Text), nil
}
