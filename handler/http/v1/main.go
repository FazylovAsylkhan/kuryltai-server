package httpRouterV1

import (
	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/token"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/profile"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/usecase/user"
	"github.com/go-chi/chi"
)

type RouterConfig struct{
	tokenMaker *token.JWTMaker
	DB *database.Queries
}

var routerCfg  = &RouterConfig{}

func Init(DB *database.Queries, secretKey string) *chi.Mux {	
	router := chi.NewRouter()
	
	routerCfg = &RouterConfig{
		tokenMaker: token.NewJWTMaker(secretKey),
		DB: DB,
	}

	user.Init(routerCfg.DB, routerCfg.tokenMaker)
	profile.Init(routerCfg.DB, routerCfg.tokenMaker)

	router.Post("/users/login", handlerLogin)
	router.Post("/users/signup", handlerSignUp)
	router.Post("/users/token/refresh", handlerRefreshToken)
	router.Delete("/users/logout", MiddlewareAuth(handlerLogoutUser, routerCfg))
	router.Post("/users/change-password", MiddlewareAuth(handlerUpdatePassword, routerCfg))
	
	router.Get("/profiles/profile/me", MiddlewareAuth(handlerGetProfile, routerCfg))
	router.Patch("/profiles/profile/edit", MiddlewareAuth(handlerUpdateProfile, routerCfg))

	return router
}
