package aluno

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

type AlunoServiceInterface interface {
	Create(ctx context.Context, Aluno model.Aluno) (*model.Aluno, error)
	Update(ctx context.Context, ID string, AlunoToChange *model.Aluno) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Aluno, error)
	GetAll(ctx context.Context, filters model.FilterAluno, limit, page int64) (*model.Paginate, error)
	GetByDocumento(ctx context.Context, Documento string) bool
}

type AlunoDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewAlunoervice(mongo_connection mongodb.MongoDBInterface) *AlunoDataService {
	return &AlunoDataService{
		mdb: mongo_connection,
	}
}

func (cat *AlunoDataService) Create(ctx context.Context, Aluno model.Aluno) (*model.Aluno, error) {
	collection := cat.mdb.GetCollection("cfStore")
	cli := model.NewAluno(Aluno)
	result, err := collection.InsertOne(ctx, cli)
	if err != nil {
		logger.Error("erro salvar  Aluno", err)
		return &Aluno, err
	}

	cli.ID = result.InsertedID.(primitive.ObjectID)

	return cli, nil
}

func (cat *AlunoDataService) Update(ctx context.Context, ID string, Aluno *model.Aluno) (bool, error) {
	collection := cat.mdb.GetCollection("cfStore")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
		{Key: "data_type", Value: "Aluno"},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: Aluno.Nome},
			{Key: "enabled", Value: Aluno.Enabled},
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

func (cat *AlunoDataService) GetByID(ctx context.Context, ID string) (*model.Aluno, error) {

	collection := cat.mdb.GetCollection("cfStore")

	Aluno := &model.Aluno{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "Aluno"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Aluno)
	if err != nil {
		logger.Error("erro ao consultar Aluno", err)
		return nil, err
	}

	return Aluno, nil
}

func (cat *AlunoDataService) GetAll(ctx context.Context, filters model.FilterAluno, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfStore")

	query := bson.M{"data_type": "Aluno"}

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

func (cat *AlunoDataService) GetByDocumento(ctx context.Context, Doc string) bool {

	collection := cat.mdb.GetCollection("cfStore")

	// Utilizando o método CountDocuments para verificar a existência
	filter := bson.D{
		{Key: "cpf_cnpj", Value: Doc},
		{Key: "data_type", Value: "Aluno"},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("erro ao consultar Aluno pelo doc", err)
		return false
	}

	// Se count for maior que zero, o fornecedor existe
	return count > 0
}
