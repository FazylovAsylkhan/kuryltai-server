package httpRouterV1

import (
	"fmt"
	"net/http"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	authd "github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/auth"
)

type  authedHandler func(http.ResponseWriter, *http.Request, database.User)

func MiddlewareAuth(handler authedHandler, routerCfg *RouterConfig) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		accessToken, err := authd.GetAPIKey(r.Header)
		if err != nil {
			RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
	
		accessClaims, err := routerCfg.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			RespondWithError(w, 401, fmt.Sprintf("Error verifying token: %v", err))
			return
		} 

		user, err := routerCfg.DB.GetUser(r.Context(), accessClaims.Email)
		if err != nil {
			RespondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}