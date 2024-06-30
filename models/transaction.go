package models

import "time"

type Transaction struct {
	ID              string `gorm:"type:varchar(36);primaryKey"`
	Amount          float64
	Status          string
	OrderID         string
	UserID          uint
	User            User  `json:"-"`
	Order           Order `json:"-"`
	TransactionType string
	PaymentType     string
	PaymentCode     string
	ApprovalCode    string
	CardType        string
	Bank            string
	VANumber        string
	MaskedCard      string
	Issuer          string
	Acquirer        string
	Currency        string
	Store           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
