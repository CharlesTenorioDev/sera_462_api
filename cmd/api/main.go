package main

import (
	"log"
	"net/http"

	"github.com/sera_backend/internal/config"
	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/pkg/adapter/rabbitmq"

	hand_asaas "github.com/sera_backend/internal/handler/asaas"

	hand_instituicao "github.com/sera_backend/internal/handler/instituicao"
	hand_questionarioia "github.com/sera_backend/internal/handler/questionarioia"
	hand_usr "github.com/sera_backend/internal/handler/user"
	handHealthcheck "github.com/sera_backend/internal/healthcheck"

	"github.com/sera_backend/pkg/adapter/mongodb"

	"github.com/sera_backend/pkg/server"

	service_asaas "github.com/sera_backend/pkg/service/asaas"

	serviceHealthcheck "github.com/sera_backend/pkg/service/healthcheck"
	service_instituicao "github.com/sera_backend/pkg/service/instituicao"
	service_questionarioia "github.com/sera_backend/pkg/service/questionarioia"
	service_usr "github.com/sera_backend/pkg/service/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	VERSION = "0.1.0-dev"
	COMMIT  = "ABCDEFG-dev"
)

func main() {

	logger.Info("start Application Sera 462 API")
	conf := config.NewConfig()

	mogDbConn := mongodb.New(conf)

	fila := []rabbitmq.Fila{
		{
			Name:       "QUEUE_ENVIAR_IA",
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
		},
	}

	rbtMQConn := rabbitmq.NewRabbitMQ(fila, conf)

	if err := rbtMQConn.Connect(); err != nil {
		logger.Error("Falha ao conectar no RabbitMQ", err)

	}

	if err := rbtMQConn.DeclareQueues(); err != nil {
		logger.Error("Falha ao declarar filas no RabbitMQ", err)

	}

	isAlive, err := rbtMQConn.IsAlive()
	if !isAlive || err != nil {
		logger.Error("RabbitMQ não está disponível", err)

	}

	usr_service := service_usr.NewUsuarioservice(mogDbConn)

	inst_service := service_instituicao.NewInstituicaoervice(mogDbConn)
	asaas_service := service_asaas.NewClient(conf)

	handServiceHealthcheck := serviceHealthcheck.NewHealthcheckService(mogDbConn)

	questionarioia_service := service_questionarioia.NewQuestionarioervice(mogDbConn, rbtMQConn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", conf.TokenAuth))
	r.Use(middleware.WithValue("JWTTokenExp", conf.JWTTokenExp))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/", healthcheck)
	hand_usr.RegisterUsuarioAPIHandlers(r, usr_service)
	hand_instituicao.RegisterInstituicaoHandlers(r, inst_service, usr_service)
	hand_asaas.RegisterAsaasHandlers(r, asaas_service)
	handHealthcheck.RegisterHealthcheckAPIHandlers(r, handServiceHealthcheck)

	hand_questionarioia.RegisterQuestionarioHandlers(r, questionarioia_service)

	// Inicie o worker em uma goroutine

	srv := server.NewHTTPServer(r, conf)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Server Run on [Port: %s], [Mode: %s], [Version: %s], [Commit: %s]", conf.PORT, conf.Mode, VERSION, COMMIT)

	select {}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"MSG": "Server Ok", "codigo": 200}`))
}
