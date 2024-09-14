package model

import (
	"encoding/json"
	"time"

	"github.com/sera_backend/internal/config/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Turma struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ID_Instituicao primitive.ObjectID `bson:"instituicao_id" json:"instituicao_id"`
	DataType       string             `bson:"data_type" json:"-"`
	Nome           string             `bson:"nome" json:"nome"`
	Horario        string             `bson:"horario" json:"horario"`
	Materias       []Materia          `bson:"materias" json:"materias"`
	Turmas         []Turma            `bson:"Turmas" json:"Turmas"`
	Enabled        bool               `bson:"enabled" json:"enabled"`
	CreatedAt      string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      string             `bson:"updated_at" json:"updated_at,omitempty"`
}

func (t Turma) TurmaConvet() string {
	data, err := json.Marshal(t)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterTurma struct {
	Nome           string             `json:"nome"`
	DataType       string             `bson:"data_type" json:"-"`
	Horario        string             `json:"horario"`
	ID_Instituicao primitive.ObjectID `bson:"instituicao_id" json:"instituicao_id"`
	Enabled        string             `json:"enabled"`
}

func NewTurma(Turma_request Turma) *Turma {
	return &Turma{
		ID:             primitive.NewObjectID(),
		ID_Instituicao: Turma_request.ID_Instituicao,
		DataType:       "Turma",
		Nome:           Turma_request.Nome,
		Horario:        Turma_request.Horario,
		Materias:       Turma_request.Materias,
		Turmas:         Turma_request.Turmas,
		Enabled:        true,
		CreatedAt:      time.Now().String(),
	}
}
