package store

import (
	"context"
	"errors"

	"github.com/Tushar456/go-carzone/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

// type CarStore struct {
// 	db *sql.DB
// }

// func NewCarStore(db *sql.DB) *CarStore {
// 	return &CarStore{
// 		db: db,
// 	}
// }

type CarStore struct {
	db *gorm.DB
}

func NewCarStore(db *gorm.DB) *CarStore {
	return &CarStore{
		db: db,
	}
}

// func (s *CarStore) GetCarById(ctx context.Context, id string) (*models.Car, error) {

// 	var car models.Car

// 	query := `Select c.id, c.brand, c.name, c.year, c.price, c.fuel_type, c.created_at, c.updated_at, c.engine_id, e.id, e.displacement, e.no_of_cylinders, e.car_range
// 	 			FROM car c
// 				LEFT JOIN
// 				engine e ON c.engine_id = e.id WHERE c.id=$1`

// 	err := s.db.QueryRowContext(ctx, query, id).Scan(
// 		&car.ID,
// 		&car.Brand,
// 		&car.Name,
// 		&car.Year,
// 		&car.Price,
// 		&car.FuelType,
// 		&car.CreatedAt,
// 		&car.UpdatedAt,
// 		&car.Engine.EngineID,
// 		&car.Engine.EngineID,
// 		&car.Engine.Displacement,
// 		&car.Engine.NoOfCylinders,
// 		&car.Engine.CarRange,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return &models.Car{}, nil
// 		}
// 		return &models.Car{}, err
// 	}
// 	return &car, nil
// }

func (s *CarStore) GetCarById(ctx context.Context, id string) (*models.Car, error) {
	ctx, span := otel.Tracer("carstore").Start(ctx, "GetCarById")
	defer span.End()
	var car models.Car

	if err := s.db.WithContext(ctx).Preload("Engine").First(&car, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &car, nil
		}
		return &car, err
	}
	return &car, nil
}

// func (s *CarStore) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {

// 	var cars []models.Car

// 	var query string

// 	if isEngine {
// 		query = `Select c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at,
// 		e.id, e.displacement, e.no_of_cylinders, e.car_range
// 		FROM car c LEFT JOIN engine e ON c.engine_id = e.id WHERE c.brand=$1`
// 	} else {
// 		query = `Select  c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at
// 		FROM car c WHERE c.brand=$1`
// 	}

// 	rows, err := s.db.QueryContext(ctx, query, brand)
// 	if err != nil {
// 		return cars, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var car models.Car
// 		if isEngine {
// 			err = rows.Scan(
// 				&car.ID,
// 				&car.Name,
// 				&car.Year,
// 				&car.Brand,
// 				&car.FuelType,
// 				&car.Engine.EngineID,
// 				&car.Price,
// 				&car.CreatedAt,
// 				&car.UpdatedAt,
// 				&car.Engine.EngineID,
// 				&car.Engine.Displacement,
// 				&car.Engine.NoOfCylinders,
// 				&car.Engine.CarRange,
// 			)
// 		} else {
// 			err = rows.Scan(
// 				&car.ID,
// 				&car.Name,
// 				&car.Year,
// 				&car.Brand,
// 				&car.FuelType,
// 				&car.Engine.EngineID,
// 				&car.Price,
// 				&car.CreatedAt,
// 				&car.UpdatedAt,
// 			)
// 		}
// 		if err != nil {
// 			return cars, err
// 		}
// 		cars = append(cars, car)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return cars, err
// 	}

// 	return cars, nil

// }

func (s *CarStore) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "GetCarByBrand")
	defer span.End()
	var cars []models.Car

	query := s.db.WithContext(ctx).Where("brand = ?", brand)

	if isEngine {
		query = query.Preload("Engine")
	}

	if err := query.Find(&cars).Error; err != nil {
		return cars, err
	}

	return cars, nil

}

// func (s *CarStore) CreateCar(ctx context.Context, carRequest *models.CarRequest) (*models.Car, error) {
// 	var createdCar models.Car
// 	var engineId uuid.UUID

// 	err := s.db.QueryRowContext(ctx, "SELECT id FROM engine WHERE id=$1", carRequest.Engine.EngineID).Scan(&engineId)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return &createdCar, errors.New("engine not found")
// 		}
// 		return &createdCar, err
// 	}

// 	carId := uuid.New()
// 	createdAt := time.Now()
// 	updatedAt := createdAt

// 	newCar := models.Car{
// 		ID:        carId,
// 		Brand:     carRequest.Brand,
// 		Name:      carRequest.Name,
// 		Year:      carRequest.Year,
// 		Price:     carRequest.Price,
// 		FuelType:  carRequest.FuelType,
// 		CreatedAt: createdAt,
// 		UpdatedAt: updatedAt,
// 		Engine:    carRequest.Engine,
// 	}

// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return &createdCar, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 			return
// 		}
// 		tx.Commit()
// 	}()

// 	query := `INSERT INTO car (id, name, year, brand, fuel_type, price, engine_id, created_at, updated_at)
// 	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
// 	 RETURNING id, name, year, brand, fuel_type, price, engine_id, created_at, updated_at`

