package models

import "time"

type CreditType string
type CreditStatus string

const (
	CreditTypeAuto       CreditType = "AUTO"
	CreditTypeMortgage   CreditType = "MORTGAGE"
	CreditTypeCommercial CreditType = "COMMERCIAL"
)

const (
	CreditStatusPending  CreditStatus = "PENDING"
	CreditStatusApproved CreditStatus = "APPROVED"
	CreditStatusRejected CreditStatus = "REJECTED"
)

type Credit struct {
	ID         int          `json:"id"`
	ClientID   int          `json:"client_id"`
	BankID     int          `json:"bank_id"`
	MinPayment float64      `json:"min_payment"`
	MaxPayment float64      `json:"max_payment"`
	TermMonths int          `json:"term_months"`
	CreditType CreditType   `json:"credit_type"`
	Status     CreditStatus `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
}