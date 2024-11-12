package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Questionario struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	IDTurma       primitive.ObjectID `bson:"turma_id" json:"turma_id"`
	IDMateria     primitive.ObjectID `bson:"materia_id" json:"materia_id"`
	IDProfessor   primitive.ObjectID `bson:"professor_id" json:"professor_id"`
	Titulo        string             `bson:"titulo" json:"titulo"`
	Questoes      []Questao          `bson:"questoes" json:"questoes"`
	EnvaidaParaIA bool               `bson:"enviado_ia" json:"enviado_ia"`
	CreatedAt     string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt     string             `bson:"updated_at" json:"updated_at,omitempty"`
}

type Questao struct {
	Pergunta     string   `bson:"pergunta" json:"pergunta"`
	Alternativas []string `bson:"alternativas" json:"alternativas"`
	Correta      string   `bson:"correta" json:"correta"`
}

type Resposta struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	IDQuestionario primitive.ObjectID `bson:"questionario_id" json:"questionario_id"`
	IDAluno        primitive.ObjectID `bson:"aluno_id" json:"aluno_id"`
	Respostas      map[string]string  `bson:"respostas" json:"respostas"`
	CreatedAt      string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      string             `bson:"updated_at" json:"updated_at,omitempty"` // Map of question to answer
}
