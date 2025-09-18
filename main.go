package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/Tushar456/go-carzone/docs"
	"github.com/Tushar456/go-carzone/driver"
	carHandler "github.com/Tushar456/go-carzone/handler/car"
	engineHandler "github.com/Tushar456/go-carzone/handler/engine"
	loginHanler "github.com/Tushar456/go-carzone/handler/login"
	"github.com/Tushar456/go-carzone/middleware"
	"github.com/Tushar456/go-carzone/models"
	carRepository "github.com/Tushar456/go-carzone/repository/car-repository"
	engineRepository "github.com/Tushar456/go-carzone/repository/engine-repository"
	"github.com/Tushar456/go-carzone/service/carService"
	"github.com/Tushar456/go-carzone/service/engineService"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.15.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// @title			Carzone API
// @version		1.0
// @description	API for managing cars and engines.
// @host			localhost:8080
// @BasePath		/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	err := godotenv.Load() // Load environment variables from .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	traceProvider, err := startTracing()
	if err != nil {
		log.Fatalf("Error starting tracing: %v", err)
	}

	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()

	otel.SetTracerProvider(traceProvider)

	db, err := driver.InitDB()

	if err != nil {
		log.Fatalf("Error initializing DB: %v", err)
	}

	fmt.Println("Migrating database...")
	err = db.AutoMigrate(&models.Engine{})
	if err != nil {
		log.Fatalf("Error migrating engine table: %v", err)
	}
	err = db.AutoMigrate(&models.Car{})
	if err != nil {
		log.Fatalf("Error migrating car table: %v", err)
	}
	fmt.Println("Migration successful!")

	// schemaFile := "store/schema.sql"
	// if err = executeSchemaFile(db, schemaFile); err != nil {
	// 	log.Fatalf("Error executing schema file: %v", err)
	// }

	carRepository := carRepository.NewCarRepository(db)
	carService := carService.NewCarService(carRepository)

	engineRepository := engineRepository.NewEngineRepository(db)
	engineService := engineService.NewEngineService(engineRepository)

	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	router := gin.Default()

	router.Use(otelgin.Middleware("carzone"))

	// Middleware to add TraceID to response header
	router.Use(func(c *gin.Context) {
		span := oteltrace.SpanFromContext(c.Request.Context())
		if span.SpanContext().HasTraceID() {
			c.Header("X-Trace-ID", span.SpanContext().TraceID().String())
		}
		c.Next()
	})

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// router := mux.NewRouter()

	// router.HandleFunc("/cars/{id}", carHandler.GetCarByIdHandler).Methods("GET")
	// router.HandleFunc("/cars/brand/{brand}", carHandler.GetCarByBrandHandler).Methods("GET")
	// router.HandleFunc("/cars", carHandler.CreateCarHandler).Methods("POST")
	// router.HandleFunc("/cars/{id}", carHandler.UpdateCarHandler).Methods("PUT")
	// router.HandleFunc("/cars/{id}", carHandler.DeleteCarHandler).Methods("DELETE")

	// router.HandleFunc("/engines/{id}", engineHandler.GetEngineByIdHandler).Methods("GET")
	// router.HandleFunc("/engines", engineHandler.CreateEngineHandler).Methods("POST")
	// router.HandleFunc("/engines/{id}", engineHandler.UpdateEngineHandler).Methods("PUT")
	// router.HandleFunc("/engines/{id}", engineHandler.DeleteEngineHandler).Methods("DELETE")

	router.POST("/login", func(c *gin.Context) {
		loginHanler.LoginHandler(c)
	})

	carRouter := router.Group("/cars").Use(middleware.AuthMiddleware())

	carRouter.GET("/:id", func(c *gin.Context) {
		carHandler.GetCarByIdHandler(c)
	})
	carRouter.GET("/brand/:brand", func(c *gin.Context) {
		carHandler.GetCarByBrandHandler(c)
	})
	carRouter.POST("", func(c *gin.Context) {
		carHandler.CreateCarHandler(c)
	})
	carRouter.PUT("/:id", func(c *gin.Context) {
		carHandler.UpdateCarHandler(c)
	})
	carRouter.DELETE("/:id", func(c *gin.Context) {
		carHandler.DeleteCarHandler(c)
	})

	engineRouter := router.Group("/engines").Use(middleware.AuthMiddleware())

	engineRouter.GET("/:id", func(c *gin.Context) {
		engineHandler.GetEngineByIdHandler(c)
	})
	engineRouter.POST("", func(c *gin.Context) {
		engineHandler.CreateEngineHandler(c)
	})
	engineRouter.PUT("/:id", func(c *gin.Context) {
		engineHandler.UpdateEngineHandler(c)
	})
	engineRouter.DELETE("/:id", func(c *gin.Context) {
		engineHandler.DeleteEngineHandler(c)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(router.Run(":" + port))

}

func executeSchemaFile(db *sql.DB, schemaFile string) error {
	sqlFile, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlFile))

	if err != nil {
		return err
	}

	return nil

}

func startTracing() (*trace.TracerProvider, error) {

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	exporter, err := otlptrace.New(context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("jaeger:4318"),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)

	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
		),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("carzone"),
		),
		),
	)

	return traceProvider, nil

}