// 	err = tx.QueryRowContext(ctx, query,
// 		newCar.ID,
// 		newCar.Name,
// 		newCar.Year,
// 		newCar.Brand,
// 		newCar.FuelType,
// 		newCar.Price,
// 		newCar.Engine.EngineID,
// 		newCar.CreatedAt,
// 		newCar.UpdatedAt,
// 	).Scan(
// 		&createdCar.ID,
// 		&createdCar.Name,
// 		&createdCar.Year,
// 		&createdCar.Brand,
// 		&createdCar.FuelType,
// 		&createdCar.Price,
// 		&createdCar.Engine.EngineID,
// 		&createdCar.CreatedAt,
// 		&createdCar.UpdatedAt,
// 	)
// 	if err != nil {
// 		return &createdCar, err
// 	}

// 	return &createdCar, nil

// }

func (s *CarStore) CreateCar(ctx context.Context, carRequest *models.CarRequest) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "CreateCar")
	defer span.End()

	var engine models.Engine
	if err := s.db.WithContext(ctx).First(&engine, "engine_id = ?", carRequest.EngineID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("engine not found")
		}
		return nil, err
	}

	car := models.Car{
		ID:       uuid.New(),
		Name:     carRequest.Name,
		Year:     carRequest.Year,
		Brand:    carRequest.Brand,
		FuelType: carRequest.FuelType,
		EngineID: engine.EngineID,
		Price:    carRequest.Price,
	}

	if err := s.db.WithContext(ctx).Create(&car).Error; err != nil {
		return nil, err
	}

	var createdCar models.Car

	if err := s.db.WithContext(ctx).Preload("Engine").First(&createdCar, "id = ?", car.ID).Error; err != nil {
		return nil, err
	}

	return &createdCar, nil

}

// func (s *CarStore) UpdateCar(ctx context.Context, id string, updateCarRequest *models.CarRequest) (*models.Car, error) {

// 	var updatedCar models.Car

// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return &updatedCar, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 			return
// 		}
// 		tx.Commit()
// 	}()

// 	query := `
// 		UPDATE car SET name=$1, year=$2, brand=$3, fuel_type=$4, price=$5, engine_id=$6, updated_at=$7
// 		WHERE id=$8
// 		RETURNING id, name, year, brand, fuel_type, price, engine_id, created_at, updated_at`

// 	err = tx.QueryRowContext(ctx, query,
// 		updateCarRequest.Name,
// 		updateCarRequest.Year,
// 		updateCarRequest.Brand,
// 		updateCarRequest.FuelType,
// 		updateCarRequest.Price,
// 		updateCarRequest.Engine.EngineID,
// 		time.Now(),
// 		id,
// 	).Scan(
// 		&updatedCar.ID,
// 		&updatedCar.Name,
// 		&updatedCar.Year,
// 		&updatedCar.Brand,
// 		&updatedCar.FuelType,
// 		&updatedCar.Price,
// 		&updatedCar.Engine.EngineID,
// 		&updatedCar.CreatedAt,
// 		&updatedCar.UpdatedAt,
// 	)

// 	if err != nil {
// 		return &updatedCar, err
// 	}

// 	return &updatedCar, nil

// }

func (s *CarStore) UpdateCar(ctx context.Context, id string, updateCarRequest *models.CarRequest) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "UpdateCar")
	defer span.End()
	var car models.Car
	if err := s.db.WithContext(ctx).First(&car, "id = ?", id).Error; err != nil {
		return nil, err
	}

	car.Name = updateCarRequest.Name
	car.Year = updateCarRequest.Year
	car.Brand = updateCarRequest.Brand
	car.FuelType = updateCarRequest.FuelType
	car.Price = updateCarRequest.Price
	engineID, err := uuid.Parse(updateCarRequest.EngineID)
	if err != nil {
		return nil, err
	}
	car.EngineID = engineID

	if err := s.db.WithContext(ctx).Save(&car).Error; err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Preload("Engine").First(&car, "id = ?", car.ID).Error; err != nil {
		return nil, err
	}

	return &car, nil

}

// func (s *CarStore) DeleteCar(ctx context.Context, id string) (*models.Car, error) {

// 	var deletedCar models.Car

// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return &deletedCar, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 			return
// 		}
// 		tx.Commit()
// 	}()

// 	err = tx.QueryRowContext(ctx, "SELECT id, name, year, brand, fuel_type, price, engine_id, created_at, updated_at FROM car WHERE id=$1", id).Scan(
// 		&deletedCar.ID,
// 		&deletedCar.Name,
// 		&deletedCar.Year,
// 		&deletedCar.Brand,
// 		&deletedCar.FuelType,
// 		&deletedCar.Price,
// 		&deletedCar.Engine.EngineID,
// 		&deletedCar.CreatedAt,
// 		&deletedCar.UpdatedAt,
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return &deletedCar, errors.New("car not found")
// 		}
// 		return &deletedCar, err
// 	}

// 	result, err := tx.ExecContext(ctx, "DELETE FROM car WHERE id=$1", id)
// 	if err != nil {
// 		return &deletedCar, err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return &deletedCar, err
// 	}
// 	if rowsAffected == 0 {
// 		return &deletedCar, errors.New("car not found")
// 	}
// 	return &deletedCar, nil
// }

func (s *CarStore) DeleteCar(ctx context.Context, id string) (*models.Car, error) {
	ctx, span := otel.Tracer("carservice").Start(ctx, "DeleteCar")
	defer span.End()
	var car models.Car
	if err := s.db.WithContext(ctx).Preload("Engine").First(&car, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("car not found")
		}
		return nil, err
	}

	if err := s.db.WithContext(ctx).Delete(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}
