package healthcheck

import (
	"github.com/go-chi/chi/v5"
	"github.com/sera_backend/pkg/service/healthcheck"
)

func RegisterHealthcheckAPIHandlers(r chi.Router, service healthcheck.HealthcheckServiceInterface) {
	r.Route("/int/v1", func(r chi.Router) {
		r.Get("/healthcheck", healthchecka(service))
	})
}
