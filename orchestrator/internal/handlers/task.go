package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"orchestrator/internal/constValues"
	"orchestrator/internal/handlers/models"
	"reflect"
	"strconv"
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
	taskIds, err := a.Tasks.Keys(c.Context(), "*").Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Error{
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
	}

	for _, taskId := range taskIds {
		taskStr, err := a.Tasks.Get(c.Context(), taskId).Result()
		if err != nil {
			continue
		}

		var task models.InternalTask
		if err := json.Unmarshal([]byte(taskStr), &task); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				&fiber.Error{
					Message: err.Error(),
					Code:    fiber.StatusInternalServerError,
				})
		}

		if task.Result == "" {
			if reflect.TypeOf(task.Arg1).String() == "string" {
				arg1Str, err := a.Tasks.Get(c.Context(), task.Arg1.(string)).Result()
				if err != nil {
					continue
				}
				var arg1task models.InternalTask
				if err := json.Unmarshal([]byte(arg1Str), &arg1task); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(
						&fiber.Error{
							Message: err.Error(),
							Code:    fiber.StatusInternalServerError,
						})
				}
				if arg1task.Result == constValues.Processing || arg1task.Result == "" {
					continue
				}

				if arg1task.Result == constValues.Error {
					task.Result = constValues.Error
				}

				var arg1 float64
				if reflect.TypeOf(arg1task.Result).String() == "string" {
					arg1, err = strconv.ParseFloat(arg1task.Result.(string), 64)
					if err != nil {
						continue
					}
				} else {
					arg1 = arg1task.Result.(float64)
				}
				task.Arg1 = arg1
			}
			if reflect.TypeOf(task.Arg2).String() == "string" {
				arg2str, err := a.Tasks.Get(c.Context(), task.Arg2.(string)).Result()
				if err != nil {
					continue
				}
				var arg2task models.InternalTask
				if err := json.Unmarshal([]byte(arg2str), &arg2task); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(
						&fiber.Error{
							Message: err.Error(),
							Code:    fiber.StatusInternalServerError,
						})
				}
				if arg2task.Result == constValues.Processing || arg2task.Result == "" {
					continue
				}

				if arg2task.Result == constValues.Error {
					task.Result = constValues.Error
				}

				var arg2 float64
				if reflect.TypeOf(arg2task.Result).String() == "string" {
					arg2, err = strconv.ParseFloat(arg2task.Result.(string), 64)
					if err != nil {
						continue
					}
				} else {
					arg2 = arg2task.Result.(float64)
				}
				task.Arg2 = arg2
			}

			if task.Result != constValues.Error {
				task.Result = constValues.Processing
			} else {
				if value, err := a.Results.Get(c.Context(), task.ID).Result(); err == nil {
					var taskErr models.InternalTask
					if err := json.Unmarshal([]byte(value), &taskErr); err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(
							&fiber.Error{
								Message: err.Error(),
								Code:    fiber.StatusInternalServerError,
							})
					}
					taskErr.Result = constValues.Error
					if err := a.Results.Set(c.Context(), task.ID, value, 0).Err(); err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(
							&fiber.Error{
								Message: err.Error(),
								Code:    fiber.StatusInternalServerError,
							})
					}
				}
			}

			taskBytes, err := json.Marshal(&task)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Error{
						Message: err.Error(),
						Code:    fiber.StatusInternalServerError,
					})
			}
			if err := a.Tasks.Set(c.Context(), taskId, string(taskBytes), 0).Err(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Error{
						Message: err.Error(),
						Code:    fiber.StatusInternalServerError,
					})
			}

			if task.Result == constValues.Error {
				if _, err := a.Results.Get(c.Context(), task.ID).Result(); err == nil {
					if err := a.Results.Set(c.Context(), task.ID, task.Result, 0).Err(); err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(
							&fiber.Error{
								Message: err.Error(),
								Code:    fiber.StatusInternalServerError,
							})
					}
				}
				continue
			}

			return c.Status(fiber.StatusOK).JSON(
				&models.TaskResponse{
					ID:            task.ID,
					Arg1:          task.Arg1.(float64),
					Arg2:          task.Arg2.(float64),
					Operation:     task.Operation,
					OperationTime: a.cfg.GetOperationTime(task.Operation),
				})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(
		&fiber.Error{
			Message: "Task not found",
			Code:    fiber.StatusNotFound,
		})
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
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			&fiber.Error{
				Message: constValues.ContentTypeError.Error(),
				Code:    fiber.StatusUnprocessableEntity,
			},
		)
	}

	var body models.TaskRequest
	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			&fiber.Error{
				Message: constValues.InvalidJsonError.Error(),
				Code:    fiber.StatusUnprocessableEntity,
			})
	}

	taskStr, err := a.Tasks.Get(c.Context(), body.ID).Result()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			&fiber.Error{
				Message: err.Error(),
				Code:    fiber.StatusNotFound,
			})
	}

	var task models.InternalTask
	if err := json.Unmarshal([]byte(taskStr), &task); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Error{
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
	}

	task.Result = body.Result
	marshal, err := json.Marshal(&task)
	if err != nil {
		return err
	}

	if _, err := a.Results.Get(c.Context(), body.ID).Result(); err == nil {
		if err := a.Results.Set(c.Context(), body.ID, task.Result, 0).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				&fiber.Error{
					Message: err.Error(),
					Code:    fiber.StatusInternalServerError,
				})
		}
	}

	if err := a.Tasks.Set(c.Context(), body.ID, string(marshal), 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Error{
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			})
	}

	return c.Status(fiber.StatusOK).JSON(
		&fiber.Error{
			Message: "ok",
			Code:    fiber.StatusOK,
		})
}
