package gpt

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sera_backend/internal/config"
	"github.com/sera_backend/internal/config/logger"
)

type GptClientInterface interface {
	DoRequest(method, messages string) (map[string]interface{}, error)
}

type ClientGpt struct {
	apiKey  string
	baseUrl string
	model   string
	cliente *http.Client
}

// Ensure Client implements GptClientInterface
var _ GptClientInterface = (*ClientGpt)(nil)

func NewClient(cfg *config.Config) *ClientGpt {
	return &ClientGpt{
		apiKey:  cfg.GptConfig.SRV_GPT_API_KEY,
		baseUrl: cfg.GptConfig.SRV_GPT_URL,
		model:   cfg.GptConfig.SRV_GPT_MODEL,

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

func (c *ClientGpt) DoRequest(method, messages string) (map[string]interface{}, error) {
	url := c.baseUrl

	// Create the payload with the model from the config
	payload := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": messages},
		},
	}

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.Error("error to convert Client to JSON para GPT", err)
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		logger.Error("error to create request GPT", err)
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	resp, err := c.cliente.Do(req)
	if err != nil {
		logger.Error("error to send request GPT", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response body
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("error decoding response body", err)
		return nil, err
	}

	logger.Info("Response Status Code: " + strconv.Itoa(resp.StatusCode))
	return result, nil
}
