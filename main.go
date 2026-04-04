package main

import (
	"database/sql"
	"fmt"
	"jahitin_be/api"
	"jahitin_be/config"
	"jahitin_be/token"
	"log"

	db "jahitin_be/database/repository"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// connect to database
	dbURL := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	dbConn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	store := db.NewStore(dbConn)

	// token maker object
	tokenMaker, err := token.NewMaker(config.TokenSecretKey)
	if err != nil {
		log.Fatalf("failed to create token maker: %v", err)
	}

	server, err := api.NewApiServer(*config, store, *tokenMaker)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
