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
	DoRequest(method, endpoint string, messages []map[string]string) (*http.Response, error)
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

func (c *ClientGpt) DoRequest(method, endpoint string, messages []map[string]string) (*http.Response, error) {
	url := c.baseUrl + endpoint

	// Create the payload with the model from the config
	payload := map[string]interface{}{
		"model":    c.model,
		"messages": messages,
	}

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("access_token", c.apiKey)

	logger.Info("api key" + c.apiKey)

	// Logging the request for debugging
	logger.Info("Sending request to URL: " + url)

	resp, err := c.cliente.Do(req)
	if err != nil {
		return nil, err
	}

	logger.Info("Response Status Code: " + strconv.Itoa(resp.StatusCode))
	return resp, nil
}
