package gpt

import (
	"encoding/json"

	"net/http"

	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/internal/dto"
	"github.com/sera_backend/pkg/service/gpt"
)

func createQuestion(service gpt.GptClientInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		perg := &dto.PerguntIADTO{}

		err := json.NewDecoder(r.Body).Decode(perg)
		if err != nil {
			logger.Error("erro ao converte para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}

		result, err := service.DoRequest("POST", perg.Perguntas)
		if err != nil {
			logger.Error("erro ao criar usuario", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to Insert the User", "codigo": 500}`))
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converte esultada para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}
	})
}
