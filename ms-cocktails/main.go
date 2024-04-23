package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/gorm"

	"ms-cocktails/database"
	"ms-cocktails/services"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func NewRouter() *mux.Router {
	return mux.NewRouter()
}

func registerRoutes(router *mux.Router, db *gorm.DB) {
	router.HandleFunc("/cocktails", services.CocktailHandler).Methods("GET")
	router.HandleFunc("/cocktails/detail", services.CocktailDetail).Methods("GET")
	router.HandleFunc("/cocktails/user", services.CocktailsByUser(db)).Methods("GET")
	router.HandleFunc("/cocktails", services.AddCocktail(db)).Methods("POST")
	router.HandleFunc("/cocktails/{cocktail_id}/user/{user_id}", services.RemoveFavorite(db)).Methods("DELETE")
}

func startServer(router *mux.Router) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "Authorization"}),
		handlers.AllowCredentials(),
	)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: cors(router),
	}

	log.Printf("Listening on http://localhost:%s/", port)
	log.Fatal(server.ListenAndServe())
}

func main() {
	fx.New(
		database.Module,
		fx.Provide(NewRouter),
		fx.Invoke(registerRoutes, startServer),
	).Run()
}
