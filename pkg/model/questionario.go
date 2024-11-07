package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Questionario struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	IDInstituicao primitive.ObjectID   `bson:"id_instituicao" json:"id_instituicao"`
	IDMateria     primitive.ObjectID   `bson:"id_materia" json:"id_materia"`
	Assunto       string               `bson:"assunto" json:"assunto"`
	Turmas        []primitive.ObjectID `bson:"turmas" json:"turmas"`
	Responsavel   string               `bson:"responsavel" json:"responsavel"`
	QtdTextoLivre int                  `bson:"qtd_texto_livre" json:"qtd_texto_livre"`
	QtdMultiplas  int                  `bson:"qtd_multiplas" json:"qtd_multiplas"`
	Perguntas     []Pergunta           `bson:"perguntas" json:"perguntas"`
	Enabled       bool                 `bson:"enabled" json:"enabled"`
	CreatedAt     string               `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt     string               `bson:"updated_at" json:"updated_at,omitempty"`
}

type Pergunta struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Tipo            string             `bson:"tipo" json:"tipo"` // "texto_livre" ou "multipla_escolha"
	Conteudo        string             `bson:"conteudo" json:"conteudo"`
	Opcoes          []string           `bson:"opcoes,omitempty" json:"opcoes,omitempty"` // Usado apenas para múltipla escolha
	RespostaCorreta string             `bson:"resposta_correta" json:"resposta_correta"`
	CreatedAt       string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt       string             `bson:"updated_at" json:"updated_at,omitempty"`
}

type Resposta struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	IDAluno        primitive.ObjectID `bson:"id_aluno" json:"id_aluno"`
	IDQuestionario primitive.ObjectID `bson:"id_questionario" json:"id_questionario"`
	Respostas      []RespostaPergunta `bson:"respostas" json:"respostas"`
	CreatedAt      string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      string             `bson:"updated_at" json:"updated_at,omitempty"`
}

type RespostaPergunta struct {
	IDPergunta   primitive.ObjectID `bson:"id_pergunta" json:"id_pergunta"`
	RespostaDada string             `bson:"resposta_dada" json:"resposta_dada"`
	Correta      bool               `bson:"correta" json:"correta"`
}

type PerguntaGeminai struct {
	Perguntas string `json:"perguntas"`
}
