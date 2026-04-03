package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"jahitin_be/config"
	db "jahitin_be/database/repository"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
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
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5400)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "da7sduhasd")
	v.SetDefault("DB_NAME", "jahitin_db")

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) {
			log.Printf("warning: unable to read .env config: %v", err)
		}
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Printf("warning: unable to unmarshal config, using defaults/env values: %v", err)
	}

	if cfg.DBPort == 0 {
		cfg.DBPort = v.GetInt("DB_PORT")
	}

	return &cfg
}
