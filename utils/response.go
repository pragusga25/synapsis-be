package utils

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Success bool `json:"success"`

	Error string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Success bool `json:"success"`

	Data interface{} `json:"data,omitempty"`
}

func ErrorResponseJSON(ctx *fiber.Ctx, status int, err string) error {
	return ctx.Status(status).JSON(ErrorResponse{
		Success: false,
		Error:   err,
	})
}

func SuccessResponseJSON(ctx *fiber.Ctx, status int, data interface{}) error {
	return ctx.Status(status).JSON(SuccessResponse{
		Success: true,
		Data:    data,
	})
}
