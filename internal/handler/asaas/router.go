package asaas

import (
	"github.com/go-chi/chi/v5"
	"github.com/sera_backend/pkg/service/asaas"
)

func RegisterAsaasHandlers(r chi.Router, service asaas.AsaasClientInterface) {
	r.Route("/api/v1/asaas", func(r chi.Router) {
		r.Post("/add", createClienteAsaas(service))

	})
}
