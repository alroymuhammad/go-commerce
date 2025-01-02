package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DB *sql.DB

func ConnectDB() *sql.DB { // Changed return type to *sql.DB
	dbConfig := DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "admin",
		DBName:   "commerce",
		SSLMode:  "disable",
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	fmt.Println("Successfully connected to database!")

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	return DB // Return the DB connection
}
