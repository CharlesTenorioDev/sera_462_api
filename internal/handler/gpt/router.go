package gpt

import (
	"github.com/go-chi/chi/v5"
	"github.com/sera_backend/pkg/service/gpt"
)

func RegisterGPTAPIHandlers(r chi.Router, service gpt.GptClientInterface) {
	r.Route("/api/v1/gpt", func(r chi.Router) {
		r.Post("/add", createQuestion(service))

	})
}
