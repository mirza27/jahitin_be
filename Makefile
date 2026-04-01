

migration:
	migrate create -ext sql -dir database/migration -seq

migrate: 
	migrate -path database/migration -database "postgresql://postgres:da7sduhasd@localhost:5400/jahitin_db?sslmode=disable" -verbose up

rollback:
	migrate -path database/migration -database "postgresql://postgres:da7sduhasd@localhost:5400/jahitin_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

sqlc:
	sqlc generate