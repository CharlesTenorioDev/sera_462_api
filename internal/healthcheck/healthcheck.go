package healthcheck

import (
	"net/http"

	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/pkg/service/healthcheck"
)

func healthchecka(service healthcheck.HealthcheckServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok, err := service.CheckDB()

		if err != nil && !ok {
			logger.Error("erro ao verificar o banco de dados", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Erro ao verificar o banco de dados", "codigo": 500}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"MSG": "Banco de dados MongoDB ok", "codigo": 200}`))
	}
}
