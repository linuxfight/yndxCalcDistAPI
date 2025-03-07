package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"orchestrator/internal/constValues"
	"orchestrator/internal/handlers/models"
)

// GetTask @Summary      Получить выражение на выполнение
// @Tags         internal
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.TaskResponse
// @Failure      404  {object}  models.ApiError
// @Failure      500  {object}  models.ApiError
// @Router       /internal/task [get]
func (a *Controller) GetTask(c fiber.Ctx) error {
	ctx := c.Context()
	taskIds, err := a.Tasks.Keys(ctx, "*").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return sendError(c, fiber.StatusNotFound, constValues.NotFoundError)
		}
		return sendError(c, fiber.StatusInternalServerError, err)
	}

	for _, taskId := range taskIds {
		task, err := a.getTask(ctx, taskId)
		if err != nil || task.Result != "" {
			continue
		}

		hasError, err := a.processTaskArguments(ctx, taskId, task)
		if err != nil {
			return sendError(c, fiber.StatusInternalServerError, err)
		}
		if hasError {
			continue
		}

		if resp := a.getTaskResponse(task); resp != nil {
			task.Result = constValues.Processing
			err := a.updateTask(ctx, taskId, task)
			if err != nil {
				return sendError(c, fiber.StatusInternalServerError, err)
			}
			return c.Status(fiber.StatusOK).JSON(&resp)
		}
		continue
	}

	return sendError(c, fiber.StatusNotFound, constValues.NotFoundError)
}

// SetTask @Summary      Обновить результат выражения
// @Tags         internal
// @Accept       json
// @Produce      json
// @Param        body body  models.TaskRequest true  "Объект, содержащий в себе результат части выражения"
// @Success      200  {object}  models.ApiError
// @Failure      404  {object}  models.ApiError
// @Failure      422  {object}  models.ApiError
// @Failure      500  {object}  models.ApiError
// @Router       /internal/task [post]
func (a *Controller) SetTask(c fiber.Ctx) error {
	if c.Get("Content-Type") != "application/json" {
		return sendError(c, fiber.StatusUnprocessableEntity, constValues.ContentTypeError)
	}

	var body models.TaskRequest
	if err := c.Bind().JSON(&body); err != nil {
		return sendError(c, fiber.StatusUnprocessableEntity, constValues.InvalidJsonError)
	}

	taskStr, err := a.Tasks.Get(c.Context(), body.ID).Result()
	if err != nil {
		return sendError(c, fiber.StatusNotFound, err)
	}

	var task models.InternalTask
	if err := json.Unmarshal([]byte(taskStr), &task); err != nil {
		return sendError(c, fiber.StatusInternalServerError, err)
	}

	task.Result = body.Result
	marshal, err := json.Marshal(&task)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, err)
	}

	if _, err := a.Results.Get(c.Context(), body.ID).Result(); err == nil {
		if err := a.Results.Set(c.Context(), body.ID, task.Result, 0).Err(); err != nil {
			return sendError(c, fiber.StatusInternalServerError, err)
		}
	}

	if err := a.Tasks.Set(c.Context(), body.ID, string(marshal), 0).Err(); err != nil {
		return sendError(c, fiber.StatusInternalServerError, err)
	}

	return c.Status(fiber.StatusOK).JSON(
		&fiber.Error{
			Message: "ok",
			Code:    fiber.StatusOK,
		})
}
