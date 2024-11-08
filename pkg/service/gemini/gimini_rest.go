package gemini

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sera_backend/internal/config"
	"github.com/sera_backend/internal/config/logger"
)

// Add these structs
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type Content struct {
	Parts []string `json:"parts"`
}

type GminiClientInterface interface {
	DoRequest(method string, contents []map[string]interface{}) (string, error)
}

type ClientGemini struct {
	apiKey  string
	baseUrl string
	cliente *http.Client
}

// Ensure Client implements GminiClientInterface
var _ GminiClientInterface = (*ClientGemini)(nil)

func NewClient(cfg *config.Config) *ClientGemini {
	return &ClientGemini{
		apiKey:  cfg.GiminiConfig.API_KEY,
		baseUrl: cfg.GiminiConfig.URL,

		cliente: &http.Client{
			Timeout: time.Duration(cfg.ASAAS_TIMEOUT) * time.Second,
			Transport: &http.Transport{
				ForceAttemptHTTP2:   false,
				MaxConnsPerHost:     1,
				MaxIdleConns:        1,
				MaxIdleConnsPerHost: 1,
				TLSHandshakeTimeout: time.Duration(10) * time.Second,
			},
		},
	}
}

func cleanString(input string) string {
	cleaned := strings.ReplaceAll(input, "\\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\\t", "\t")
	cleaned = strings.ReplaceAll(cleaned, "\\\"", "\"")
	cleaned = strings.ReplaceAll(cleaned, "**\n", "")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.ReplaceAll(cleaned, "**", "")

	return cleaned
}

func (r *GeminiResponse) GetPartsContent() string {
	if len(r.Candidates) == 0 || len(r.Candidates[0].Content.Parts) == 0 {
		return ""
	}

	// Get the content from the first part
	content := r.Candidates[0].Content.Parts[0]

	// Clean the content
	cleaned := cleanString(content)

	return cleaned
}

// Modify DoRequest method to include response processing
func (c *ClientGemini) DoRequest(method string, contents []map[string]interface{}) (string, error) {
	url := c.baseUrl + c.apiKey

	payload := map[string]interface{}{
		"contents": contents,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	resp, err := c.cliente.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse a resposta JSON
	var response GeminiResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return "", err
	}

	// Extrair o conteúdo presente em Parts
	partsContent := ""
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		partsContent = response.Candidates[0].Content.Parts[0]
	}

	// Limpeza adicional se necessário
	partsContent = cleanString(partsContent)

	logger.Info("Parts Content: " + partsContent)
	return partsContent, nil
}
