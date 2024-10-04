package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Open-Code-Zone/cms/config"
	"github.com/Open-Code-Zone/cms/handlers"
	"github.com/Open-Code-Zone/cms/services/auth"
	"github.com/Open-Code-Zone/cms/store"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// SQLite config
	db, err := store.NewSQLiteStorage("cms.db")
	if err != nil {
		log.Fatal(err)
	}

	// db store
	store := store.NewStore(db)
	initStorage(db)

	// session store and auth service
  sessionStore := auth.NewCoockieOptions(auth.SessionOptions{
    CookiesKey: config.Envs.CookiesAuthSecret,
    MaxAge: config.Envs.CookiesAuthAgeInSeconds,
    Secure: config.Envs.CookiesAuthIsSecure,
    HttpOnly: config.Envs.CookiesAuthIsHttpOnly,
  })

  authService := auth.NewAuthService(sessionStore)


	router := mux.NewRouter()

	handler := handlers.New(store, authService)

	// static assets
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// blog posts
	router.HandleFunc("/blog-post", auth.RequireAuth(handler.HandleShowAllBlogPostsPage, authService)).Methods("GET")
	router.HandleFunc("/blog-post/new", handler.HandleNewBlogPostPage).Methods("GET")
	router.HandleFunc("/blog-post", handler.HandleNewBlogPost).Methods("POST")
	router.HandleFunc("/blog-post/edit/{id}", handler.HandleBlogPostEditPage).Methods("GET")
	router.HandleFunc("/blog-post/{id}", handler.HandleDeleteBlogPost).Methods("DELETE")
	router.HandleFunc("/blog-post/{id}", handler.HandleUpdateBlogPost).Methods("PUT")

	// auth
	router.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	router.HandleFunc("/auth/logout/{provider}", nil).Methods("GET")
	router.HandleFunc("/login", handler.HandleLogin).Methods("GET")

	log.Printf("Server: Listening on %s:%s\n", config.Envs.PublicHost, config.Envs.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", config.Envs.Port), router))
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database!")
}
