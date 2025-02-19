package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/auth"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/slug"
	"github.com/google/uuid"
)

type SignUpReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRes struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	Email      string `json:"email"`
}

func SignUp(r *http.Request, password string, email string) (SignUpRes, error) {
	hashed, err := auth.HashPassword(password) // Changed from auth.hashedPassword
	if err != nil {
		return SignUpRes{}, fmt.Errorf("Error hashing password: %v", err)
	}
	userID := uuid.New()
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email:     email,
		Password:  hashed,
	})
	if err != nil {
		return SignUpRes{}, fmt.Errorf("Couldn't create user: %v", err)
	}

	_, err = apiCfg.DB.CreateProfile(r.Context(), database.CreateProfileParams{
		ID:       uuid.New(),
		UserID:   userID,
		Username: slug.GetUsernameFromEmail(user.Email), // Changed from auth.GetUsernameFrom
		Slug:     slug.GenerateSlug(), // Changed from auth.GenerateRandomSlug
	})
	if err != nil {
		return SignUpRes{}, fmt.Errorf("Couldn't create profile: %v", err)
	}

	return getUserRes(user), nil
}

func getUserRes(dbUser database.User) SignUpRes {
	return SignUpRes{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}