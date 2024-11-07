package gemini

import (
	"encoding/json"

	"net/http"

	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/pkg/model"
	"github.com/sera_backend/pkg/service/gemini"
)

func createQuestion(service gemini.GeminiServiceInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		perg := &model.PerguntaGeminai{}

		err := json.NewDecoder(r.Body).Decode(perg)

		if err != nil {
			logger.Error("erro ao converte para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}
		logger.Info(perg.Perguntas)
		// if perg.Perguntas != " " {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte(`{"MSG": "Pergunta e obrigatória", "codigo": 400}`))
		// 	return
		// }

		result, err := service.MontarQuestionario(r.Context(), perg.Perguntas)
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
