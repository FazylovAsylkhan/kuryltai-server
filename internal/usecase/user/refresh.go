package user

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/token"
)


type RefreshTokenReq struct {
	RefreshToken string `json:"refresh"`
}

type RefreshTokenRes struct {
	AccessToken string `json:"access"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func CheckRefreshToken(r *http.Request, refreshToken string)  (*token.UserClaims, error) {
	refreshClaims, err := apiCfg.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("Error verifying token: %v", err)
	}

	session, err := apiCfg.DB.GetSession(r.Context(), refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return nil, fmt.Errorf("Error getting session: %v", err)
	}
	if session.IsRevoked {
		return nil, errors.New("Session revoked")
	}
	if session.UserEmail != refreshClaims.Email {
		return nil, errors.New("Invalid session")
	}
	return refreshClaims, nil
}

func CreateAccessToken(refreshClaims *token.UserClaims) (RefreshTokenRes, error) {
	accessToken, accessClaims, err := apiCfg.tokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, 15 * time.Minute)
	if err != nil {
		return RefreshTokenRes{}, fmt.Errorf("error creating access token: %v", err)
	}
	res := RefreshTokenRes{
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}
	return res, nil
}
