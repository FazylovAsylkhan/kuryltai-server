package main

import (
	"fmt"
	"net/http"

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