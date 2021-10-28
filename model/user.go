package model

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	UserID       uuid.UUID
	CreatedAt    int64
	UpdatedAt    int64
	Fullname     string
	Phone        string
	Email        string
	Encrypted_PW string `json:"-"`
	Salt         string `json:"-"`
	Budget       int64
	Status       int8
}
