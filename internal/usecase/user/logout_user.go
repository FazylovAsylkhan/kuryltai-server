package user

import (
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
)

func LogoutUser(r *http.Request, user database.User) error {
	_, err := apiCfg.DB.DeleteSession(r.Context(), user.ID)

	if err != nil {
		return fmt.Errorf("error deleting session: %v", err)
	}
	
	return nil
}