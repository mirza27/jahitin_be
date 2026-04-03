package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
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

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("begin tx: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	queries := db.New(tx)

	if err := truncateSeedTables(ctx, tx); err != nil {
		log.Fatalf("truncate tables: %v", err)
	}

	categories, err := queries.SeedClothesCategories(ctx)
	if err != nil {
		log.Fatalf("seed clothes categories: %v", err)
	}

	serviceTypes, err := queries.SeedServiceTypes(ctx)
	if err != nil {
		log.Fatalf("seed service types: %v", err)
	}

	seedRand := rand.New(rand.NewSource(42))
	statusOptions := []string{"pending", "in_progress", "completed", "picked_up"}

	users, err := seedUsers(ctx, queries)
	if err != nil {
		log.Fatalf("seed users: %v", err)
	}

	for userIdx, user := range users {
		customers, err := seedCustomers(ctx, queries, user.ID, userIdx)
		if err != nil {
			log.Fatalf("seed customers for user %d: %v", user.ID, err)
		}

		for orderIdx := 0; orderIdx < 5; orderIdx++ {
			customer := customers[orderIdx%len(customers)]
			status := statusOptions[(userIdx+orderIdx)%len(statusOptions)]
			orderName := fmt.Sprintf("Order %02d-%02d", userIdx+1, orderIdx+1)

			order, err := queries.CreateOrder(ctx, db.CreateOrderParams{
				UserID:     user.ID,
				Name:       orderName,
				CustomerID: customer.ID,
				Deadline:   sql.NullTime{Valid: false},
				Status:     status,
			})
			if err != nil {
				log.Fatalf("seed order for user %d: %v", user.ID, err)
			}

			itemCount := 1 + seedRand.Intn(3)
			if err := seedOrderItems(ctx, queries, order.ID, itemCount, categories, serviceTypes, seedRand); err != nil {
				log.Fatalf("seed order items for order %d: %v", order.ID, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("commit seed tx: %v", err)
	}

	fmt.Println("seed completed")
}

func truncateSeedTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `TRUNCATE TABLE
		clothes_categories,
		service_types,
		notification_logs,
		order_items,
		orders,
		customers,
		users
	RESTART IDENTITY CASCADE;`)
	return err
}

func seedUsers(ctx context.Context, queries *db.Queries) ([]db.User, error) {
	users := make([]db.User, 0, 10)

	for i := 1; i <= 10; i++ {
		user, err := queries.CreateBaseUser(ctx, db.CreateBaseUserParams{
			Name:     fmt.Sprintf("User %02d", i),
			DeviceID: sql.NullString{String: fmt.Sprintf("device-%02d-%s", i, strings.Repeat("x", i%3+1)), Valid: true},
			UserType: "tailor",
			Phone:    sql.NullString{Valid: false},
		})
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func seedCustomers(ctx context.Context, queries *db.Queries, userID int64, userIdx int) ([]db.Customer, error) {
	customerCount := 2 + (userIdx % 2)
	customers := make([]db.Customer, 0, customerCount)

	for i := 1; i <= customerCount; i++ {
		var phone sql.NullString
		if (userIdx+i)%2 == 0 {
			phone = sql.NullString{String: fmt.Sprintf("08%02d%02d%02d%02d", userIdx+1, i, userIdx+2, i+3), Valid: true}
		}

		var notes sql.NullString
		if (userIdx+i)%3 != 0 {
			notes = sql.NullString{String: fmt.Sprintf("Customer note %d-%d", userIdx+1, i), Valid: true}
		}

		customer, err := queries.CreateCustomer(ctx, db.CreateCustomerParams{
			Name:   fmt.Sprintf("Customer %02d-%02d", userIdx+1, i),
			UserID: userID,
			Phone:  phone,
			Notes:  notes,
		})
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func seedOrderItems(
	ctx context.Context,
	queries *db.Queries,
	orderID int64,
	itemCount int,
	categories []db.ClothesCategory,
	serviceTypes []db.ServiceType,
	rnd *rand.Rand,
) error {
	statusOptions := []string{"pending", "in_progress", "completed", "picked_up"}

	for i := 1; i <= itemCount; i++ {
		category := categories[rnd.Intn(len(categories))]
		serviceType := serviceTypes[rnd.Intn(len(serviceTypes))]
		status := statusOptions[(int(orderID)+i)%len(statusOptions)]

		customServiceName := sql.NullString{Valid: false}
		categoryID := sql.NullInt64{Int64: category.ID, Valid: true}
		serviceTypeID := sql.NullInt64{Int64: serviceType.ID, Valid: true}

		if i%3 == 0 {
			customServiceName = sql.NullString{String: fmt.Sprintf("Custom Service %d", i), Valid: true}
			serviceTypeID = sql.NullInt64{Valid: false}
		}

		var notes sql.NullString
		if i%2 == 0 {
			notes = sql.NullString{String: fmt.Sprintf("Notes for order %d item %d", orderID, i), Valid: true}
		}

		_, err := queries.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:           orderID,
			CategoryID:        categoryID,
			ServiceTypeID:     serviceTypeID,
			ClothesFor:        fmt.Sprintf("Customer item %d", i),
			CustomServiceName: customServiceName,
			Notes:             notes,
			Price:             int64(50000 + rnd.Intn(150000)),
			Status:            status,
		})
		if err != nil {
			return err
		}
	}

	return nil
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
