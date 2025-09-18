package engineService

import (
	"context"

	"github.com/Tushar456/go-carzone/models"
	"github.com/Tushar456/go-carzone/repository"
	"go.opentelemetry.io/otel"
)

type EngineService struct {
	store repository.EngineRepositoryInterface
}

func NewEngineService(store repository.EngineRepositoryInterface) *EngineService {
	return &EngineService{
		store: store,
	}
}

func (es *EngineService) GetEngineById(ctx context.Context, id string) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "GetEngineById")
	defer span.End()
	engine, err := es.store.GetEngineById(ctx, id)
	if err != nil {
		return &models.Engine{}, err
	}
	return engine, nil
}

func (es *EngineService) CreateEngine(ctx context.Context, engine *models.EngineRequest) (*models.Engine, error) {

	ctx, span := otel.Tracer("engineservice").Start(ctx, "CreateEngine")
	defer span.End()

	if err := engine.Validate(); err != nil {
		return &models.Engine{}, err
	}
	createdEngine, err := es.store.CreateEngine(ctx, engine)
	if err != nil {
		return &models.Engine{}, err
	}
	return createdEngine, nil

}

func (es *EngineService) UpdateEngine(ctx context.Context, id string, engineRequest *models.EngineRequest) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "UpdateEngine")
	defer span.End()
	if err := engineRequest.Validate(); err != nil {
		return &models.Engine{}, err
	}
	updatedEngine, err := es.store.UpdateEngine(ctx, id, engineRequest)
	if err != nil {
		return &models.Engine{}, err
	}
	return updatedEngine, nil
}

func (es *EngineService) DeleteEngine(ctx context.Context, id string) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "DeleteEngine")
	defer span.End()
	deletedEngine, err := es.store.DeleteEngine(ctx, id)
	if err != nil {
		return &models.Engine{}, err
	}
	return deletedEngine, nil
}
