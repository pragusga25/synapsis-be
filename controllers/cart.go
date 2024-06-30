package controllers

import (
	"synapsis/services"
	"synapsis/utils"
	"synapsis/validations"

	"github.com/gofiber/fiber/v2"
)

type CartController struct {
	cartService services.CartService
}

func NewCartController(cartService services.CartService) *CartController {
	return &CartController{cartService: cartService}
}

func (c *CartController) AddToCart(ctx *fiber.Ctx) error {
	adds := new(validations.AddToCartDto)

	if err := ctx.BodyParser(&adds); err != nil {
		utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	userID := ctx.Locals("userID").(uint)

	if err := c.cartService.AddToCart(userID, adds.ProductID, adds.Quantity); err != nil {
		cerr, ok := err.(*utils.CustomError)
		if ok {
			return utils.ErrorResponseJSON(ctx, cerr.Status, cerr.Inner.Error())
		}
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, nil)
}

func (c *CartController) GetCart(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	cart, err := c.cartService.GetCartByUserID(userID)
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, "Failed to retrieve cart")
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, cart)
}

func (c *CartController) UpdateCart(ctx *fiber.Ctx) error {
	updates := new(validations.UpdateCartDto)

	if err := ctx.BodyParser(&updates); err != nil {
		utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	userID := ctx.Locals("userID").(uint)

	if err := c.cartService.UpdateCart(userID, updates.ProductID, updates.Quantity); err != nil {
		cerr, ok := err.(*utils.CustomError)
		if ok {
			return utils.ErrorResponseJSON(ctx, cerr.Status, cerr.Inner.Error())
		}
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, nil)
}
