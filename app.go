package main

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type Application struct {
	Addr        string
	Client      *openai.Client
	AssistantID string
	ThreadID    string
	FileID      string
}

func NewApplication() (*Application, error) {
	client := openai.NewClient(os.Getenv("API_KEY"))

	assistant, err := client.CreateAssistant(context.Background(), openai.AssistantRequest{
		Name:         &assistantName,
		Description:  &description,
		Model:        openai.GPT4TurboPreview,
		Instructions: &instructions,
		Tools:        []openai.AssistantTool{{Type: openai.AssistantToolTypeRetrieval}},
	})
	if err != nil {
		return nil, err
	}

	return &Application{
		Addr:        os.Getenv("ADDR"),
		Client:      client,
		AssistantID: assistant.ID,
	}, nil
}
