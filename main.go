package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sougandhini/rssagg/internal/database"

	_ "github.com/lib/pq"
)

// Note: we are building this rss aggregator on chi server- a lightweight server
type apiConfig struct {
	DB *database.Queries //saying that this is a pointer to another structure called queries thats present in database folder
}

func main() {

	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in .env") // this is cut the program with exit status 1 and returns
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in .env") //
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Unable to connect to Database")
	}

	// this is my link to connect to DB through go
	apiCfg := apiConfig{
		DB: database.New(conn), // this takes type: *database.Queries but what we have is conn which is *sql.DB type, so we need to typecast it
	}

	router := chi.NewRouter() // Mother router

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
	v1Router.Get("/healthz", handlerReadiness) // scopes the handler to only fire on GET requests
	v1Router.Get("/error", handlerErr)

	// CRUD for users
	// 1. Create
	v1Router.Post("/users", apiCfg.handlerCreateUser) // now this handler will have access to DB

	// 2. Get
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser)) // we're using auth middleware only for those functions which requires authentication

	// Feed handlers
	// 1. Create feed
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))

	// 2. Get Feed
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	// Feed follow
	// 1. Create a feed follow
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))

	// 2. List all the feeds the the user is currently following
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))

	// 3. Delete feed follow
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router) // nesting v1 router for /v1 path and then see /healthz - first child router to mother router

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("server starting on port %v", portString)
	errListenNServe := srv.ListenAndServe() // this line will run forever
	if errListenNServe != nil {
		log.Fatal(errListenNServe)
	} // this means if anything goes wrong while handling the http request in the mentioned port we will log and return the program
}
