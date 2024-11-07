package gemini

import (
	"github.com/go-chi/chi/v5"
	"github.com/sera_backend/pkg/service/gemini"
)

func RegisterGeminiAPIHandlers(r chi.Router, service gemini.GeminiServiceInterface) {
	r.Route("/api/v1/gemini", func(r chi.Router) {
		r.Post("/add", createQuestion(service))

	})
}
