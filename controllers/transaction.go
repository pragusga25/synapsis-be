package controllers

import (
	"synapsis/services"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	transactionService services.TransactionService
	midtransService    services.MidtransService
}

func NewTransactionController(transactionService services.TransactionService, midtransService services.MidtransService) *TransactionController {
	return &TransactionController{transactionService: transactionService, midtransService: midtransService}
}

func (c *TransactionController) CreateTransaction(ctx *fiber.Ctx) error {
	var req map[string]interface{}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, err.Error())
	}

	err := c.midtransService.HandleNotificationPayload(req)

	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, nil)

}

func (c *TransactionController) GetTransactions(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	transactions, err := c.transactionService.GetTransactionsByUserID(userID)
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, transactions)
}
