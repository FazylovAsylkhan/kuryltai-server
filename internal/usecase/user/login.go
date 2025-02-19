package user

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/auth"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User SignUpRes `json:"user"`
}

func CheckPassword(r *http.Request, password string, email string)  (database.User, error) {
	user, err := apiCfg.DB.GetUser(r.Context(), email)
	if err != nil {
		return database.User{}, fmt.Errorf("couldn't find user: %v", err)
	}
	err = auth.CheckPassword(password, user.Password)
	if err != nil {
		return database.User{}, errors.New("wrong password")
	}

	return user, nil
}

func Login(r *http.Request, user database.User, email string) (LoginRes, error) {
	accessToken, accessClaims, err := apiCfg.tokenMaker.CreateToken(user.ID, email, 15*time.Minute)
	if err != nil {
		return LoginRes{}, fmt.Errorf("error creating access token: %v", err)
	}

	refreshToken, refreshClaims, err := apiCfg.tokenMaker.CreateToken(user.ID, email, 24*time.Hour)
	if err != nil {
		return LoginRes{}, fmt.Errorf("error creating refresh token: %v", err)
	}

	session, err := apiCfg.DB.CreateSession(r.Context(), database.CreateSessionParams{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserID:       user.ID,
		UserEmail:    email,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		return LoginRes{}, fmt.Errorf("error creating session: %v", err)
	}

	res := LoginRes{
		RefreshToken:          session.RefreshToken,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                  getUserRes(user),
	}

	return res, nil
}