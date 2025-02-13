package main

import (
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/util"
)

type  authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		accessToken, err := util.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
	
		accessClaims, err := apiCfg.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			respondWithError(w, 401, fmt.Sprintf("Error verifying token: %v", err))
			return
		} 

		user, err := apiCfg.DB.GetUser(r.Context(), accessClaims.Email)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}