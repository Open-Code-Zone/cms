package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"githum.com/Open-Code-Zone/cms/handlers"
	"githum.com/Open-Code-Zone/cms/store"
)

func main() {
	// SQLite config
	db, err := store.NewSQLiteStorage("cms.db")
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStore(db)
	initStorage(db)

	router := mux.NewRouter()
	handler := handlers.New(store)

	router.HandleFunc("/", handler.HandleHome).Methods("GET")
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	fmt.Printf("Listening on %v\n", "localhost:7000")
	http.ListenAndServe(":7000", router)
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database!")
}
