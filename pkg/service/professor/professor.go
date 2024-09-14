package professor

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

type ProfessorServiceInterface interface {
	Create(ctx context.Context, Professor model.Professor) (*model.Professor, error)
	Update(ctx context.Context, ID string, ProfessorToChange *model.Professor) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Professor, error)
	GetAll(ctx context.Context, filters model.FilterProfessor, limit, page int64) (*model.Paginate, error)
	GetByDocumento(ctx context.Context, Documento string) bool
}

type ProfessorDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewProfessorervice(mongo_connection mongodb.MongoDBInterface) *ProfessorDataService {
	return &ProfessorDataService{
		mdb: mongo_connection,
	}
}

func (cat *ProfessorDataService) Create(ctx context.Context, Professor model.Professor) (*model.Professor, error) {
	collection := cat.mdb.GetCollection("cfStore")
	cli := model.NewProfessor(Professor)
	result, err := collection.InsertOne(ctx, cli)
	if err != nil {
		logger.Error("erro salvar  Professor", err)
		return &Professor, err
	}

	cli.ID = result.InsertedID.(primitive.ObjectID)

	return cli, nil
}

func (cat *ProfessorDataService) Update(ctx context.Context, ID string, Professor *model.Professor) (bool, error) {
	collection := cat.mdb.GetCollection("cfStore")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
		{Key: "data_type", Value: "Professor"},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: Professor.Nome},
			{Key: "enabled", Value: Professor.Enabled},
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

func (cat *ProfessorDataService) GetByID(ctx context.Context, ID string) (*model.Professor, error) {

	collection := cat.mdb.GetCollection("cfStore")

	Professor := &model.Professor{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "Professor"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Professor)
	if err != nil {
		logger.Error("erro ao consultar Professor", err)
		return nil, err
	}

	return Professor, nil
}

func (cat *ProfessorDataService) GetAll(ctx context.Context, filters model.FilterProfessor, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfStore")

	query := bson.M{"data_type": "Professor"}

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
		logger.Error("erro ao consultar todas as Professors", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Professor, 0)
	for curr.Next(ctx) {
		cat := &model.Professor{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todas as Professors", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (cat *ProfessorDataService) GetByDocumento(ctx context.Context, Doc string) bool {

	collection := cat.mdb.GetCollection("cfStore")

	// Utilizando o método CountDocuments para verificar a existência
	filter := bson.D{
		{Key: "cpf_cnpj", Value: Doc},
		{Key: "data_type", Value: "Professor"},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("erro ao consultar Professor pelo doc", err)
		return false
	}

	// Se count for maior que zero, o fornecedor existe
	return count > 0
}
