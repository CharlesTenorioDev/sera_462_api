package questionarioia

import (
	"encoding/json"
	"strings"

	"net/http"

	"github.com/sera_backend/internal/config/logger"

	"github.com/sera_backend/pkg/service/questionarioia"

	"github.com/sera_backend/pkg/model"
)

func createQuestionario(service questionarioia.QuestionarioServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		quest := &model.Questionario{}

		err := json.NewDecoder(r.Body).Decode(&quest)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		titulo := strings.ToUpper(quest.Titulo)

		// GetByQuetionario(ctx context.Context, IDTurma, IDMateria, IDProfessor, Titulo string)

		peruntaExist := service.GetByQuetionario(r.Context(), quest.IDTurma.Hex(), quest.IDMateria.Hex(), quest.IDProfessor.Hex(), titulo)

		if peruntaExist {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"MSG": "Caro Professor, essa pergunta ja foi cadastrada para essa turma", "codigo": 400}`))
			return
		}

		_, err = service.Create(r.Context(), *quest)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg", err)
			http.Error(w, "Error ou salvar Instituicao"+err.Error(), http.StatusInternalServerError)
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

// func updateInstituicao(service instituicao.InstituicaoServiceInterface) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		idp := chi.URLParam(r, "id")
// 		logger.Info("PEGANDO O PARAMENTRO")

// 		_, err := service.GetByID(r.Context(), idp)
// 		if err != nil {
// 			http.Error(w, "Instituicao nao encontrada", http.StatusNotFound)
// 			return
// 		}

// 		mpg := &model.Instituicao{}
// 		nome := chi.URLParam(r, "nome")
// 		logger.Info("PEGANDO O NOME")
// 		logger.Info(nome)
// 		if nome == "" {
// 			http.Error(w, "o Nome do curso e obrigatório", http.StatusBadRequest)
// 			return
// 		}

// 		mpg.Nome = nome
// 		id, err := primitive.ObjectIDFromHex(idp)
// 		if err != nil {
// 			http.Error(w, "erro ao converter id", http.StatusBadRequest)

// 			return
// 		}

// 		mpg.ID = id
// 		_, err = service.Update(r.Context(), idp, *&mpg)
// 		if err != nil {
// 			logger.Error("erro ao acessar a camada de service do mpg no upd", err)
// 			http.Error(w, "Error ao atualizar meio de pagamento", http.StatusInternalServerError)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(map[string]interface{}{"MSG": "Success", "codigo": 1})
// 	}
// }

// func getByIdInstituicao(service instituicao.InstituicaoServiceInterface) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		idp := chi.URLParam(r, "id")
// 		logger.Info("PEGANDO O PARAMENTRO NA CONSULTA")
// 		result, err := service.GetByID(r.Context(), idp)
// 		if err != nil {
// 			logger.Error("erro ao acessar a camada de service da Instituicao no por id", err)
// 			w.WriteHeader(http.StatusNotFound)
// 			w.Write([]byte(`{"MSG": "Instituicao não encontrada", "codigo": 404}`))
// 			return
// 		}

// 		err = json.NewEncoder(w).Encode(result)
// 		if err != nil {
// 			logger.Error("erro ao converter em json", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(`{"MSG": "Error to parse Bot to JSON", "codigo": 500}`))
// 			return
// 		}
// 	}
// }

// func getAllInstituicao(service instituicao.InstituicaoServiceInterface) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		filters := model.FilterInstituicao{
// 			Nome:    chi.URLParam(r, "nome"),
// 			Enabled: chi.URLParam(r, "enable"),
// 		}

// 		limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
// 		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)

// 		result, err := service.GetAll(r.Context(), filters, limit, page)
// 		if err != nil {
// 			logger.Error("erro ao acessar a camada de service do mpg no upd", err)
// 			w.WriteHeader(http.StatusNotFound)
// 			w.Write([]byte(`{"MSG": "User not found", "codigo": 404}`))
// 			return
// 		}

// 		// Configurando o cabeçalho para resposta JSON usando o middleware
// 		w.Header().Set("Content-Type", "application/json; charset=utf-8")

// 		// Escrevendo a resposta JSON
// 		err = json.NewEncoder(w).Encode(result)
// 		if err != nil {
// 			logger.Error("erro ao converter para json", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
// 			return
// 		}
// 	})
// }
