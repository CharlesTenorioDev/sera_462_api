package questionarioia

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/internal/dto"
	"github.com/sera_backend/pkg/adapter/mongodb"
	"github.com/sera_backend/pkg/adapter/rabbitmq"
	"github.com/sera_backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuestionarioServiceInterface interface {
	Create(ctx context.Context, Questionario model.Questionario) (*model.Questionario, error)
	Update(ctx context.Context, ID string, QuestionarioToChange *model.Questionario) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Questionario, error)
	GetAll(ctx context.Context, filters model.FilterQuestionario, limit, page int64) (*model.Paginate, error)
	GetByDocumento(ctx context.Context, Documento string) bool
}

type QuestionarioDataService struct {
	mdb mongodb.MongoDBInterface
	rmb rabbitmq.RabbitInterface
}

func NewQuestionarioervice(mongo_connection mongodb.MongoDBInterface, rabbitmq_connection rabbitmq.RabbitInterface) *QuestionarioDataService {
	return &QuestionarioDataService{
		mdb: mongo_connection,
		rmb: rabbitmq_connection,
	}
}

func (cat *QuestionarioDataService) Create(ctx context.Context, Questionario model.Questionario) (*model.Questionario, error) {
	collection := cat.mdb.GetCollection("cfSera")
	questFila := dto.QuestionarioParaFilaDTO{}
	questFila.ID = Questionario.ID
	questFila.Titulo = Questionario.Titulo

	cli := model.NewQuestionario(Questionario)
	result, err := collection.InsertOne(ctx, cli)
	if err != nil {
		logger.Error("erro salvar  Questionario", err)
		return &Questionario, err
	}

	jsonData, err := json.Marshal(questFila)
	if err != nil {
		logger.Error("Erro ao converter para JSON questinoario de envio p fila:", err)
		return &Questionario, err
	}

	msg := &rabbitmq.Message{
		Data:        jsonData,
		ContentType: "application/json; charset=utf-8",
	}

	err = cat.rmb.Connect()
	if err != nil {
		logger.Error("deu ruim na conexao como RabbitMQ", err)
	}

	err = cat.rmb.SenderRb(ctx, "amq.direct", "QUEUE_ENVIAR_IA", msg)
	if err != nil {
		logger.Error("Erro ao enviar Questionario para fila:", err)
		return &Questionario, err
	}

	cli.ID = result.InsertedID.(primitive.ObjectID)

	return cli, nil
}

func (cat *QuestionarioDataService) Update(ctx context.Context, ID string, Questionario *model.Questionario) (bool, error) {
	collection := cat.mdb.GetCollection("cfSera")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
		{Key: "data_type", Value: "Questionario"},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "titulo", Value: Questionario.Titulo},
			{Key: "enabled", Value: Questionario.Enabled},
			{Key: "updated_at", Value: time.Now().Format(time.RFC3339)},
		},
	}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		logger.Error("Error while updating data", err)

		return false, err
	}

	return true, nil
}

func (cat *QuestionarioDataService) GetByID(ctx context.Context, ID string) (*model.Questionario, error) {

	collection := cat.mdb.GetCollection("cfSera")

	Questionario := &model.Questionario{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "Questionario"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Questionario)
	if err != nil {
		logger.Error("erro ao consultar Questionario", err)
		return nil, err
	}

	return Questionario, nil
}

func (cat *QuestionarioDataService) GetAll(ctx context.Context, filters model.FilterQuestionario, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfSera")

	query := bson.M{"data_type": "Questionario"}

	if filters.Titulo != "" || filters.Enabled != "" {
		if filters.Titulo != "" {
			query["nome"] = bson.M{"$regex": fmt.Sprintf(".*%s.*", filters.Titulo), "$options": "i"}
		}
		if filters.Enabled != "" {
			enable, err := strconv.ParseBool(filters.Enabled)
			if err != nil {
				logger.Error("erro converter campo enabled", err)
				return nil, err
			}
			query["enabled"] = enable
		}
	}
	count, err := collection.CountDocuments(ctx, query, &options.CountOptions{})

	if err != nil {
		logger.Error("erro ao consultar todas as Alunos", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Aluno, 0)
	for curr.Next(ctx) {
		cat := &model.Aluno{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todas as Alunos", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (cat *QuestionarioDataService) GetByDocumento(ctx context.Context, Doc string) bool {

	collection := cat.mdb.GetCollection("cfSera")

	// Utilizando o método CountDocuments para verificar a existência
	filter := bson.D{
		{Key: "cpf_cnpj", Value: Doc},
		{Key: "data_type", Value: "Questionario"},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("erro ao consultar Questionario pelo doc", err)
		return false
	}

	// Se count for maior que zero, o fornecedor existe
	return count > 0
}
