package asaas

import (
	"encoding/json"

	"net/http"

	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/internal/dto"
	"github.com/sera_backend/pkg/service/asaas"
)

func createClienteAsaas(service asaas.AsaasClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliAsaas := &dto.CustomerInputAsaasDTO{}

		err := json.NewDecoder(r.Body).Decode(&cliAsaas)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		criado, err := service.CreateCliente(r.Context(), *cliAsaas)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg", err)
			http.Error(w, "Error ou salvar Instituicao"+err.Error(), http.StatusInternalServerError)
			return
		}

		if !criado {
			http.Error(w, "Error ou salvar LinkedInuicao"+err.Error(), http.StatusInternalServerError)
			return
		}

		type Response struct {
			Message string `json:"message"`
		}

		// Crie uma instância da estrutura com a mensagem desejada.
		msg := Response{
			Message: "Dados gravados com sucesso",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	}
}
