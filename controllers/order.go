package controllers

import (
	"synapsis/services"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	orderService    services.OrderService
	midtransService services.MidtransService
}

func NewOrderController(orderService services.OrderService, midtransService services.MidtransService) *OrderController {
	return &OrderController{orderService: orderService, midtransService: midtransService}
}

func (c *OrderController) Checkout(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	order, err := c.orderService.Checkout(userID)
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, err.Error())
	}

	res, _ := c.midtransService.SnapRequest(order.ID)

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, res)
}

func (c *OrderController) GetOrders(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	orders, err := c.orderService.GetOrdersByUserID(userID)
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, orders)
}

func (c *OrderController) GetOrder(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id")

	userID := ctx.Locals("userID").(uint)
	order, err := c.orderService.GetOrderByIDAndUserID(userID, orderID)
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, order)
}

func (c *OrderController) ConfirmOrder(ctx *fiber.Ctx) error {
	orderID := ctx.Params("id")

	err := c.orderService.ConfirmOrder(orderID)
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, nil)
}
