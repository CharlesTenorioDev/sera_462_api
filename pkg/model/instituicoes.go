package model

import (
	"encoding/json"
	"time"

	"github.com/sera_backend/internal/config/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instituicao struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	DataType  string             `bson:"data_type" json:"-"`
	IDUsuario primitive.ObjectID `bson:"user_id " json:"id_usr"`
	Nome      string             `bson:"nome" json:"nome"`
	Telefone  string             `bson:"telefone" json:"telefone"`
	CNPJ      string             `bson:"cnpj" json:"cnpj"`
	Endereco  Endereco           `bson:"endereco" json:"endereco"`
	Enabled   bool               `bson:"enabled" json:"enabled"`
	CreatedAt string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at" json:"updated_at,omitempty"`
}

func (f Instituicao) IntitucaoConvet() string {
	data, err := json.Marshal(f)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterInstituicao struct {
	Nome      string             `json:"nome"`
	IDUsuario primitive.ObjectID `bson:"user_id " json:"id_usr"`
	CNPJ      string             `bson:"cnpj" json:"cnpj"`
	Enabled   string             `json:"enabled"`
}

func NewIntituicao(instituicao_request Instituicao) *Instituicao {
	return &Instituicao{
		ID:        primitive.NewObjectID(),
		DataType:  "Institucao",
		IDUsuario: instituicao_request.IDUsuario,
		Nome:      instituicao_request.Nome,
		CNPJ:      instituicao_request.CNPJ,
		Enabled:   true,
		CreatedAt: time.Now().String(),
	}
}
