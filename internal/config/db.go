package config

import (
	"database/sql"
	"errors"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
)

func ConnectToDB(dbUrl string) (*database.Queries, error) {

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, errors.New("can't connect to database")
	}

	return database.New(conn), nil
}