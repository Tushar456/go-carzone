package carService

import (
	"context"

	"github.com/Tushar456/go-carzone/models"
	"github.com/Tushar456/go-carzone/repository"
	"go.opentelemetry.io/otel"
)

type CarService struct {
	store repository.CarRepositoryInterface
}

func NewCarService(store repository.CarRepositoryInterface) *CarService {
	return &CarService{
		store: store,
	}
}

func (cs *CarService) GetCarById(ctx context.Context, id string) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "GetCarById")
	defer span.End()
	car, err := cs.store.GetCarById(ctx, id)
	if err != nil {
		return &models.Car{}, err
	}
	return car, nil
}

func (cs *CarService) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "GetCarByBrand")
	defer span.End()
	cars, err := cs.store.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		return []models.Car{}, err
	}
	return cars, nil
}

func (cs *CarService) CreateCar(ctx context.Context, car *models.CarRequest) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "CreateCar")
	defer span.End()

	if err := car.Validate(); err != nil {
		return &models.Car{}, err
	}
	createdCar, err := cs.store.CreateCar(ctx, car)
	if err != nil {
		return &models.Car{}, err
	}
	return createdCar, nil

}

func (cs *CarService) UpdateCar(ctx context.Context, id string, carRequest *models.CarRequest) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "UpdateCar")
	defer span.End()
	if err := carRequest.Validate(); err != nil {
		return &models.Car{}, err
	}
	car, err := cs.store.UpdateCar(ctx, id, carRequest)
	if err != nil {
		return &models.Car{}, err
	}
	return car, nil

}

func (cs *CarService) DeleteCar(ctx context.Context, id string) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "DeleteCar")
	defer span.End()
	car, err := cs.store.DeleteCar(ctx, id)
	if err != nil {
		return &models.Car{}, err
	}
	return car, nil
}
