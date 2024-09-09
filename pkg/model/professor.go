package model

import (
	"encoding/json"
	"time"

	"github.com/sera_backend/internal/config/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Professor struct {
	ID             primitive.ObjectID   `bson:"_id" json:"_id" validate:"required"`
	ID_instituicao primitive.ObjectID   `bson:"instituicao_id" json:"instituicao_id" validate:"required"`
	DataType       string               `bson:"data_type" json:"-"`
	IDUsuario      primitive.ObjectID   `bson:"user_id " json:"id_usr" validate:"required"`
	Nome           string               `bson:"nome" json:"nome" validate:"required"`
	Email          string               `bson:"email" json:"email" validate:"required,email"`
	Sexo           string               `bson:"sexo" json:"sexo"`
	Telefone       string               `bson:"telefone" json:"telefone"`
	Tipo           string               `bson:"tipo" json:"tipo"`
	Documento      string               `bson:"cpf_cnpj" json:"cpf_cnpj"`
	Materias       []primitive.ObjectID `bson:"materias" json:"materias"`
	Enabled        bool                 `bson:"enabled" json:"enabled"`
	CreatedAt      string               `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      string               `bson:"updated_at" json:"updated_at,omitempty"`
}

func (c Professor) ProfessorConvet() string {
	data, err := json.Marshal(c)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterProfessor struct {
	Nome      string             `json:"nome"`
	IDUsuario primitive.ObjectID `bson:"user_id " json:"id_usr"`
	Documento string             `bson:"documento" json:"documento"`
	Enabled   string             `json:"enabled"`
}

func NewProfessor(Professor_request Professor) *Professor {
	return &Professor{
		ID:        primitive.NewObjectID(),
		DataType:  "Professor",
		IDUsuario: Professor_request.IDUsuario,
		Nome:      Professor_request.Nome,
		Materias:  Professor_request.Materias,
		Enabled:   true,
		CreatedAt: time.Now().String(),
	}
}
