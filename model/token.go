package model

import uuid "github.com/satori/go.uuid"

type Token struct {
	UserID    uuid.UUID `gorm:"PRIMARY_KEY"`
	Token     string
	CreatedAt int64
}
