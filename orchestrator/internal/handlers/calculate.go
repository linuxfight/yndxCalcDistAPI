package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"orchestrator/internal/calc"
	"orchestrator/internal/constValues"
	"orchestrator/internal/handlers/models"
	"strings"
)

// PostExpression @Summary      Добавить выражение в очередь на выполнение
// @Tags         calculate
// @Accept       json
// @Produce      json
// @Param        body body  models.CalculateRequest true  "Объект, содержащий в себе выражение"
// @Success      200  {object}  models.CalculateResponse
// @Success      201  {object}  models.CalculateResponse
// @Failure      422  {object}  models.ApiError
// @Failure      500  {object}  models.ApiError
// @Router       /api/v1/calculate [post]
func (a Controller) PostExpression(c fiber.Ctx) error {
	if c.Get("Content-Type") != "application/json" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			&fiber.Error{
				Message: constValues.ContentTypeError.Error(),
				Code:    fiber.StatusUnprocessableEntity,
			},
		)
	}

	var body models.CalculateRequest
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			&fiber.Error{
				Message: constValues.InvalidJsonError.Error(),
				Code:    fiber.StatusUnprocessableEntity,
			})
	}

	body.Expression = strings.ReplaceAll(body.Expression, " ", "")
	body.Expression = strings.ReplaceAll(body.Expression, ",", ".")

	result, err := a.Expressions.Get(c.Context(), body.Expression).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Error{
				Message: constValues.InvalidExpressionError.Error(),
				Code:    fiber.StatusInternalServerError,
			})
	} else if errors.Is(err, redis.Nil) {
		id := uuid.New().String()

		tasks, err := calc.ParseExpression(body.Expression)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(
				&fiber.Error{
					Message: constValues.InvalidExpressionError.Error(),
					Code:    fiber.StatusUnprocessableEntity,
				})
		}

		if a.Expressions.Set(c.Context(), body.Expression, id, 0).Err() != nil ||
			a.Results.Set(c.Context(), id, constValues.Processing, 0).Err() != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				&fiber.Error{
					Message: constValues.InvalidExpressionError.Error(),
					Code:    fiber.StatusInternalServerError,
				})
		}

		for i, task := range tasks {
			var taskString []byte
			if i == len(tasks)-1 {
				task.ID = id
			}
			if taskString, err = json.Marshal(task); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Error{
						Message: constValues.InvalidExpressionError.Error(),
						Code:    fiber.StatusInternalServerError,
					})
			}
			if a.Tasks.Set(c.Context(), task.ID, string(taskString), 0).Err() != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Error{
						Message: constValues.InvalidExpressionError.Error(),
						Code:    fiber.StatusInternalServerError,
					})
			}
		}

		return c.Status(fiber.StatusCreated).JSON(models.CalculateResponse{
			Id: id,
		})
	}

	return c.Status(fiber.StatusOK).JSON(&models.CalculateResponse{Id: result})
}
