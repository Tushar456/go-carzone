package models

import (
	"errors"

	"github.com/google/uuid"
)

// type Engine struct {
// 	EngineID      uuid.UUID `json:"engine_id"`
// 	Displacement  int       `json:"displacement"`
// 	NoOfCylinders int       `json:"no_of_cylinders"`
// 	CarRange      int       `json:"car_range"`
// }

type Engine struct {
	EngineID      uuid.UUID `json:"engine_id" gorm:"type:uuid;primaryKey"`
	Displacement  int       `json:"displacement"`
	NoOfCylinders int       `json:"no_of_cylinders"`
	CarRange      int       `json:"car_range"`
}

type EngineRequest struct {
	Displacement  int `json:"displacement"`
	NoOfCylinders int `json:"no_of_cylinders"`
	CarRange      int `json:"car_range"`
}

func (e *EngineRequest) Validate() error {
	if err := validateDisplacement(e.Displacement); err != nil {
		return err
	}

	if err := validateNoOfCylinders(e.NoOfCylinders); err != nil {
		return err
	}

	if err := validateCarRange(e.CarRange); err != nil {
		return err
	}

	return nil
}

func validateDisplacement(displacement int) error {
	if displacement <= 0 {
		return errors.New("displacement must be greater than 0")
	}
	return nil
}

func validateNoOfCylinders(noOfCylinders int) error {
	if noOfCylinders <= 0 {
		return errors.New("number of cylinders must be greater than 0")
	}
	return nil
}

func validateCarRange(carRange int) error {
	if carRange < 0 {
		return errors.New("car range cannot be negative")
	}
	return nil
}
