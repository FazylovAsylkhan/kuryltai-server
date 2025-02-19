package user

import (
	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/token"
)

type apiConfig struct {
	DB *database.Queries
	tokenMaker *token.JWTMaker
}


var apiCfg apiConfig

func Init(db *database.Queries, token *token.JWTMaker) {
	apiCfg.DB = db
	apiCfg.tokenMaker = token
}
