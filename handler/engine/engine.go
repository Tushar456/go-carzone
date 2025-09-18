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

type EngineHandler struct {
	engineService service.EngineServiceInterface
}

func NewEngineHandler(engineService service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		engineService: engineService,
	}
}

// GetEngineByIdHandler godoc
// @Summary      Get engine by ID
// @Description  get engine by ID
// @Tags         engines
// @Param        id   path      string  true  "Engine ID"
// @Success      200  {object}  models.Engine
// @Failure      404  {object}  map[string]string
// @Router       /engines/{id} [get]
// @Security     BearerAuth
func (eh *EngineHandler) GetEngineByIdHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	engine, err := eh.engineService.GetEngineById(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("Error fetching engine by ID: %v", err)
		return
	}
	if engine.EngineID.String() == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Engine not found"})
		return
	}
	body, err := json.Marshal(engine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error marshalling engine data: %v", err)
		return
	}
	c.Data(http.StatusOK, "application/json", body)

}

// CreateEngineHandler godoc
// @Summary      Create engine
// @Description  create engine
// @Tags         engines
// @Accept       json
// @Produce      json
// @Param        engine  body      models.EngineRequest  true  "Engine Request"
// @Success      201     {object}  models.Engine
// @Failure      400     {object}  map[string]string
// @Router       /engines [post]
// @Security     BearerAuth
func (eh *EngineHandler) CreateEngineHandler(c *gin.Context) {

	ctx, span := otel.Tracer("engineservice").Start(c.Request.Context(), "CreateEngineHandler")
	defer span.End()

	var engineRequest models.EngineRequest
	err := json.NewDecoder(c.Request.Body).Decode(&engineRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		log.Printf("Error decoding engine request: %v", err)
		return
	}

	createdEngine, err := eh.engineService.CreateEngine(ctx, &engineRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("Error creating engine: %v", err)
		return
	}
	body, err := json.Marshal(createdEngine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error marshalling created engine data: %v", err)
		return
	}
	c.Data(http.StatusCreated, "application/json", body)

}

// UpdateEngineHandler godoc
// @Summary      Update engine
// @Description  update engine
// @Tags         engines
// @Accept       json
// @Produce      json
// @Param        id      path      string              true  "Engine ID"
// @Param        engine  body      models.EngineRequest  true  "Engine Request"
// @Success      200     {object}  models.Engine
// @Failure      400     {object}  map[string]string
// @Security     BearerAuth
func (eh *EngineHandler) UpdateEngineHandler(c *gin.Context) {

	ctx, span := otel.Tracer("engineservice").Start(c.Request.Context(), "UpdateEngineHandler")
	defer span.End()

	var engineRequest models.EngineRequest
	id := c.Param("id")
	err := json.NewDecoder(c.Request.Body).Decode(&engineRequest)
	if err != nil {

		log.Printf("Error decoding engine request: %v", err)
		return
	}

	updatedEngine, err := eh.engineService.UpdateEngine(ctx, id, &engineRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Printf("Error updating engine: %v", err)
		return
	}
	body, err := json.Marshal(updatedEngine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error marshalling updated engine data: %v", err)
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}

// DeleteEngineHandler godoc
// @Summary      Delete engine
// @Description  delete engine
// @Tags         engines
// @Param        id   path      string  true  "Engine ID"
// @Success      200  {object}  models.Engine
// @Failure      404  {object}  map[string]string
// @Router       /engines/{id} [delete]
// @Security     BearerAuth
func (eh *EngineHandler) DeleteEngineHandler(c *gin.Context) {
	ctx, span := otel.Tracer("engineservice").Start(c.Request.Context(), "DeleteEngineHandler")
	defer span.End()
	id := c.Param("id")

	deletedEngine, err := eh.engineService.DeleteEngine(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Engine ID"})
		log.Printf("Error deleting engine: %v", err)
		return
	}
	if deletedEngine.EngineID.String() == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Engine not found"})
		return
	}
	body, err := json.Marshal(deletedEngine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Error marshalling deleted engine data: %v", err)
		return
	}
	c.Data(http.StatusOK, "application/json", body)

}
