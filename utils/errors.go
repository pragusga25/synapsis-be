package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrDatabaseOperationFailed = errors.New("database operation failed")
	ErrProductNotFound         = errors.New("product not found")
	ErrCategoryNotFound        = errors.New("category not found")
	ErrProductQuantityExceeded = errors.New("product quantity exceeded")
	ErrEmptyCart               = errors.New("cart is empty")
	ErrTransactionNotFound     = errors.New("transaction not found")
	ErrTransactionFailed       = errors.New("transaction failed")
	ErrEmptyOrderID            = errors.New("order id is empty")
	ErrOrderNotFound           = errors.New("order does not exist")
)

type CustomError struct {
	Status int
	Inner  error
}

// Error is mark the struct as an error.
func (e *CustomError) Error() string {
	return e.Inner.Error()
}

func WrapWithCustomeError(err error, status int) error {
	return &CustomError{
		Inner:  err,
		Status: status,
	}
}

func PostgresErrorHandler(err error) error {
	var pgErr *pgconn.PgError

	if ok := errors.As(err, &pgErr); !ok {
		return ErrDatabaseOperationFailed
	}

	return pgErr
}
