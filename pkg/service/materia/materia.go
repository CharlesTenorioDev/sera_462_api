package materia

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

type MateriaServiceInterface interface {
	Create(ctx context.Context, Materia model.Materia) (*model.Materia, error)
	Update(ctx context.Context, ID string, meioToChange *model.Materia) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Materia, error)
	GetAll(ctx context.Context, filters model.FilterMateria, limit, page int64) (*model.Paginate, error)
	CheckExists(ctx context.Context, meio string) bool
}

type MateriaDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewMateriaService(mongo_connection mongodb.MongoDBInterface) *MateriaDataService {
	return &MateriaDataService{
		mdb: mongo_connection,
	}
}

func (mpg *MateriaDataService) Create(ctx context.Context, Materia model.Materia) (*model.Materia, error) {
	collection := mpg.mdb.GetCollection("cfStore")
	meio := model.NewMateria(Materia)

	result, err := collection.InsertOne(ctx, meio)
	if err != nil {
		logger.Error("erro salvar meio de pagamento", err)
		return &Materia, err
	}

	Materia.ID = result.InsertedID.(primitive.ObjectID)

	return &Materia, nil
}

func (mpg *MateriaDataService) Update(ctx context.Context, ID string, Materia *model.Materia) (bool, error) {
	collection := mpg.mdb.GetCollection("cfStore")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: Materia.Nome},
			{Key: "enabled", Value: Materia.Enabled},
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

func (mpg *MateriaDataService) GetByID(ctx context.Context, ID string) (*model.Materia, error) {

	collection := mpg.mdb.GetCollection("cfStore")

	Materia := &model.Materia{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Materia)
	if err != nil {
		logger.Error("erro ao consultar meio de pagamento", err)
		return nil, err
	}

	return Materia, nil
}

func (mpg *MateriaDataService) GetAll(ctx context.Context, filters model.FilterMateria, limit, page int64) (*model.Paginate, error) {
	collection := mpg.mdb.GetCollection("cfStore")

	query := bson.M{}

	if filters.Nome != "" || filters.Enabled != "" {
		if filters.Nome != "" {
			query["meio_pg"] = bson.M{"$regex": fmt.Sprintf(".*%s.*", filters.Nome), "$options": "i"}
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
		logger.Error("erro ao consulta todos Materia de pg", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Materia, 0)
	for curr.Next(ctx) {
		mpg := &model.Materia{}
		if err := curr.Decode(mpg); err != nil {
			logger.Error("erro ao consulta todos Materia de pg", err)
		}
		result = append(result, mpg)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (mpg *MateriaDataService) CheckExists(ctx context.Context, meio string) bool {
	collection := mpg.mdb.GetCollection("cfStore")

	query := bson.M{
		"data_type": "meio_pg",
		"meio_pg":   meio,
	}

	count, err := collection.CountDocuments(ctx, query)
	if err != nil {
		logger.Error("Erro ao verificar a existência do meio de pagamento", err)
		return false
	}

	return count > 0
}
