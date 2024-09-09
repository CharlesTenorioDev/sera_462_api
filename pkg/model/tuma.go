package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Turma struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	DataType  string             `bson:"data_type" json:"-"`
	Nome      string             `bson:"nome" json:"nome"`
	Enabled   bool               `bson:"enabled" json:"enabled"`
	Materias  []Materia          `bson:"materias" json:"materias"`
	ALunos    []Aluno            `bson:"alunos" json:"alunos"`
	CreatedAt string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at" json:"updated_at,omitempty"`
}
