package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"jahitin_be/config"
	db "jahitin_be/database/repository"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := loadConfig()
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	queries := db.New(conn)

	if _, err := queries.SeedClothesCategories(ctx); err != nil {
		log.Fatalf("seed clothes categories: %v", err)
	}

	if _, err := queries.SeedServiceTypes(ctx); err != nil {
		log.Fatalf("seed service types: %v", err)
	}

	fmt.Println("seed completed")
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig(".")
	if err == nil {
		return cfg
	}

	port, convErr := strconv.Atoi(defaultValue(os.Getenv("DB_PORT"), "5400"))
	if convErr != nil {
		port = 5400
	}

	return &config.Config{
		DBHost:     defaultValue(os.Getenv("DB_HOST"), "localhost"),
		DBPort:     port,
		DBUser:     defaultValue(os.Getenv("DB_USER"), "postgres"),
		DBPassword: defaultValue(os.Getenv("DB_PASSWORD"), "da7sduhasd"),
		DBName:     defaultValue(os.Getenv("DB_NAME"), "jahitin_db"),
	}
}

func defaultValue(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
