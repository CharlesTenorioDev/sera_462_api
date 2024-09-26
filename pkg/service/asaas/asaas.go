package asaas

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sera_backend/internal/config"
	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/internal/dto"
)

// Define an interface for the Asaas client
type AsaasClientInterface interface {
	DoRequest(method, endpoint string, payload io.Reader) (*http.Response, error)
	GetClienteByID(clienteID string) (bool, error)
	CreateCliente(ctx context.Context, conclienteData dto.CustomerInputAsaasDTO) (bool, error)
	UpdateCliente(clienteID string, clienteData dto.CustomerInputAsaasDTO) (bool, error)
	NovaAssinatura(subscriptionData dto.SubscriptionInputDTO) (bool, error)
	ListaAssinaturas(billingType, status string) (*http.Response, error)
	GetAssinatura(subscriptionID string) (*http.Response, error)
}

type Client struct {
	apiKey  string
	baseUrl string
	wallet  string
	cliente *http.Client
}

// Ensure Client implements AsaasClientInterface
var _ AsaasClientInterface = (*Client)(nil)

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

func (c *Client) DoRequest(method, endpoint string, payload io.Reader) (*http.Response, error) {
	url := c.baseUrl + endpoint

	req, err := http.NewRequest(method, url, payload)
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

func (c *Client) CreateCliente(ctx context.Context, clienteData dto.CustomerInputAsaasDTO) (bool, error) {
	endpoint := "/api/v3/customers"

	// Construct the payload as a JSON string
	payload := strings.NewReader(`{
		"name": "` + clienteData.Name + `",
		"cpfCnpj": "` + clienteData.CpfCnpj + `",
		"email": "` + clienteData.Email + `",
		"phone": "` + clienteData.Phone + `",
		"mobilePhone": "` + clienteData.MobilePhone + `",
		"address": "` + clienteData.Address + `",
		"addressNumber": "` + clienteData.AddressNumber + `",
		"province": "` + clienteData.Province + `",
		"postalCode": "` + clienteData.PostalCode + `",
		"externalReference": "` + clienteData.ExternalReference + `"
	}`)

	resp, err := c.DoRequest(http.MethodPost, endpoint, payload)
	if err != nil {
		logger.Error("failed to create cliente no asas_api_key", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, errors.New("failed to create cliente")
}

func (c *Client) UpdateCliente(clienteID string, clienteData dto.CustomerInputAsaasDTO) (bool, error) {
	endpoint := "/api/v3/customers/" + clienteID

	payload := strings.NewReader(`{
		"name": "` + clienteData.Name + `",
		"cpfCnpj": "` + clienteData.CpfCnpj + `",
		"email": "` + clienteData.Email + `",
		"phone": "` + clienteData.Phone + `",
		"mobilePhone": "` + clienteData.MobilePhone + `",
		"address": "` + clienteData.Address + `",
		"addressNumber": "` + clienteData.AddressNumber + `",
		"province": "` + clienteData.Province + `",
		"postalCode": "` + clienteData.PostalCode + `",
		"externalReference": "` + clienteData.ExternalReference + `"
	}`)
	resp, err := c.DoRequest(http.MethodPut, endpoint, payload)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, errors.New("failed to update cliente")
}

func (c *Client) NovaAssinatura(subscriptionData dto.SubscriptionInputDTO) (bool, error) {
	endpoint := "/api/v3/subscriptions"

	// Construct the JSON payload manually
	payload := strings.NewReader(`{
		"billingType": "` + subscriptionData.BillingType + `",
		"cycle": "` + subscriptionData.Cycle + `",
		"customer": "` + subscriptionData.Customer + `",
		"value": ` + fmt.Sprintf("%f", subscriptionData.Value) + `,
		"nextDueDate": "` + subscriptionData.NextDueDate + `",
		"description": "` + subscriptionData.Description + `",
		"maxPayments": ` + fmt.Sprintf("%d", subscriptionData.MaxPayments) + `,
		"externalReference": "` + subscriptionData.ExternalReference + `"
	}`)

	resp, err := c.DoRequest(http.MethodPost, endpoint, payload)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, errors.New("failed to create subscription")
}

func (c *Client) ListaAssinaturas(billingType, status string) (*http.Response, error) {
	endpoint := "/api/v3/subscriptions?billingType=" + billingType + "&status=" + status
	return c.DoRequest(http.MethodGet, endpoint, nil)
}

func (c *Client) GetAssinatura(subscriptionID string) (*http.Response, error) {
	endpoint := "/api/v3/subscriptions/" + subscriptionID
	return c.DoRequest(http.MethodGet, endpoint, nil)
}
