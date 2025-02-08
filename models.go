package main

import (
	"time"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/google/uuid"
)

type UserRes struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	Email      string `json:"email"`
}

func databaseUserToUser(dbUser database.User) UserRes {
	return UserRes{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}
}