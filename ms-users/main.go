package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/gorm"

	"ms-user/database"
	"ms-user/services"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

func NewRouter() *mux.Router {
	return mux.NewRouter()
}

func registerRoutes(router *mux.Router, db *gorm.DB) {
	router.HandleFunc("/register", services.RegisterUser(db)).Methods("POST")
	router.HandleFunc("/user", services.GetUserById(db)).Methods("GET")
	router.HandleFunc("/login", services.Login(db)).Methods("POST")
}

func startServer(router *mux.Router) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
