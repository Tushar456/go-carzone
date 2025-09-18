package repository

import (
	"context"

	"github.com/Tushar456/go-carzone/models"
)

type CarRepositoryInterface interface {
	GetCarById(ctx context.Context, id string) (*models.Car, error)
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	CreateCar(ctx context.Context, car *models.CarRequest) (*models.Car, error)
	UpdateCar(ctx context.Context, id string, updateCar *models.CarRequest) (*models.Car, error)
	DeleteCar(ctx context.Context, id string) (*models.Car, error)
}

type EngineRepositoryInterface interface {
	GetEngineById(ctx context.Context, id string) (*models.Engine, error)
	CreateEngine(ctx context.Context, engine *models.EngineRequest) (*models.Engine, error)
	UpdateEngine(ctx context.Context, id string, updateEngine *models.EngineRequest) (*models.Engine, error)
	DeleteEngine(ctx context.Context, id string) (*models.Engine, error)
}
