package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	monitorWare "github.com/gofiber/contrib/monitor"
	"github.com/gofiber/fiber/v3"
	corsWare "github.com/gofiber/fiber/v3/middleware/cors"
	healthWare "github.com/gofiber/fiber/v3/middleware/healthcheck"
	loggerWare "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/redis/go-redis/v9"
	"orchestrator/internal/constValues"
	"orchestrator/internal/handlers/middlewares"
	"orchestrator/internal/handlers/models"
	"orchestrator/internal/logger"
	"os"
	"strconv"
)

type Controller struct {
	app         *fiber.App
	cfg         *Config
	Expressions *redis.Client
	Results     *redis.Client
	Tasks       *redis.Client
	Validator   *validator.Validate
}

func (a *Controller) Start() {
	err := a.app.Listen(":9090")
	if err != nil {
		logger.Log.Fatal(err)
	}
}

func New() *Controller {
	a := fiber.New()

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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	logger.Log.Info("Initializing redis")

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

	timeAddStr := os.Getenv("TIME_ADDITION_MS")
	timeAdd := 1000
	if timeAddStr != "" {
		var err error
		timeAdd, err = strconv.Atoi(timeAddStr)
		if err != nil {
			logger.Log.Fatal(err)
		}
	}
	timeSubStr := os.Getenv("TIME_SUBTRACTION_MS")
	timeSub := 1000
	if timeSubStr != "" {
		var err error
		timeSub, err = strconv.Atoi(timeSubStr)
		if err != nil {
			logger.Log.Fatal(err)
		}
	}
	timeMulStr := os.Getenv("TIME_MULTIPLICATIONS_MS")
	timeMul := 1000
	if timeMulStr != "" {
		var err error
		timeMul, err = strconv.Atoi(timeMulStr)
		if err != nil {
			logger.Log.Fatal(err)
		}
	}
	timeDivStr := os.Getenv("TIME_DIVISIONS_MS")
	timeDiv := 1000
	if timeDivStr != "" {
		var err error
		timeDiv, err = strconv.Atoi(timeDivStr)
		if err != nil {
			logger.Log.Fatal(err)
		}
	}

	// create api controller
	h := &Controller{
		Expressions: redisExpressions,
		Results:     redisResults,
		Tasks:       redisTasks,
		Validator:   newValidator,
		app:         a,
		cfg: &Config{
			TimeAdditionMS:       timeAdd,
			TimeSubtractionMS:    timeSub,
			TimeMultiplicationMS: timeMul,
			TimeDivisionMS:       timeDiv,
		},
	}

	// map api routes
	a.Post("/api/v1/calculate", h.PostExpression)
	a.Get("/api/v1/expressions", h.ListExpressions)
	a.Get("/api/v1/expressions/:id", h.GetById)
	a.Get("/internal/task", h.GetTask)
	a.Post("/internal/task", h.SetTask)

	return h
}

type Config struct {
	TimeAdditionMS       int
	TimeSubtractionMS    int
	TimeMultiplicationMS int
	TimeDivisionMS       int
}

func (c *Config) GetOperationTime(operation string) int {
	switch operation {
	case "+":
		return c.TimeAdditionMS
	case "-":
		return c.TimeSubtractionMS
	case "*":
		return c.TimeMultiplicationMS
	case "/":
		return c.TimeDivisionMS
	}

	return 0
}

func (a *Controller) getTask(ctx context.Context, taskId string) (models.InternalTask, error) {
	taskStr, err := a.Tasks.Get(ctx, taskId).Result()
	if err != nil {
		return models.InternalTask{}, err
	}

	var task models.InternalTask
	if err := json.Unmarshal([]byte(taskStr), &task); err != nil {
		return models.InternalTask{}, err
	}
	return task, nil
}

func (a *Controller) processTaskArguments(ctx context.Context, taskId string, task *models.InternalTask) (bool, error) {
	hasError := false

	if err := processArgument(ctx, a, task, &task.Arg1); err != nil {
		return false, err
	}
	if task.Result == constValues.Error {
		hasError = true
	}

	if err := processArgument(ctx, a, task, &task.Arg2); err != nil {
		return false, err
	}
	if task.Result == constValues.Error {
		hasError = true
	}

	if hasError {
		return true, a.updateErrorTask(ctx, taskId, task)
	}

	if task.Result == "" {
		task.Result = constValues.Processing
	}
	return false, a.updateTask(ctx, taskId, task)
}

func processArgument(ctx context.Context, a *Controller, task *models.InternalTask, arg *interface{}) error {
	argStr, ok := (*arg).(string)
	if !ok {
		return nil
	}

	argTask, err := a.getTask(ctx, argStr)
	if err != nil || argTask.Result == constValues.Processing || argTask.Result == "" {
		return err
	}

	if argTask.Result == constValues.Error {
		task.Result = constValues.Error
		return nil
	}

	result, err := convertResult(argTask.Result)
	if err != nil {
		return err
	}

	*arg = result
	return nil
}

func convertResult(result interface{}) (float64, error) {
	switch v := result.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("unexpected result type: %T", v)
	}
}

func (a *Controller) updateErrorTask(ctx context.Context, taskId string, task *models.InternalTask) error {
	if err := a.updateTask(ctx, taskId, task); err != nil {
		return err
	}

	if err := a.Results.Get(ctx, task.ID).Err(); errors.Is(err, redis.Nil) {
		return a.Results.Set(ctx, task.ID, task.Result, 0).Err()
	}
	return nil
}

func (a *Controller) updateTask(ctx context.Context, taskId string, task *models.InternalTask) error {
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return a.Tasks.Set(ctx, taskId, string(taskBytes), 0).Err()
}

func (a *Controller) handleValidTask(c fiber.Ctx, task *models.InternalTask) error {
	arg1, ok1 := task.Arg1.(float64)
	arg2, ok2 := task.Arg2.(float64)
	if !ok1 || !ok2 {
		return sendError(c, fiber.StatusInternalServerError, fmt.Errorf("invalid argument types"))
	}

	return c.Status(fiber.StatusOK).JSON(&models.TaskResponse{
		ID:            task.ID,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     task.Operation,
		OperationTime: a.cfg.GetOperationTime(task.Operation),
	})
}

func sendError(c fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(&fiber.Error{
		Message: err.Error(),
		Code:    status,
	})
}
