package services

import (
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"synapsis/config"
	"synapsis/models"
	"synapsis/repositories"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	SnapRequest(orderID string) (*snap.Response, *midtrans.Error)
	VerifyMidtransSignature(orderID, statusCode, grossAmount, midtransSignature string) bool
	HandleNotificationPayload(notificationPayload map[string]interface{}) error
}

type midtransService struct {
	orderRepository       repositories.OrderRepository
	transactionRepository repositories.TransactionRepository
}

func NewMidtransService(orderRepo repositories.OrderRepository, transactionRepo repositories.TransactionRepository) MidtransService {
	return &midtransService{orderRepository: orderRepo, transactionRepository: transactionRepo}
}

func (s *midtransService) SnapRequest(orderID string) (*snap.Response, *midtrans.Error) {
	cfg, _ := config.LoadConfig()

	var midtransEnv = midtrans.Sandbox

	if cfg.Env == "production" {
		midtransEnv = midtrans.Production
	}

	order, err := s.orderRepository.GetOrderByID(orderID)
	if err != nil {
		return nil, &midtrans.Error{StatusCode: fiber.StatusNotFound, Message: utils.ErrOrderNotFound.Error()}
	}

	// 1. Initiate Snap client
	var sn = &snap.Client{}
	sn.New(cfg.MidtransServerKey, midtransEnv)
	// Use to midtrans.Production if you want Production Environment (accept real transaction).

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.ID,
			GrossAmt: int64(order.TotalPrice),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.User.Name,
			Email: order.User.Email,
		},
		UserId: strconv.Itoa(int(order.UserID)),
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	return sn.CreateTransaction(req)
}

func (s *midtransService) HandleNotificationPayload(notificationPayload map[string]interface{}) error {
	cfg, _ := config.LoadConfig()

	var midtransEnv = midtrans.Sandbox

	if cfg.Env == "production" {
		midtransEnv = midtrans.Production
	}

	// 1. Initiate Snap client
	var c = &coreapi.Client{}
	c.New(cfg.MidtransServerKey, midtransEnv)

	// Get order-id from payload
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		return utils.WrapWithCustomeError(utils.ErrEmptyOrderID, fiber.StatusBadRequest)
	}

	// Check transaction to Midtrans with param orderId
	transactionStatusResp, e := c.CheckTransaction(orderId)
	if e != nil {
		return e
	}

	if transactionStatusResp != nil {
		signatureKey := transactionStatusResp.SignatureKey

		isValid := s.VerifyMidtransSignature(orderId, transactionStatusResp.StatusCode, transactionStatusResp.GrossAmount, signatureKey)

		if !isValid {
			return utils.WrapWithCustomeError(utils.ErrTransactionFailed, fiber.StatusBadRequest)
		}

		transactionID := transactionStatusResp.TransactionID
		amount, _ := strconv.ParseFloat(transactionStatusResp.GrossAmount, 64)
		transactionStatus := transactionStatusResp.TransactionStatus
		bank := transactionStatusResp.Bank

		vaNumbers := transactionStatusResp.VaNumbers
		var vaNumber string

		if vaNumbers != nil {
			bank = vaNumbers[0].Bank
			vaNumber = vaNumbers[0].VANumber
		}

		if transactionStatusResp.PermataVaNumber != "" {
			vaNumber = transactionStatusResp.PermataVaNumber
			bank = "permata"
		}

		order, err := s.orderRepository.GetOrderByID(orderId)

		if err != nil {
			return utils.WrapWithCustomeError(utils.ErrOrderNotFound, fiber.StatusNotFound)
		}

		// Get transaction from database
		transaction := &models.Transaction{
			ID:              transactionID,
			OrderID:         orderId,
			Amount:          amount,
			Status:          transactionStatus,
			Bank:            bank,
			VANumber:        vaNumber,
			PaymentType:     transactionStatusResp.PaymentType,
			MaskedCard:      transactionStatusResp.MaskedCard,
			CardType:        transactionStatusResp.CardType,
			Issuer:          transactionStatusResp.Issuer,
			Acquirer:        transactionStatusResp.Acquirer,
			Currency:        transactionStatusResp.Currency,
			PaymentCode:     transactionStatusResp.PaymentCode,
			ApprovalCode:    transactionStatusResp.ApprovalCode,
			Store:           transactionStatusResp.Store,
			TransactionType: transactionStatusResp.TransactionType,
			UserID:          order.UserID,
		}

		// Do set transaction status based on response from check transaction status
		if transactionStatus == "capture" {
			if transactionStatusResp.FraudStatus == "challenge" {
				// TODO set transaction status on your database to 'challenge'
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			} else if transactionStatusResp.FraudStatus == "accept" {
				// TODO set transaction status on your database to 'success'
				transaction.Status = "success"
				order.Status = models.OrderStatusConfirmed
			}
		} else if transactionStatus == "settlement" {
			// TODO set transaction status on your databaase to 'success'
			transaction.Status = "success"
			order.Status = models.OrderStatusConfirmed
		} else if transactionStatus == "deny" {
			// TODO you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
		} else if transactionStatus == "cancel" || transactionStatus == "expire" {
			// TODO set transaction status on your databaase to 'failure'
			transaction.Status = "failure"
			order.Status = models.OrderStatusFailed
		} else if transactionStatus == "pending" {
			// TODO set transaction status on your databaase to 'pending' / waiting payment
			transaction.Status = "pending"
		}

		if err := s.orderRepository.UpdateOrder(order); err != nil {
			return utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
		}

		if err := s.transactionRepository.UpsertTransaction(transaction); err != nil {
			return utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
		}
	}

	return nil
}

func (s *midtransService) VerifyMidtransSignature(orderID, statusCode, grossAmount, midtransSignature string) bool {
	cfg, _ := config.LoadConfig()
	serverKey := cfg.MidtransServerKey

	hash := sha512.New()
	hash.Write([]byte(orderID + statusCode + grossAmount + serverKey))
	signature := hex.EncodeToString(hash.Sum(nil))

	return signature == midtransSignature
}
