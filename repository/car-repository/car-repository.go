package carRepository

import (
	"context"
	"errors"

	"github.com/Tushar456/go-carzone/models"
	"github.com/Tushar456/go-carzone/repository"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type CarRepository struct {
	carRepo    *repository.Repository[models.Car]
	engineRepo *repository.Repository[models.Engine]
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{
		carRepo:    repository.New[models.Car](db),
		engineRepo: repository.New[models.Engine](db),
	}
}

func (s *CarRepository) GetCarById(ctx context.Context, id string) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "GetCarById")
	defer span.End()

	var car models.Car
	if err := s.carRepo.GetWithPreload(ctx, &car, []string{"Engine"}, "id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Note: Returning a non-nil empty struct on "not found" can be misleading.
			return &car, nil
		}
		return &car, err
	}
	return &car, nil
}

func (s *CarRepository) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "GetCarByBrand")
	defer span.End()

	var cars []models.Car
	var err error

	if isEngine {
		err = s.carRepo.FindWithPreload(ctx, &cars, []string{"Engine"}, "brand = ?", brand)
	} else {
		err = s.carRepo.Find(ctx, &cars, "brand = ?", brand)
	}

	if err != nil {
		return nil, err
	}

	return cars, nil
}

func (s *CarRepository) CreateCar(ctx context.Context, carRequest *models.CarRequest) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "CreateCar")
	defer span.End()

	var engine models.Engine
	if err := s.engineRepo.Get(ctx, &engine, "engine_id = ?", carRequest.EngineID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("engine not found")
		}
		return nil, err
	}

	car := &models.Car{
		ID:       uuid.New(),
		Name:     carRequest.Name,
		Year:     carRequest.Year,
		Brand:    carRequest.Brand,
		FuelType: carRequest.FuelType,
		EngineID: engine.EngineID, // Use the validated engine's ID
		Price:    carRequest.Price,
	}

	if err := s.carRepo.Create(ctx, car); err != nil {
		return nil, err
	}

	var createdCar models.Car
	if err := s.carRepo.GetWithPreload(ctx, &createdCar, []string{"Engine"}, "id = ?", car.ID); err != nil {
		return nil, err
	}

	return &createdCar, nil
}

func (s *CarRepository) UpdateCar(ctx context.Context, id string, updateCarRequest *models.CarRequest) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "UpdateCar")
	defer span.End()

	var car models.Car
	if err := s.carRepo.Get(ctx, &car, "id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("car not found")
		}
		return nil, err
	}

	// Validate that the new engine exists before updating.
	engineID, err := uuid.Parse(updateCarRequest.EngineID)
	if err != nil {
		return nil, err
	}
	var engine models.Engine
	if err := s.engineRepo.Get(ctx, &engine, "engine_id = ?", engineID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("engine not found")
		}
		return nil, err
	}

	car.Name = updateCarRequest.Name
	car.Year = updateCarRequest.Year
	car.Brand = updateCarRequest.Brand
	car.FuelType = updateCarRequest.FuelType
	car.Price = updateCarRequest.Price
	car.EngineID = engineID

	if err := s.carRepo.Update(ctx, &car); err != nil {
		return nil, err
	}

	// Reload the car with the engine association to return the full object.
	if err := s.carRepo.GetWithPreload(ctx, &car, []string{"Engine"}, "id = ?", car.ID); err != nil {
		return nil, err
	}

	return &car, nil
}

func (s *CarRepository) DeleteCar(ctx context.Context, id string) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "DeleteCar")
	defer span.End()

	var car models.Car
	// First, find the car to return it after deletion.
	if err := s.carRepo.GetWithPreload(ctx, &car, []string{"Engine"}, "id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("car not found")
		}
		return nil, err
	}

	if err := s.carRepo.Delete(ctx, &car); err != nil {
		return nil, err
	}
	return &car, nil
}
