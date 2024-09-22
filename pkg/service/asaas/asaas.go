package asaas

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/sera_backend/internal/config"
	"github.com/sera_backend/internal/dto"
)

type Client struct {
	apiKey  string
	baseUrl string
	wallet  string
	cliente *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiKey:  cfg.ASAAS_API_KEY,
		baseUrl: cfg.URL_ASAAS,
		wallet:  cfg.ASAAS_WALLET_ID,
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

func (c *Client) DoRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	url := c.baseUrl + endpoint

	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("access_token", c.apiKey)

	return c.cliente.Do(req)
}

func (c *Client) GetClienteByID(clienteID string) (bool, error) {
	endpoint := "/api/v3/customers/" + clienteID
	resp, err := c.DoRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, errors.New("unexpected status code")
}

// Function to create a new client
func (c *Client) CreateCliente(clienteData dto.CustomerInputAsaasDTO) (*http.Response, error) {
	endpoint := "/api/v3/customers"
	return c.DoRequest(http.MethodPost, endpoint, clienteData)
}

// Function to update an existing client
func (c *Client) UpdateCliente(clienteID string, clienteData dto.CustomerInputAsaasDTO) (*http.Response, error) {
	endpoint := "/api/v3/customers/" + clienteID
	return c.DoRequest(http.MethodPut, endpoint, clienteData)
}

// Function to create a new subscription
func (c *Client) NovaAssinatura(subscriptionData dto.SubscriptionInputDTO) (*http.Response, error) {
	endpoint := "/api/v3/subscriptions"
	return c.DoRequest(http.MethodPost, endpoint, subscriptionData)
}

// client.ListaAssinaturas("CREDIT_CARD", "ACTIVE")
func (c *Client) ListaAssinaturas(billingType, status string) (*http.Response, error) {
	endpoint := "/api/v3/subscriptions?billingType=" + billingType + "&status=" + status
	return c.DoRequest(http.MethodGet, endpoint, nil)
}

// F client.GetAssinatura("sub_5u51vszvwg3ab71g")unction to get a specific subscription by ID
func (c *Client) GetAssinatura(subscriptionID string) (*http.Response, error) {
	endpoint := "/api/v3/subscriptions/" + subscriptionID
	return c.DoRequest(http.MethodGet, endpoint, nil)
}
