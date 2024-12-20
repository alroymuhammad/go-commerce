package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "commerce"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	}
	fmt.Println("Successfully connected to the database!")
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Printf("Starting server on %s", *addr)
	err = http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Go Commerce!")
}
