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

	// TODO: currently FileSystemStore is used since CookieStore doesn't able to store cookie of larger size
	sessionStore := auth.NewFileSystemStore(auth.SessionOptions{
		CookiesKey: config.Envs.CookiesAuthSecret,
		MaxAge:     config.Envs.CookiesAuthAgeInSeconds,
		Secure:     config.Envs.CookiesAuthIsSecure,
		HttpOnly:   config.Envs.CookiesAuthIsHttpOnly,
	})

	auth.NewAuthService(sessionStore)

	router := mux.NewRouter()

	handler := handlers.New(store)

	// static assets
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// blog posts
	router.HandleFunc("/blog-post", auth.RequireAuth(handler.HandleShowAllBlogPostsPage)).Methods("GET")
	router.HandleFunc("/blog-post/new", auth.RequireAuth(handler.HandleNewBlogPostPage)).Methods("GET")
	router.HandleFunc("/blog-post", auth.RequireAuth(handler.HandleNewBlogPost)).Methods("POST")
	router.HandleFunc("/blog-post/edit/{id}", auth.RequireAuth(handler.HandleBlogPostEditPage)).Methods("GET")
	router.HandleFunc("/blog-post/{id}", auth.RequireAuth(handler.HandleDeleteBlogPost)).Methods("DELETE")
	router.HandleFunc("/blog-post/{id}", auth.RequireAuth(handler.HandleUpdateBlogPost)).Methods("PUT")

	// auth
	router.HandleFunc("/login", handler.HandleLoginPage).Methods("GET")
	router.HandleFunc("/auth/{provider}", handler.HandleProviderLogin).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", handler.HandleAuthCallbackFunction).Methods("GET")
	router.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")

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
