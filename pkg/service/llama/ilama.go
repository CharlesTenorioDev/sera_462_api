package llama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sera_backend/internal/config"
	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/pkg/model"
)

type LlhamaClientInterface interface {
	DoRequest(method, endpoint string, payload io.Reader) (*http.Response, error)
	CreateQuestion(question string) ([]model.Pergunta, error)
}

type ClientLlama struct {
	apiKey  string
	baseUrl string
	cliente *http.Client
}

type LlamaResponse struct {
	Questions []struct {
		Content string   `json:"content"`
		Options []string `json:"options,omitempty"`
		Answer  string   `json:"answer"`
		Type    string   `json:"type"`
	} `json:"questions"`
}

var _ LlhamaClientInterface = (*ClientLlama)(nil)

func NewClient(cfg *config.Config) *ClientLlama {
	return &ClientLlama{
		apiKey:  cfg.LlamaConfig.SRV_Llama_API_KEY,
		baseUrl: cfg.LlamaConfig.SRV_Llama_URL,
		cliente: &http.Client{
			Timeout: time.Duration(cfg.AsaasConfig.ASAAS_TIMEOUT) * time.Second,
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

type QuestionPayload struct {
	Question string `json:"question"`
}

func (c *ClientLlama) DoRequest(method, endpoint string, payload io.Reader) (*http.Response, error) {
	url := c.baseUrl + endpoint

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("access_token", c.apiKey)

	logger.Info("api key" + c.apiKey)
	logger.Info("Sending request to URL: " + url)

	resp, err := c.cliente.Do(req)
	if err != nil {
		return nil, err
	}

	logger.Info("Response Status Code: " + strconv.Itoa(resp.StatusCode))
	return resp, nil
}

func (c *ClientLlama) CreateQuestion(question string) ([]model.Pergunta, error) {
	payload := QuestionPayload{
		Question: question,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar payload: %v", err)
	}

	resp, err := c.DoRequest("POST", "/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição à API Llama: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta da API: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Llama retornou status %d: %s", resp.StatusCode, string(body))
	}

	var llamaResp LlamaResponse
	err = json.Unmarshal(body, &llamaResp)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar resposta da API: %v", err)
	}

	var perguntas []model.Pergunta
	for _, q := range llamaResp.Questions {
		pergunta := model.Pergunta{
			Tipo:            q.Type,
			Conteudo:        q.Content,
			RespostaCorreta: q.Answer,
			Opcoes:          q.Options,
		}
		perguntas = append(perguntas, pergunta)
	}

	return perguntas, nil
}
