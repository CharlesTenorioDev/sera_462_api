package model

import (
	"encoding/json"
	"time"

	"github.com/sera_backend/internal/config/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Aluno struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	ID_Instituicao primitive.ObjectID `bson:"instituicao_id" json:"instituicao_id"`
	Matricula      string             `bson:"matricula" json:"matricula"`
	DataType       string             `bson:"data_type" json:"-"`
	IDUsuario      primitive.ObjectID `bson:"user_id " json:"id_usr"`
	Nome           string             `bson:"nome" json:"nome"`
	DataNasc       time.Time          `bson:"data_nasc" json:"data_nasc"`
	Sexo           string             `bson:"sexo" json:"sexo"`
	Enabled        bool               `bson:"enabled" json:"enabled"`
	CreatedAt      string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      string             `bson:"updated_at" json:"updated_at,omitempty"`
}

func (c Aluno) AlunoConvet() string {
	data, err := json.Marshal(c)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterAluno struct {
	Nome           string             `json:"nome"`
	Matricula      string             `json:"matricula"`
	ID_Instituicao primitive.ObjectID `bson:"instituicao_id" json:"instituicao_id"`
	DataNasc       time.Time          `json:"data_nasc"`
	IDUsuario      primitive.ObjectID `bson:"user_id " json:"id_usr"`
	Documento      string             `bson:"documento" json:"documento"`
	Enabled        string             `json:"enabled"`
}

func NewAluno(Aluno_request Aluno) *Aluno {
	return &Aluno{
		ID:             primitive.NewObjectID(),
		ID_Instituicao: Aluno_request.ID_Instituicao,
		Matricula:      Aluno_request.Matricula,
		DataType:       "Aluno",
		IDUsuario:      Aluno_request.IDUsuario,
		Nome:           Aluno_request.Nome,
		DataNasc:       Aluno_request.DataNasc,
		Enabled:        true,
		CreatedAt:      time.Now().String(),
	}
}
