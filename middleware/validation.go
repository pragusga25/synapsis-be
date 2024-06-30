package middleware

import (
	"fmt"
	"synapsis/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateBody(s interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(s); err != nil {
			return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, err.Error())
		}

		if err := validate.Struct(s); err != nil {
			valErr := err.(validator.ValidationErrors)[0]
			errMessage := fmt.Sprintf("Error on field %s, condition: %s", valErr.Field(), valErr.ActualTag())

			return utils.ErrorResponseJSON(c, fiber.StatusBadRequest, errMessage)
		}

		return c.Next()
	}
}
