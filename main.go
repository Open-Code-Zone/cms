package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Open-Code-Zone/cms/config"
	"github.com/Open-Code-Zone/cms/handlers"
	"github.com/Open-Code-Zone/cms/internal/database"
	"github.com/Open-Code-Zone/cms/services/auth"
	"github.com/Open-Code-Zone/cms/store"
	"github.com/Open-Code-Zone/cms/utils"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// SQLite config
	db, err := store.NewSQLiteStorage("cms.db")
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)
	// db store
	store := store.NewStore(queries, config.Envs.CollectionConfig)
	pingStorage(db)

	// TODO: currently FileSystemStore is used since CookieStore doesn't able to store cookie of larger size
	sessionStore := auth.NewFileSystemStore(auth.SessionOptions{
		CookiesKey: config.Envs.CookiesAuthSecret,
		MaxAge:     config.Envs.CookiesAuthAgeInSeconds,
		Secure:     config.Envs.CookiesAuthIsSecure,
		HttpOnly:   config.Envs.CookiesAuthIsHttpOnly,
	})

	// GitHub Client
	githubClient, err := utils.NewGitHubClient()
	if err != nil {
		log.Printf("Error creating GitHub client: %v", err)
	}

	auth.NewAuthService(sessionStore)
	router := mux.NewRouter()

	handler := handlers.New(store, githubClient)

	// static assets
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// auth
	router.HandleFunc("/login", handler.LoginPage).Methods("GET")                       // login page
	router.HandleFunc("/auth/{provider}", handler.ProviderLogin).Methods("GET")         // login with provider
	router.HandleFunc("/auth/{provider}/callback", handler.AuthCallback).Methods("GET") // callback from provider
	router.HandleFunc("/auth/logout/{provider}", handler.Logout).Methods("GET")         // logout

	// routes for every collection
	router.HandleFunc("/{collection}", auth.RequireAuth(handler.Index)).Methods("GET")           // all the posts
	router.HandleFunc("/{collection}/new", auth.RequireAuth(handler.New)).Methods("GET")         // new post page
	router.HandleFunc("/{collection}", auth.RequireAuth(handler.Create)).Methods("POST")         // create new post htmx endpoint
	router.HandleFunc("/{collection}/edit/{id}", auth.RequireAuth(handler.Edit)).Methods("GET")  // edit post page
	router.HandleFunc("/{collection}/{id}", auth.RequireAuth(handler.Destroy)).Methods("DELETE") // delete post htmx endpoint
	router.HandleFunc("/{collection}/{id}", auth.RequireAuth(handler.Update)).Methods("PUT")     // update post htmx endpoint

	log.Printf("Server: Listening on %s:%s\n", config.Envs.PublicHost, config.Envs.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", config.Envs.Port), router))
}

func pingStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database!")
}
