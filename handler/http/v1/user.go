package httpRouterV1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/user"
)

func handlerSignUp(w http.ResponseWriter, r *http.Request) {
	params := user.SignUpReq{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	res, err := user.SignUp(r, params.Password, params.Email)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}
	RespondWithJSON(w, 201, res)
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	params := user.LoginReq{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	userData, err := user.CheckPassword(r, params.Password, params.Email)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	res, err := user.Login(r, userData, params.Email)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, 200, res)
}
func handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	params := user.RefreshTokenReq{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	refreshClaims, err := user.CheckRefreshToken(r, params.RefreshToken)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}
	res, err := user.CreateAccessToken(refreshClaims)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	RespondWithJSON(w, 200, res)
}

func handlerLogoutUser(w http.ResponseWriter, r *http.Request, userData database.User) {
	err := user.LogoutUser(r, userData)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}
	RespondWithJSON(w, 200, struct{}{})
}

func handlerUpdatePassword(w http.ResponseWriter, r *http.Request, userData database.User) {
	params := user.UpdatePasswordReq{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	hashedPassword, err := user.ValidatePassword(params.OldPassword, params.NewPassword, userData)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return 
	}
	
	err = user.UpdatePassword(r, hashedPassword, userData)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return 
	}

	RespondWithJSON(w, 200, struct{ message string }{message: "Password changed successful"})
}