package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/util"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerSignup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	hashed, err := util.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error hashing password: %v", err))
		return
	}

	userID := uuid.New()
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email: params.Email,
		Password: hashed,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	_, err = apiCfg.DB.CreateProfile(r.Context(), database.CreateProfileParams{
		ID: uuid.New(),
		UserID: userID,
		Username: util.GetUsernameFrom(user.Email),
		Slug: util.GenerateRandomSlug(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't find user: %v", err))
		return
	}
	err = util.CheckPassword(params.Password, user.Password)
	fmt.Println(err)
	if err != nil {
		respondWithError(w, 400, "wrong password")
		return
	}

	accessToken, accessClaims, err := apiCfg.tokenMaker.CreateToken(user.ID, user.Email, 15 * time.Minute)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error creating access token: %v", err))
		return
	}
	refreshToken, refreshClaims, err := apiCfg.tokenMaker.CreateToken(user.ID, user.Email, 24 * time.Hour)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error creating refresh token: %v", err))
		return
	}
	session, err := apiCfg.DB.CreateSession(r.Context(), database.CreateSessionParams{
		ID: refreshClaims.RegisteredClaims.ID,
		UserID: user.ID,
		UserEmail: user.Email,
		RefreshToken: refreshToken,
		ExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error creating session: %v", err))
		return
	}

	res := LoginUserRes{
		SessionID: session.ID,
		RefreshToken: session.RefreshToken,
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User: databaseUserToUser(user),
	}

	respondWithJSON(w, 200, res)
}

func (apiCfg *apiConfig) handlerUpdatePassword(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	err = util.CheckPassword(params.OldPassword, user.Password)
	if err != nil {
		respondWithError(w, 400, "wrong password")
		return
	}

	hashed, err := util.HashPassword(params.NewPassword)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error hashing password: %v", err))
		return
	}

	user, err = apiCfg.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		ID: user.ID,
		Password: hashed,
	})
	fmt.Println(user)
	fmt.Println(err)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't change password: %s", err))
		return
	}

	respondWithJSON(w, 200, struct{message string}{message: "Password changed successful"})
}

func (apiCfg *apiConfig) handlerLogoutUser(w http.ResponseWriter, r *http.Request, user database.User) {
    _, err := apiCfg.DB.DeleteSession(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting session: %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}

func (apiCfg *apiConfig) renewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req RenewAccessTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	refreshClaims, err := apiCfg.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error verifying token: %v", err))
		return
	}

	session, err := apiCfg.DB.GetSession(r.Context(), refreshClaims.RegisteredClaims.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting session: %v", err))
		return
	}
	if session.IsRevoked {
		respondWithError(w, 400, "Session revoked")
		return
	}
	if session.UserEmail != refreshClaims.Email {
		respondWithError(w, 400, "Invalid session")
		return
	}

	accessToken, accessClaims, err := apiCfg.tokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, 15 * time.Minute)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("error creating access token: %v", err))
		return
	}
	res := RenewAccessTokenRes{
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}
	respondWithJSON(w, 200, res)
}

func (apiCfg *apiConfig) revokeSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, 400, "Missing session ID")
		return
	}
    _, err := apiCfg.DB.RevokeSession(r.Context(), id)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error revoking session: %v", err))
		return
	}	
	respondWithJSON(w, 200, struct{}{})
}