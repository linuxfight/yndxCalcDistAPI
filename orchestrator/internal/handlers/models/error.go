package models

// ApiError needs to exist, because Swag cannot parse fiber.Error
type ApiError struct {
	Message string `json:"message"`
	Code    int    `json:"status"`
}
