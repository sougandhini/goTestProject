package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

// NOte: we are building this rss aggregator on chi server- a lightweight server
func main() {

	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in .env") // this is cut the program and returns
	}

	router := chi.NewRouter()

	// this router.Use is wrt what request should be given response to
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	v1Router := chi.NewRouter()
	// v1Router.HandleFunc("/healthz", handlerReadiness) // this handleFunc allows all type of api requests post, put, delete everything same... but we want healthz to only be accessed to GET so change it to next line....
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerErr)

	router.Mount("/v1", v1Router) // nesting v1 router for /v1 path and then see /healthz

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("server starting on port %v", portString)
	err := srv.ListenAndServe() // this line will run forever
	if err != nil {
		log.Fatal(err)
	} // this means if anything goes wrong while handing the http request in the mentioned port we will log and return the program
}
