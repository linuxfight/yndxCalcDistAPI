package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"log"
)

func NewRecovery() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				_ = c.Status(fiber.StatusInternalServerError).JSON(
					&fiber.Error{
						Message: err.(error).Error(),
						Code:    fiber.StatusInternalServerError},
				)
			}
		}()

		return c.Next()
	}
}
