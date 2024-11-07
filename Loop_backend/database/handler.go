package database

import (
	"log"
	"net/http"
)

func StartServer() {
	err := InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	//defer DB.Close()
}

func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Loop Backend API"))
}
