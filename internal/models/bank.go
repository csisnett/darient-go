package models

import "time"

type BankType string

const (
	BankTypePrivate    BankType = "PRIVATE"
	BankTypeGovernment BankType = "GOVERNMENT"
)

type Bank struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      BankType  `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}