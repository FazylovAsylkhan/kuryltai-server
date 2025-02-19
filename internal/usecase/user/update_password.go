package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	authd "github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/auth"
)

type UpdatePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ValidatePassword(oldPassword string, newPassword string, user database.User) (string, error) {
	err := authd.CheckPassword(oldPassword, user.Password)
	if err != nil {
		return "", errors.New("wrong password")
	}
	hashed, err := authd.HashPassword(newPassword)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %v", err)
	}

	return hashed, nil
}


func UpdatePassword(r *http.Request, hashedPassword string, user database.User) error {

	user, err := apiCfg.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		ID:       user.ID,
		Password: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("couldn't change password: %s", err)
	}
	return nil
}