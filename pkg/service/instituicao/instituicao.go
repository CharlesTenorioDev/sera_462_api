package instituicao

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

type InstituicaoServiceInterface interface {
	Create(ctx context.Context, Instituicao model.Instituicao) (*model.Instituicao, error)
	Update(ctx context.Context, ID string, InstituicaoToChange *model.Instituicao) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Instituicao, error)
	GetAll(ctx context.Context, filters model.FilterInstituicao, limit, page int64) (*model.Paginate, error)
	GetByDocumento(ctx context.Context, Documento string) bool
}

type InstituicaoDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewInstituicaoervice(mongo_connection mongodb.MongoDBInterface) *InstituicaoDataService {
	return &InstituicaoDataService{
		mdb: mongo_connection,
	}
}

func (cat *InstituicaoDataService) Create(ctx context.Context, Instituicao model.Instituicao) (*model.Instituicao, error) {
	collection := cat.mdb.GetCollection("cfSera")
	cli := model.NewIntituicao(Instituicao)
	result, err := collection.InsertOne(ctx, cli)
	if err != nil {
		logger.Error("erro salvar  Instituicao", err)
		return &Instituicao, err
	}

	cli.ID = result.InsertedID.(primitive.ObjectID)

	return cli, nil
}

func (cat *InstituicaoDataService) Update(ctx context.Context, ID string, Instituicao *model.Instituicao) (bool, error) {
	collection := cat.mdb.GetCollection("cfSera")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
		{Key: "data_type", Value: "Instituicao"},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: Instituicao.Nome},
			{Key: "enabled", Value: Instituicao.Enabled},
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

func (cat *InstituicaoDataService) GetByID(ctx context.Context, ID string) (*model.Instituicao, error) {

	collection := cat.mdb.GetCollection("cfSera")

	Instituicao := &model.Instituicao{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "Instituicao"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Instituicao)
	if err != nil {
		logger.Error("erro ao consultar Instituicao", err)
		return nil, err
	}

	return Instituicao, nil
}

func (cat *InstituicaoDataService) GetAll(ctx context.Context, filters model.FilterInstituicao, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfSera")

	query := bson.M{"data_type": "Instituicao"}

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
		logger.Error("erro ao consultar todas as Instituicaos", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Instituicao, 0)
	for curr.Next(ctx) {
		cat := &model.Instituicao{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todas as Instituicaos", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (cat *InstituicaoDataService) GetByDocumento(ctx context.Context, Doc string) bool {

	collection := cat.mdb.GetCollection("cfSera")

	// Utilizando o método CountDocuments para verificar a existência
	filter := bson.D{
		{Key: "cnpj", Value: Doc},
		{Key: "data_type", Value: "Instituicao"},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("erro ao consultar Instituicao pelo doc", err)
		return false
	}

	// Se count for maior que zero, o fornecedor existe
	return count > 0
}
