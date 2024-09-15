package instituicao

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sera_backend/pkg/service/instituicao"
)

func RegisterInstituicaoHandlers(r chi.Router, service instituicao.InstituicaoServiceInterface) {
	r.Route("/api/v1/instituicao", func(r chi.Router) {
		r.Post("/add", createInstituicao(service))
		r.Put("/update/{id}/{nome}", updateInstituicao(service))
		r.Get("/getbyid/{id}", getByIdInstituicao(service))
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllInstituicao(service)
			handler.ServeHTTP(w, r)
		})
	})
}
