package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tushar456/go-carzone/models"
	"github.com/Tushar456/go-carzone/service"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

type CarHandler struct {
	carService service.CarServiceInterface
}

func NewCarHandler(carService service.CarServiceInterface) *CarHandler {
	return &CarHandler{
		carService: carService,
	}
}

// GetCarByIdHandler godoc
//
//	@Summary		Get car by ID
//	@Description	get car by ID
//	@Tags			cars
//	@Param			id	path		string	true	"Car ID"
//	@Success		200	{object}	models.Car
//	@Failure		404	{object}	map[string]string
//	@Router			/cars/{id} [get]
//
// @Security     BearerAuth
func (ch *CarHandler) GetCarByIdHandler(c *gin.Context) {
	ctx, span := otel.Tracer("carservice").Start(c.Request.Context(), "GetCarByIdHandler")
	defer span.End()
	id := c.Param("id")
	car, err := ch.carService.GetCarById(ctx, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("Error fetching car by ID: %v", err)
		return
	}

	if car.ID.String() == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "not found"})
		return
	}

	body, err := json.Marshal(car)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error marshalling car data: %v", err)
		return
	}
	c.Data(http.StatusOK, "application/json", body)

}

// GetCarByBrandHandler godoc
//
//	@Summary		Get cars by brand
//	@Description	get cars by brand
//	@Tags			cars
//	@Param			brand		path		string	true	"Brand"
//	@Param			isEngine	query		bool	false	"Include engine"
//	@Success		200			{array}		models.Car
//	@Failure		404			{object}	map[string]string
//	@Router			/cars/brand/{brand} [get]
//
// @Security     BearerAuth
func (ch *CarHandler) GetCarByBrandHandler(c *gin.Context) {
	ctx, span := otel.Tracer("carservice").Start(c.Request.Context(), "GetCarByBrandHandler")
	defer span.End()

	brand := c.Param("brand")
	isEngine := c.DefaultQuery("isEngine", "false") == "true"
	cars, err := ch.carService.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("Error fetching cars by brand: %v", err)
		return
	}

	if len(cars) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No cars found"})
		return
	}

	c.JSON(http.StatusOK, cars)
}

// CreateCarHandler godoc
//
//	@Summary		Create car
//	@Description	create car
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			car	body		models.CarRequest	true	"Car Request"
//	@Success		201	{object}	models.Car
//	@Failure		400	{object}	map[string]string
//	@Router			/cars [post]
//
// @Security BearerAuth
func (ch *CarHandler) CreateCarHandler(c *gin.Context) {
	ctx, span := otel.Tracer("carservice").Start(c.Request.Context(), "CreateCarHandler")
	defer span.End()

	var carRequest models.CarRequest
	err := c.ShouldBindJSON(&carRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Error decoding car request: %v", err)
		return
	}

	createdCar, err := ch.carService.CreateCar(ctx, &carRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("Error creating car: %v", err)
		return
	}
	body, err := json.Marshal(createdCar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error marshalling created car data: %v", err)
		return
	}
	c.Data(http.StatusCreated, "application/json", body)

}

// UpdateCarHandler godoc
//
//	@Summary		Update car
//	@Description	update car
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Car ID"
//	@Param			car	body		models.CarRequest	true	"Car Request"
//	@Success		200	{object}	models.Car
//	@Failure		400	{object}	map[string]string
//	@Router			/cars/{id} [put]
//
// @Security BearerAuth
func (ch *CarHandler) UpdateCarHandler(c *gin.Context) {
	ctx, span := otel.Tracer("carservice").Start(c.Request.Context(), "UpdateCarHandler")
	defer span.End()

	var carRequest models.CarRequest
	id := c.Param("id")
	err := c.ShouldBindJSON(&carRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Error decoding car request: %v", err)
		return
	}

	updatedCar, err := ch.carService.UpdateCar(ctx, id, &carRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("Error updating car: %v", err)
		return
	}
	body, err := json.Marshal(updatedCar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error marshalling updated car data: %v", err)
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}

// DeleteCarHandler godoc
//
//	@Summary		Delete car
//	@Description	delete car
//	@Tags			cars
//	@Param			id	path		string	true	"Car ID"
//	@Success		200	{object}	models.Car
//	@Failure		404	{object}	map[string]string
//	@Router			/cars/{id} [delete]
//
// @Security BearerAuth
func (ch *CarHandler) DeleteCarHandler(c *gin.Context) {
	ctx, span := otel.Tracer("carservice").Start(c.Request.Context(), "DeleteCarHandler")
	defer span.End()
	id := c.Param("id")

	deletedCar, err := ch.carService.DeleteCar(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("Error deleting car: %v", err)
		return
	}
	if deletedCar.ID.String() == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}
	body, err := json.Marshal(deletedCar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error marshalling deleted car data: %v", err)
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}
