package httpRouter

import (
	"net/http"

	handlerV1Router "github.com/FazylovAsylkhan/kuryltai-server/handler/http/v1"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)


func Init(DB *database.Queries, secretKey string) *chi.Mux {
	router := chi.NewRouter()
	options := cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	router.Use(cors.Handler(options))

	v1Router := handlerV1Router.Init(DB, secretKey)
	router.Mount("/v1", v1Router)

	router.Handle("/assets/image/*", http.StripPrefix("/assets/image", http.FileServer(http.Dir("./assets/image"))))
	router.Handle("/assets/video/*", http.StripPrefix("/assets/video", http.FileServer(http.Dir("./assets/video"))))
	router.Handle("/assets/audio/*", http.StripPrefix("/assets/audio", http.FileServer(http.Dir("./assets/audio"))))

	return router
}
