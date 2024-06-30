package controllers

import (
	"errors"
	"synapsis/models"
	"synapsis/services"
	"synapsis/utils"
	"synapsis/validations"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, err.Error())
	}

	err := c.authService.Register(&user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return utils.ErrorResponseJSON(ctx, fiber.StatusConflict, "Email already registered")
			}
		}

		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, "Failed to register user")
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusCreated, "User registered successfully")
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {

	var req validations.LoginDto

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, err.Error())
	}

	token, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusUnauthorized, "Invalid credentials")
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, fiber.Map{"token": token})
}
