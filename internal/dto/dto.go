package dto

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
