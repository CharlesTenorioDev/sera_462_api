package model

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Questionario struct {
	ID              primitive.ObjectID
	IDTurma         primitive.ObjectID `bson:"turma_id" json:"turma_id"`
	IDMateria       primitive.ObjectID `bson:"materia_id" json:"materia_id"`
	IDProfessor     primitive.ObjectID `bson:"professor_id" json:"professor_id"`
	DataType        string             `bson:"data_type" json:"-"`
	Titulo          string             `bson:"titulo" json:"titulo"`
	Quantidade      int                `bson:"quantidade" json:"quantidade"`
	Questoes        []Questao          `bson:"questoes" json:"questoes"`
	EnvaidaParaIA   bool               `bson:"enviado_ia" json:"enviado_ia"`
	EnviadaParaFila bool               `bson:"enviado_fila" json:"enviado_fila"`
	CreatedAt       string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt       string             `bson:"updated_at" json:"updated_at,omitempty"`
	Enabled         bool               `bson:"enabled" json:"enabled"`
	PerguntarParaIA string             `bson:"perguntar_para_ia" json:"perguntar_para_ia"`
}

type Questao struct {
	Pergunta     string   `bson:"pergunta" json:"pergunta"`
	DataType     string   `bson:"data_type" json:"-"`
	Alternativas []string `bson:"alternativas" json:"alternativas"`
	Correta      string   `bson:"correta" json:"correta"`
}

type Resposta struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	IDQuestionario primitive.ObjectID `bson:"questionario_id" json:"questionario_id"`
	DataType       string             `bson:"data_type" json:"-"`
	IDAluno        primitive.ObjectID `bson:"aluno_id" json:"aluno_id"`
	Respostas      map[string]string  `bson:"respostas" json:"respostas"`
	CreatedAt      string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      string             `bson:"updated_at" json:"updated_at,omitempty"` // Map of question to answer
}

type RespostaAluno struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	DataType       string             `bson:"data_type" json:"-"`
	IDQuestionario primitive.ObjectID `bson:"questionario_id" json:"questionario_id"`
	IDAluno        primitive.ObjectID `bson:"aluno_id" json:"aluno_id"`
	Respostas      map[string]string  `bson:"respostas" json:"respostas"`
	CreatedAt      string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt      string             `bson:"updated_at" json:"updated_at,omitempty"` // Map of question to answer
}

type FilterQuestionario struct {
	ID          primitive.ObjectID
	IDTurma     primitive.ObjectID `bson:"turma_id" json:"turma_id"`
	DataType    string             `bson:"data_type" json:"-"`
	IDMateria   primitive.ObjectID `bson:"materia_id" json:"materia_id"`
	IDProfessor primitive.ObjectID `bson:"professor_id" json:"professor_id"`
	Titulo      string             `json:"titulo"`
	Enabled     string             `json:"enabled"`
}

func NewQuestionario(questionario_request Questionario) *Questionario {
	return &Questionario{
		ID:              primitive.NewObjectID(),
		IDTurma:         questionario_request.IDTurma,
		IDMateria:       questionario_request.IDMateria,
		IDProfessor:     questionario_request.IDProfessor,
		DataType:        "Questionario",
		Titulo:          strings.ToUpper(questionario_request.Titulo),
		Quantidade:      questionario_request.Quantidade,
		Questoes:        questionario_request.Questoes,
		EnvaidaParaIA:   questionario_request.EnvaidaParaIA,
		PerguntarParaIA: montarPerguntaParaIA(questionario_request.Titulo, questionario_request.Quantidade),
		CreatedAt:       time.Now().String(),
	}
}

func montarPerguntaParaIA(titulo string, quantidade int) string {
	var builder strings.Builder

	builder.WriteString("onde Cada pergunta deve ser representada em formato JSON, ")
	builder.WriteString("onde cada pergunta é uma tag principal contendo um array de respostas associadas. ")
	builder.WriteString("O JSON deve ser estruturado de forma aninhada, com cada pergunta como um objeto que inclui o texto da pergunta, ")
	builder.WriteString("as opcoes de resposta, e a identificacao da resposta correta?")
	resultStr := builder.String()

	pergunta_ia := ""
	resposta := "composto por"
	switch quantidade {

	case 1:
		resposta = " composto por uma pergunta "
	case 2:
		resposta = " composto por duas perguntas"
	case 3:
		resposta = " composto por três perguntas"
	case 4:
		resposta = " composto por quatro perguntas"
	case 5:
		resposta = " composto por cinco perguntas"
	case 6:
		resposta = " composto por seis perguntas"
	case 7:
		resposta = " composto por sete perguntas"
	case 8:
		resposta = " composto por oito perguntas"
	case 9:
		resposta = " composto por nove perguntas"
	case 10:
		resposta = " composto por dez perguntas"
	}

	pergunta_ia = "Gostaria de um questionário sobre" + " " + strings.ToUpper(titulo) + " " + resposta + " " + resultStr
	return pergunta_ia
}
