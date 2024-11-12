package turma

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sera_backend/internal/config/logger"
	"github.com/sera_backend/pkg/adapter/mongodb"
	"github.com/sera_backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TurmaServiceInterface interface {
	Create(ctx context.Context, Turma model.Turma) (*model.Turma, error)
	Update(ctx context.Context, ID string, TurmaToChange *model.Turma) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Turma, error)
	GetAll(ctx context.Context, filters model.FilterTurma, limit, page int64) (*model.Paginate, error)
	GetByDocumento(ctx context.Context, Documento string) bool
}

type TurmaDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewTurmaervice(mongo_connection mongodb.MongoDBInterface) *TurmaDataService {
	return &TurmaDataService{
		mdb: mongo_connection,
	}
}

func (cat *TurmaDataService) Create(ctx context.Context, Turma model.Turma) (*model.Turma, error) {
	collection := cat.mdb.GetCollection("cfSera")
	cli := model.NewTurma(Turma)
	result, err := collection.InsertOne(ctx, cli)
	if err != nil {
		logger.Error("erro salvar  Turma", err)
		return &Turma, err
	}

	cli.ID = result.InsertedID.(primitive.ObjectID)

	return cli, nil
}

func (cat *TurmaDataService) Update(ctx context.Context, ID string, Turma *model.Turma) (bool, error) {
	collection := cat.mdb.GetCollection("cfSera")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
		{Key: "data_type", Value: "Turma"},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: Turma.Nome},
			{Key: "enabled", Value: Turma.Enabled},
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

func (cat *TurmaDataService) GetByID(ctx context.Context, ID string) (*model.Turma, error) {

	collection := cat.mdb.GetCollection("cfSera")

	Turma := &model.Turma{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "Turma"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Turma)
	if err != nil {
		logger.Error("erro ao consultar Turma", err)
		return nil, err
	}

	return Turma, nil
}

func (cat *TurmaDataService) GetAll(ctx context.Context, filters model.FilterTurma, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfSera")

	query := bson.M{"data_type": "Turma"}

	if filters.Nome != "" || filters.Enabled != "" {
		if filters.Nome != "" {
			query["nome"] = bson.M{"$regex": fmt.Sprintf(".*%s.*", filters.Nome), "$options": "i"}
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
		logger.Error("erro ao consultar todas as Turmas", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Turma, 0)
	for curr.Next(ctx) {
		cat := &model.Turma{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todas as Turmas", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (cat *TurmaDataService) GetByDocumento(ctx context.Context, Doc string) bool {

	collection := cat.mdb.GetCollection("cfSera")

	// Utilizando o método CountDocuments para verificar a existência
	filter := bson.D{
		{Key: "cpf_cnpj", Value: Doc},
		{Key: "data_type", Value: "Turma"},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("erro ao consultar Turma pelo doc", err)
		return false
	}

	// Se count for maior que zero, o fornecedor existe
	return count > 0
}
