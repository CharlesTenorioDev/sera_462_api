package gemini

import (
	"context"
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"github.com/sera_backend/internal/config"
	"google.golang.org/api/option"
)

type Gemini interface {
	MontarQuestionario(ctx context.Context, pergunta string) ([]string, error)
}

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(ctx context.Context, cfg *config.Config) (Gemini, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GiminiConfig.API_KEY))
	if err != nil {
		return nil, err
	}
	return &GeminiClient{client: client}, nil
}

func (g *GeminiClient) MontarQuestionario(ctx context.Context, pergunta string) ([]string, error) {
	model := g.client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(pergunta))
	if err != nil {
		return nil, err
	}
	resul, err := printResponse(resp)
	if err != nil {
		return nil, err
	}
	return []string{resul}, nil

}

func printResponse(resp *genai.GenerateContentResponse) (string, error) {
	// Marshal the response object into a JSON byte slice
	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	// Return the JSON string representation
	return string(data), nil
}
