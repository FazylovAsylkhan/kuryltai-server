package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	httpRouter "github.com/FazylovAsylkhan/kuryltai-server/handler/http"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/config"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

type Service struct {
	Config *config.Config
	DB     *database.Queries
	Router *chi.Mux
}

func Init(cfg *config.Config) (*Service, error) {

	conn, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %v", err)
	}

	s := &Service{
		Config: cfg,
		DB:     database.New(conn),
	}
	s.Router = httpRouter.Init(s.DB, cfg.SecretKey)

	return s, nil
}

func (s Service) Start() {
	srv := &http.Server{
		Handler: s.Router,
		Addr:    ":" + s.Config.Port,
	}

	log.Printf("Server starting on port %v", s.Config.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}