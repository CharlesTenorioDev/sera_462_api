package model

import (
	"encoding/json"
	"time"

	"github.com/sera_backend/internal/config/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Materia struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DataType string             `bson:"data_type" json:"-"`
	Nome     string             `bson:"nome" json:"nome"`
	Enabled  bool               `bson:"enabled" json:"enabled"`

	CreatedAt string `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt string `bson:"updated_at" json:"updated_at,omitempty"`
}

func (c Materia) MateriaConvet() string {
	data, err := json.Marshal(c)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterMateria struct {
	Nome    string `json:"nome"`
	Enabled string `json:"enabled"`
}

func NewMateria(client_request Materia) *Materia {
	return &Materia{
		ID:       primitive.NewObjectID(),
		DataType: "Materia",
		Nome:     client_request.Nome,
		Enabled:  true,

		CreatedAt: time.Now().String(),
	}
}
