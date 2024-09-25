package instituicao

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sera_backend/pkg/service/asaas"
	"github.com/sera_backend/pkg/service/instituicao"
	"github.com/sera_backend/pkg/service/user"
)

func RegisterInstituicaoHandlers(r chi.Router, service instituicao.InstituicaoServiceInterface, userService user.UserServiceInterface, clientAsaas asaas.AsaasClientInterface) {
	r.Route("/api/v1/instituicao", func(r chi.Router) {
		r.Post("/add", createInstituicao(service, userService, clientAsaas))
		r.Put("/update/{id}/{nome}", updateInstituicao(service))
		r.Get("/getbyid/{id}", getByIdInstituicao(service))
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllInstituicao(service)
			handler.ServeHTTP(w, r)
		})
	})
}
