package httpRouterV1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/profile"
)

func handlerGetProfile(w http.ResponseWriter, r *http.Request, userData database.User) {
	res, err := profile.GetProfile(r, userData)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}
	RespondWithJSON(w, 200, res)
}

func handlerUpdateProfile(w http.ResponseWriter, r *http.Request, userData database.User) {
	params := profile.UpdateProfileReq{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	res, err := profile.UpdateProfile(r, params, userData)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}
	RespondWithJSON(w, 200, res)
}