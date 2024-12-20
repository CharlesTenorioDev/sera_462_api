package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetJwtInput struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
	Role  string `json:"role"`
}

type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}

type CustomerInputAsaasDTO struct {
	Name              string `json:"name"`
	CpfCnpj           string `json:"cpfCnpj"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	MobilePhone       string `json:"mobilePhone"`
	Address           string `json:"address"`
	AddressNumber     string `json:"addressNumber"`
	Province          string `json:"province"`
	PostalCode        string `json:"postalCode"`
	ExternalReference string `json:"externalReference"`
}

type SubscriptionInputDTO struct {
	BillingType       string  `json:"billingType"`
	Cycle             string  `json:"cycle"`
	Customer          string  `json:"customer"`
	Value             float64 `json:"value"`
	NextDueDate       string  `json:"nextDueDate"`
	Description       string  `json:"description"`
	MaxPayments       int     `json:"maxPayments"`
	ExternalReference string  `json:"externalReference"`
}
type PerguntIADTO struct {
	Perguntas string `json:"perguntas"`
}

type QuestionarioParaFilaDTO struct {
	ID       primitive.ObjectID `json:"id"`
	Titulo   string             `json:"titulo"`
	Pergunta string             `json:"pergunta"`
}

type GeneratePayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type LlamaResponse struct {
	CreatedAt string         `json:"created_at"`
	Response  map[int]string `json:"response"`
}

type GeneratePayloadGroq struct {
	Model    string        `json:"model"`
	Messages []MessageGroq `json:"messages"`
}

type MessageGroq struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
