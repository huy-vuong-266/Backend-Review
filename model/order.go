package model

import uuid "github.com/satori/go.uuid"

type Order struct {
	OrderID       uuid.UUID `gorm:"PRIMARY_KEY"`
	CreatedAt     int64
	UpdatedAt     int64
	UserID        uuid.UUID
	Amount        int64
	Type          int8
	TransactionID uuid.UUID
	Status        int8
}

type FinOrderRequest struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
}
