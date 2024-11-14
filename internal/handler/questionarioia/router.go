package questionarioia

import (
	"github.com/go-chi/chi/v5"
	"github.com/sera_backend/pkg/service/questionarioia"
)

func RegisterQuestionarioHandlers(r chi.Router, service questionarioia.QuestionarioServiceInterface) {
	r.Route("/api/v1/questionario", func(r chi.Router) {
		r.Post("/add", createQuestionario(service))
		// r.Put("/update/{id}/{nome}", updateInstituicao(service))
		// r.Get("/getbyid/{id}", getByIdInstituicao(service))
		// r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
		// 	handler := getAllInstituicao(service)
		// 	handler.ServeHTTP(w, r)
		// })
	})
}
