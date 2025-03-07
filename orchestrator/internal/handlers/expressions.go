package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"orchestrator/internal/constValues"
	"orchestrator/internal/handlers/models"
	"strconv"
)

// ListExpressions @Summary      Получить весь список выражений
// @Tags         expressions
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ListAllExpressionsResponse
// @Failure      500  {object}  models.ApiError
// @Router       /api/v1/expressions [get]
func (a *Controller) ListExpressions(c fiber.Ctx) error {
	result, err := a.Results.Keys(c.Context(), "*").Result()
	expressions := []models.Expression{}
	if err != nil && !errors.Is(err, redis.Nil) {
		return sendError(c, fiber.StatusInternalServerError, err)
	} else if errors.Is(err, redis.Nil) {
		return c.Status(fiber.StatusOK).JSON(models.ListAllExpressionsResponse{Expressions: expressions})
	}

	for _, id := range result {
		value, err := a.Results.Get(c.Context(), id).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return sendError(c, fiber.StatusInternalServerError, err)
		}
		switch value {
		case constValues.Error:
			expressions = append(expressions, models.Expression{
				Id:     id,
				Result: 0,
				Status: constValues.Error,
			})
		case constValues.Processing:
			expressions = append(expressions, models.Expression{
				Id:     id,
				Result: 0,
				Status: value,
			})
		default:
			r, _ := strconv.ParseFloat(value, 64)
			expressions = append(expressions, models.Expression{
				Id:     id,
				Result: r,
				Status: constValues.Done,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(&models.ListAllExpressionsResponse{Expressions: expressions})
}

// GetById @Summary      Получить выражение по UUID
// @Tags         expressions
// @Accept       json
// @Produce      json
// @Param        id path  string true  "UUID выражения"
// @Success      200  {object}  models.GetByIdExpressionResponse
// @Failure      404  {object}  models.ApiError
// @Failure      422  {object}  models.ApiError
// @Failure      500  {object}  models.ApiError
// @Router       /api/v1/expressions/{id} [get]
func (a *Controller) GetById(c fiber.Ctx) error {
	id := c.Params("id")
	if uuid.Validate(id) != nil {
		return sendError(c, fiber.StatusUnprocessableEntity, constValues.InvalidUuidError)
	}

	value, err := a.Results.Get(c.Context(), id).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return sendError(c, fiber.StatusInternalServerError, err)
	} else if errors.Is(err, redis.Nil) {
		return sendError(c, fiber.StatusNotFound, constValues.NotFoundError)
	}

	expression := models.Expression{}
	switch value {
	case constValues.Error:
		expression = models.Expression{
			Id:     id,
			Result: 0,
			Status: constValues.Error,
		}
	case constValues.Processing:
		expression = models.Expression{
			Id:     id,
			Result: 0,
			Status: value,
		}
	default:
		r, _ := strconv.ParseFloat(value, 64)
		expression = models.Expression{
			Id:     id,
			Result: r,
			Status: constValues.Done,
		}
	}

	return c.Status(fiber.StatusOK).JSON(&models.GetByIdExpressionResponse{Expression: expression})
}
