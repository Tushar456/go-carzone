package store

import (
	"context"
	"errors"

	"github.com/Tushar456/go-carzone/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

// type EngineStore struct {
// 	db *sql.DB
// }

// func NewEngineStore(db *sql.DB) *EngineStore {
// 	return &EngineStore{
// 		db: db,
// 	}
// }

type EngineStore struct {
	db *gorm.DB
}

func NewEngineStore(db *gorm.DB) *EngineStore {
	return &EngineStore{
		db: db,
	}
}

func (s *EngineStore) GetEngineById(ctx context.Context, id string) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "GetEngineById")
	defer span.End()
	var engine models.Engine
	if err := s.db.WithContext(ctx).First(&engine, "engine_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Engine{}, nil
		}
		return &models.Engine{}, err
	}
	return &engine, nil

}

// func (s *EngineStore) GetEngineById(ctx context.Context, id string) (*models.Engine, error) {

// 	var engine models.Engine

// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return &engine, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			if rnErr := tx.Rollback(); rnErr != nil {
// 				fmt.Printf("failed to rollback transaction: %v\n", rnErr)
// 			}
// 		}
// 		if cmErr := tx.Commit(); cmErr != nil {
// 			fmt.Printf("failed to commit transaction: %v\n", cmErr)
// 		}
// 	}()

// 	err = tx.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id=$1", id).Scan(
// 		&engine.EngineID,
// 		&engine.Displacement,
// 		&engine.NoOfCylinders,
// 		&engine.CarRange,
// 	)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return &engine, nil

// 		}
// 		return &engine, err
// 	}

// 	return &engine, nil

// }

func (s *EngineStore) CreateEngine(ctx context.Context, engineRequest *models.EngineRequest) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "CreateEngine")
	defer span.End()
	engine := models.Engine{
		EngineID:      uuid.New(),
		Displacement:  engineRequest.Displacement,
		NoOfCylinders: engineRequest.NoOfCylinders,
		CarRange:      engineRequest.CarRange,
	}

	if err := s.db.WithContext(ctx).Create(&engine).Error; err != nil {
		return &models.Engine{}, err
	}

	return &engine, nil
}

// func (s *EngineStore) CreateEngine(ctx context.Context, engineRequest *models.EngineRequest) (*models.Engine, error) {

// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return &models.Engine{}, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			if rnErr := tx.Rollback(); rnErr != nil {
// 				fmt.Printf("failed to rollback transaction: %v\n", rnErr)
// 			}
// 		}
// 		if cmErr := tx.Commit(); cmErr != nil {
// 			fmt.Printf("failed to commit transaction: %v\n", cmErr)
// 		}
// 	}()

// 	engineID := uuid.New()
// 	_, err = tx.ExecContext(ctx, "INSERT INTO engine (id, displacement, no_of_cylinders, car_range) VALUES ($1, $2, $3, $4)",
// 		engineID, engineRequest.Displacement, engineRequest.NoOfCylinders, engineRequest.CarRange)

// 	if err != nil {
// 		return &models.Engine{}, err
// 	}

// 	createdEngine := models.Engine{
// 		EngineID:      engineID,
// 		Displacement:  engineRequest.Displacement,
// 		NoOfCylinders: engineRequest.NoOfCylinders,
// 		CarRange:      engineRequest.CarRange,
// 	}

// 	return &createdEngine, nil
// }

// func (s *EngineStore) UpdateEngine(ctx context.Context, id string, engineRequest *models.EngineRequest) (*models.Engine, error) {

// 	engineId, err := uuid.Parse(id)
// 	if err != nil {
// 		return &models.Engine{}, err
// 	}

// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return &models.Engine{}, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			if rnErr := tx.Rollback(); rnErr != nil {
// 				fmt.Printf("failed to rollback transaction: %v\n", rnErr)
// 			}
// 		}
// 		if cmErr := tx.Commit(); cmErr != nil {
// 			fmt.Printf("failed to commit transaction: %v\n", cmErr)
// 		}
// 	}()

// 	result, err := tx.ExecContext(ctx, "UPDATE engine SET displacement=$1, no_of_cylinders=$2, car_range=$3 WHERE id=$4",
// 		engineRequest.Displacement, engineRequest.NoOfCylinders, engineRequest.CarRange, engineId)

// 	if err != nil {
// 		return &models.Engine{}, err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return &models.Engine{}, err
// 	}
// 	if rowsAffected == 0 {
// 		return &models.Engine{}, errors.New("engine not found")
// 	}

// 	updatedEngine := models.Engine{
// 		EngineID:      engineId,
// 		Displacement:  engineRequest.Displacement,
// 		NoOfCylinders: engineRequest.NoOfCylinders,
// 		CarRange:      engineRequest.CarRange,
// 	}

// 	return &updatedEngine, nil

// }

func (s *EngineStore) UpdateEngine(ctx context.Context, id string, engineRequest *models.EngineRequest) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "UpdateEngine")
	defer span.End()
	engineId, err := uuid.Parse(id)
	if err != nil {
		return &models.Engine{}, err
	}

	var engine models.Engine
	if err := s.db.WithContext(ctx).First(&engine, "engine_id = ?", engineId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Engine{}, errors.New("engine not found")
		}
		return &models.Engine{}, err
	}

	engine.Displacement = engineRequest.Displacement
	engine.NoOfCylinders = engineRequest.NoOfCylinders
	engine.CarRange = engineRequest.CarRange

	if err := s.db.WithContext(ctx).Save(&engine).Error; err != nil {
		return &models.Engine{}, err
	}

	return &engine, nil
}

// func (s *EngineStore) DeleteEngine(ctx context.Context, id string) (*models.Engine, error) {
// 	var engine models.Engine

// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return &models.Engine{}, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			if rnErr := tx.Rollback(); rnErr != nil {
// 				fmt.Printf("failed to rollback transaction: %v\n", rnErr)
// 			}
// 		}
// 		if cmErr := tx.Commit(); cmErr != nil {
// 			fmt.Printf("failed to commit transaction: %v\n", cmErr)
// 		}
// 	}()

// 	err = tx.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id=$1", id).Scan(
// 		&engine.EngineID,
// 		&engine.Displacement,
// 		&engine.NoOfCylinders,
// 		&engine.CarRange,
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return &models.Engine{}, nil
// 		}
// 		return &models.Engine{}, err
// 	}

// 	result, err := tx.ExecContext(ctx, "DELETE FROM engine WHERE id=$1", id)
// 	if err != nil {
// 		return &models.Engine{}, err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return &models.Engine{}, err
// 	}
// 	if rowsAffected == 0 {
// 		return &models.Engine{}, errors.New("engine not found")
// 	}

// 	return &engine, nil

// }

func (s *EngineStore) DeleteEngine(ctx context.Context, id string) (*models.Engine, error) {
	ctx, span := otel.Tracer("engineservice").Start(ctx, "DeleteEngine")
	defer span.End()
	var engine models.Engine

	if err := s.db.WithContext(ctx).First(&engine, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Engine{}, errors.New("engine not found")
		}
		return &models.Engine{}, err
	}

	if err := s.db.WithContext(ctx).Delete(&engine).Error; err != nil {
		return &models.Engine{}, err
	}

	return &engine, nil

}
