package healthcheck

import (
	"github.com/sera_backend/pkg/adapter/mongodb"
)

type HealthcheckServiceInterface interface {
	CheckDB() (bool, error)
}

type HealthcheckService struct {
	mdb mongodb.MongoDBInterface
}

func NewHealthcheckService(mdb mongodb.MongoDBInterface) *HealthcheckService {
	return &HealthcheckService{
		mdb: mdb,
	}
}

func (h *HealthcheckService) CheckDB() (bool, error) {
	return h.mdb.CheckDB()
}
