package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/handlerContent"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
)

func (apiCfg *apiConfig) handlerGetProfile(w http.ResponseWriter, r *http.Request, user database.User) {
	profile, err := apiCfg.DB.GetProfile(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get profile: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseProfileToProfile(profile))
}

func (apiCfg *apiConfig) handlerUpdateProfile(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)
	params := UpdateReq{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	profile, _ :=apiCfg.DB.GetProfile(r.Context(), user.ID)
	parameters := database.UpdateProfileParams{
		UserID: user.ID,
		Slug: params.Slug,
		Username: params.Username,
		HeadLine: sql.NullString{
			Valid: true,
			String: params.HeadLine,
		},
		AvatarImage: sql.NullString{
			Valid: true,
			String: profile.AvatarImage.String,
		},
		CoverImage: sql.NullString{
			Valid: true,
			String: profile.CoverImage.String,
		},
	}
	if params.AvatarImage != nil {
		link, err := handlerContent.UploadFile(user.ID, *params.AvatarImage)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't upload avatar image: %v", err))
			return
		}
		parameters.AvatarImage = sql.NullString{
			Valid: true,
			String: link,
		}
	}
	if params.CoverImage != nil {
		link, err := handlerContent.UploadFile(user.ID, *params.CoverImage)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't upload cover image: %v", err))
			return
		}
		parameters.CoverImage = sql.NullString{
			Valid: true,
			String: link,
		}
	}
	
	profile, err = apiCfg.DB.UpdateProfile(r.Context(), parameters)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't update profile: %v", err))
		return
	}
	
	respondWithJSON(w, 200, databaseProfileToProfile(profile))
}