package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/FazylovAsylkhan/kuryltai-server/token"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
	tokenMaker *token.JWTMaker
}

const minSecretKeySize = 32

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB is not found in the environment")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not found in the environment")
	}
	if len(strings.Split(secretKey, "")) < minSecretKeySize {
		log.Fatalf("Secret_KEY must be at least %d characters", minSecretKeySize)
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
		tokenMaker: token.NewJWTMaker(secretKey),
	} 

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))
	v1Router := chi.NewRouter()
	v1Router.Post("/users/login", apiCfg.handlerLogin)
	v1Router.Post("/users/signup", apiCfg.handlerSignup)
	v1Router.Post("/users/token/refresh", apiCfg.renewAccessToken)
	v1Router.Post("/users/token/revoke", apiCfg.revokeSession)
	v1Router.Delete("/users/logout", apiCfg.middlewareAuth(apiCfg.handlerLogoutUser))
	v1Router.Post("/users/change-password", apiCfg.middlewareAuth(apiCfg.handlerUpdatePassword))
	
	v1Router.Get("/profiles/profile/me", apiCfg.middlewareAuth(apiCfg.handlerGetProfile))
	v1Router.Patch("/profiles/profile/edit", apiCfg.middlewareAuth(apiCfg.handlerUpdateProfile))


	router.Handle("/assets/image/*", http.StripPrefix("/assets/image", http.FileServer(http.Dir("./assets/image"))))
	router.Handle("/assets/video/*", http.StripPrefix("/assets/video", http.FileServer(http.Dir("./assets/video"))))
	router.Handle("/assets/audio/*", http.StripPrefix("/assets/audio", http.FileServer(http.Dir("./assets/audio"))))

	router.Mount("/v1", v1Router)
	
	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	 srv.ListenAndServe()
}
