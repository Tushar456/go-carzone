package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// type Car struct {
// 	ID        uuid.UUID `json:"id"`
// 	Name      string    `json:"name"`
// 	Year      string    `json:"year"`
// 	Brand     string    `json:"brand"`
// 	FuelType  string    `json:"fuel_type"`
// 	Engine    Engine    `json:"engine"`
// 	Price     float64   `json:"price"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }

// GORM-compatible Car model
type Car struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name"`
	Year      string    `json:"year"`
	Brand     string    `json:"brand"`
	FuelType  string    `json:"fuel_type"`
	EngineID  uuid.UUID `json:"engine_id" gorm:"type:uuid"`
	Engine    Engine    `json:"engine" gorm:"foreignKey:EngineID;references:ID"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CarRequest struct {
	Name     string  `json:"name"`
	Year     string  `json:"year"`
	Brand    string  `json:"brand"`
	FuelType string  `json:"fuel_type"`
	EngineID string  `json:"engine_id"`
	Price    float64 `json:"price"`
}

func (c *CarRequest) Validate() error {

	if err := validateName(c.Name); err != nil {
		return err
	}

	if err := validateYear(c.Year); err != nil {
		return err
	}

	if err := validateBrand(c.Brand); err != nil {
		return err
	}

	if err := validateFuelType(c.FuelType); err != nil {
		return err
	}

	if err := validateEngine(c.EngineID); err != nil {
		return err
	}

	if err := validatePrice(c.Price); err != nil {
		return err
	}

	return nil
}

func validateName(name string) error {

	if name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func validateYear(year string) error {

	if year == "" {
		return errors.New("year cannot be empty")
	}

	yearint, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("year must be a number")
	}

	currentYear := time.Now().Year()
	if yearint < 1886 || yearint > currentYear {
		return errors.New("year must be between 1886 and " + strconv.Itoa(currentYear))
	}
	return nil
}

func validateBrand(brand string) error {

	if brand == "" {
		return errors.New("brand cannot be empty")
	}

	return nil
}
func validateFuelType(fuelType string) error {
	validateFuelTYpes := []string{"Petrol", "Diesel", "Electric", "Hybrid"}

	if fuelType == "" {
		return errors.New("fuel type cannot be empty")
	}

	for _, validTfuelType := range validateFuelTYpes {
		if fuelType == validTfuelType {
			return nil
		}
	}
	return errors.New("fuel type must be one of Petrol, Diesel, Electric, Hybrid")
}

func validatePrice(price float64) error {

	if price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}

func validateEngine(engineId string) error {

	if engineId == "" {
		return errors.New("engine id cannot be empty")
	}

	_, err := uuid.Parse(engineId)
	if err != nil {
		return errors.New("engine id must be a valid UUID")
	}
	return nil
}
