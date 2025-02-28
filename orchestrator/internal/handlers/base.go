package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Controller struct {
	Expressions *redis.Client
	Results     *redis.Client
	Tasks       *redis.Client
	Validator   *validator.Validate
}
