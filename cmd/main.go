package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alroymuhammad/go-commerce/internal/handler"
	"github.com/alroymuhammad/go-commerce/pkg/config"
)

func main() {
	config.ConnectDB()
	defer config.DB.Close()

	mux := http.NewServeMux()
	userHandler := handler.NewUserHandler(config.DB)
	mux.Handle("/users/", userHandler)
	mux.Handle("/users", userHandler)

	port := ":8080"
	fmt.Printf("Server starting on port%s...\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
