package engineRepository

import (
	"context"

	"errors"

	"github.com/Tushar456/go-carzone/models"
	"github.com/Tushar456/go-carzone/repository"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type EngineRepository struct {
	repo *repository.Repository[models.Engine]
}

func NewEngineRepository(db *gorm.DB) *EngineRepository {
	return &EngineRepository{
		repo: repository.New[models.Engine](db),
	}
}

func (s *EngineRepository) GetEngineById(ctx context.Context, id string) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "GetEngineById")
	defer span.End()

	var engine models.Engine
	if err := s.repo.Get(ctx, &engine, "engine_id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Note: Returning an empty struct on "not found" can be misleading.
			// Consider returning (nil, gorm.ErrRecordNotFound) or (nil, nil).
			return &models.Engine{}, nil
		}
		return &models.Engine{}, err
	}
	return &engine, nil
}

func (s *EngineRepository) CreateEngine(ctx context.Context, engineRequest *models.EngineRequest) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "CreateEngine")
	defer span.End()

	engine := &models.Engine{
		EngineID:      uuid.New(),
		Displacement:  engineRequest.Displacement,
		NoOfCylinders: engineRequest.NoOfCylinders,
		CarRange:      engineRequest.CarRange,
	}

	if err := s.repo.Create(ctx, engine); err != nil {
		return &models.Engine{}, err
	}

	return engine, nil
}

func (s *EngineRepository) UpdateEngine(ctx context.Context, id string, engineRequest *models.EngineRequest) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "UpdateEngine")
	defer span.End()

	var engine models.Engine
	if err := s.repo.Get(ctx, &engine, "engine_id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Engine{}, errors.New("engine not found")
		}
		return &models.Engine{}, err
	}

	engine.Displacement = engineRequest.Displacement
	engine.NoOfCylinders = engineRequest.NoOfCylinders
	engine.CarRange = engineRequest.CarRange

	if err := s.repo.Update(ctx, &engine); err != nil {
		return &models.Engine{}, err
	}

	return &engine, nil
}

func (s *EngineRepository) DeleteEngine(ctx context.Context, id string) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "DeleteEngine")
	defer span.End()

	var engine models.Engine
	// The original implementation had a bug here using "id = ?" instead of "engine_id = ?".
	// I've corrected it to use the correct primary key column.
	if err := s.repo.Get(ctx, &engine, "engine_id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Engine{}, errors.New("engine not found")
		}
		return &models.Engine{}, err
	}

	if err := s.repo.Delete(ctx, &engine); err != nil {
		return &models.Engine{}, err
	}

	return &engine, nil
}
