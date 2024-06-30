package services

import (
	"synapsis/models"
	"synapsis/repositories"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
)

type TransactionService interface {
	CreateTransaction(Transaction *models.Transaction) error
	GetTransactionsByUserID(userID uint) ([]models.Transaction, error)
}

type transactionService struct {
	transactionRepository repositories.TransactionRepository
	orderRepository       repositories.OrderRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, orderRepo repositories.OrderRepository) TransactionService {
	return &transactionService{transactionRepository: transactionRepo, orderRepository: orderRepo}
}

func (s *transactionService) CreateTransaction(transaction *models.Transaction) error {
	order, err := s.orderRepository.GetOrderByID(transaction.OrderID)

	if err != nil || order.TotalPrice != transaction.Amount {
		return utils.WrapWithCustomeError(utils.ErrTransactionFailed, fiber.StatusInternalServerError)
	}

	order.Status = models.OrderStatusConfirmed
	err = s.orderRepository.UpdateOrder(order)

	if err != nil {
		return err
	}

	return s.transactionRepository.CreateTransaction(transaction)
}

func (s *transactionService) GetTransactionsByUserID(userID uint) ([]models.Transaction, error) {
	transactions, err := s.transactionRepository.GetTransactionsByUserID(userID)

	if err != nil {
		return nil, utils.WrapWithCustomeError(utils.ErrTransactionNotFound, fiber.StatusNotFound)
	}

	return transactions, nil
}
