package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"orchestrator/internal/handlers"
	"orchestrator/internal/handlers/middlewares"
	"orchestrator/internal/logger"
	"os"

	monitorWare "github.com/gofiber/contrib/monitor"
	corsWare "github.com/gofiber/fiber/v3/middleware/cors"
	healthWare "github.com/gofiber/fiber/v3/middleware/healthcheck"
	loggerWare "github.com/gofiber/fiber/v3/middleware/logger"

	_ "orchestrator/docs"
)

// @title           Orchestrator API
// @version         1.0
// @description 	API documentation for the Calc Orchestrator

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9090
func main() {
	a := fiber.New()

	/* TODO: add here
	TIME_ADDITION_MS
	TIME_SUBTRACTION_MS
	TIME_MULTIPLICATIONS_MS
	TIME_MULTIPLICATIONS_MS
	*/

	debug := os.Getenv("DEBUG") == "TRUE"
	timezone := os.Getenv("TIMEZONE")

	logger.New(debug, timezone)

	logger.Log.Info("Initializing validator...")
	newValidator := validator.New()

	_ = newValidator.RegisterValidation("expression", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		return len(field) >= 1
	})

	logger.Log.Info("Validator initialized")

	logger.Log.Info("Initializing redis")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisExpressions := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})

	if err := redisExpressions.Ping(context.Background()).Err(); err != nil {
		logger.Log.Fatal("Error connecting to redis expressions")
	}

	redisResults := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   1,
	})

	if err := redisResults.Ping(context.Background()).Err(); err != nil {
		logger.Log.Fatal("Error connecting to redis results")
	}

	redisTasks := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   2,
	})

	if err := redisTasks.Ping(context.Background()).Err(); err != nil {
		logger.Log.Fatal("Error connecting to redis tasks")
	}

	logger.Log.Info("Redis initialized")

	// use cors
	a.Use(corsWare.New())
	// recover from panic
	a.Use(middlewares.NewRecovery())
	// logger for requests
	a.Use(loggerWare.New())
	// monitor requests with diagram
	a.Get("/stats", monitorWare.New())
	// healthcheck for initialization
	a.Get(healthWare.DefaultStartupEndpoint, healthWare.NewHealthChecker())
	// swagger web ui
	a.Use(middlewares.NewSwagger(middlewares.SwaggerConfig{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}))

	// create api controller
	h := handlers.Controller{
		Expressions: redisExpressions,
		Results:     redisResults,
		Tasks:       redisTasks,
		Validator:   newValidator,
	}

	// map api routes
	a.Post("/api/v1/calculate", h.PostExpression)
	a.Get("/api/v1/expressions", h.ListExpressions)
	a.Get("/api/v1/expressions/:id", h.GetById)
	a.Get("/internal/task", h.GetTask)
	a.Post("/internal/task", h.SetTask)

	// start server
	err := a.Listen(":9090")
	if err != nil {
		logger.Log.Fatal(err)
	}
}
