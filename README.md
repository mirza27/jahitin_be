# Jahitin Backend

Backend service for a home-based tailoring business management mobile application. The API is intended to be consumed by a Flutter mobile client.

## Overview

Jahitin Backend provides authenticated HTTP endpoints for managing tailor users, customers, orders, order items, clothes categories, and service types. It is designed as the backend API for the [Jahitin mobile application](https://github.com/mirza27/safelyn_mobile).

The repository is not presented as production-ready. The current implementation is a work in progress and includes known gaps in automated testing and seed tooling.

## Key Features

- Local user registration and login — Implemented
- Account user registration and login — Implemented
- Local and account Bearer-token session authentication — Implemented
- Customer listing and detail retrieval — Implemented
- Order creation, listing, detail retrieval, update, deletion, and status update — Implemented
- Order-item replacement/update — Implemented
- Clothes category and service type listing — Implemented
- PostgreSQL schema migrations — Implemented

## Tech Stack

| Technology        | Purpose                                      |
| ----------------- | -------------------------------------------- |
| Go 1.25           | Backend language                             |
| Gin               | HTTP routing and request handling            |
| PostgreSQL        | Relational database                          |
| `pgx/v5`          | PostgreSQL driver through `database/sql`     |
| SQLC              | Type-safe generated database access code     |
| Viper             | Environment and `.env` configuration loading |
| HMAC-based tokens | Session token creation and verification      |

## Project Status

**Status:** In Development

The project is actively changing. API contracts, database operations, and supporting tooling may change as development continues.

| Area                                 | Status        |
| ------------------------------------ | ------------- |
| Authentication and session local     | Implemented   |
| Authentication and session account   | Unimplemented |
| User registration local              | Implemented   |
| User registration account            | Unimplemented |
| Customer management                  | Implemented   |
| Order and order-item management      | Implemented   |
| Clothes categories and service types | Implemented   |
| Order Public Monitoring              | Unimplemented |
| Order Invoicing Monitoring           | Unimplemented |
| Database migrations                  | Implemented   |

## Architecture

The codebase uses a small layered structure rather than a full framework-specific architecture:

```text
HTTP request
	|
Gin router and authentication middleware
	|
Handlers
	|
Services
	|
SQLC-generated repository
	|
PostgreSQL
```

- `api/` creates the Gin server, registers routes, and validates Bearer tokens.
- `internal/handler/` translates HTTP requests and responses.
- `internal/service/` contains application operations.
- `database/query/` contains SQL queries and `database/repository/` contains SQLC-generated access code.
- `database/migration/` contains PostgreSQL schema migrations.
- `config/` loads runtime configuration and `token/` handles session tokens.

## API

The following routes are registered in `api/routes.go`. Routes under the authenticated group require `Authorization: Bearer <token>`.

| Method | Endpoint                        | Description                    | Auth |
| ------ | ------------------------------- | ------------------------------ | ---- |
| POST   | `/auth/login/local`             | Local login                    | No   |
| POST   | `/auth/login/account`           | Account login                  | No   |
| POST   | `/user/register/local`          | Register a local user          | No   |
| POST   | `/user/register/account`        | Register an account user       | No   |
| GET    | `/auth/session`                 | Read the authenticated session | Yes  |
| GET    | `/category/list`                | List clothes categories        | Yes  |
| GET    | `/service/list`                 | List service types             | Yes  |
| POST   | `/order/create`                 | Create an order                | Yes  |
| GET    | `/order/list`                   | List orders                    | Yes  |
| GET    | `/order/detail/:order_id`       | Get order details              | Yes  |
| PUT    | `/order/update/:order_id`       | Update an order                | Yes  |
| PUT    | `/order/update/:order_id/items` | Update order items             | Yes  |
| PUT    | `/order/status/update`          | Update order status            | Yes  |
| DELETE | `/order/delete`                 | Delete an order                | Yes  |
| GET    | `/customer/list`                | List customers                 | Yes  |
| GET    | `/customer/detail`              | Get customer details           | Yes  |

This is a route overview, not a complete request and response specification. No Swagger or OpenAPI documentation was found in the current repository.

## Database

The application connects to PostgreSQL using `database/sql` with the `pgx` driver. The current migrations define users, customers, clothes categories, service types, orders, order items, and notification logs, with later migrations adding customer contact fields.

Migrations are stored in `database/migration/`. The application does not apply migrations automatically during startup.

## Configuration

Create a `.env` file in the repository root, or provide the same variables through the environment when `.env` is absent:

```dotenv
APP_NAME=jahitin
APP_VERSION=development
APP_PORT=8001
DEBUG=true
DB_HOST=localhost
DB_PORT=5432
DB_USER=<database-user>
DB_PASSWORD=<database-password>
DB_NAME=jahitin_db
TOKEN_SECRET_KEY=<token-secret>
```

Do not commit real credentials or token secrets. The repository's Docker Compose setup expects PostgreSQL and application settings from `.env`.

## Getting Started

Prerequisites: Go 1.25 and a running PostgreSQL database configured with the variables above.

1. Clone the repository and enter its directory.
2. Create `.env` using the configuration example above.
3. Apply the SQL migrations in `database/migration/` to the configured database.
4. Start the API directly or through the Makefile:

```bash
go mod download
go run main.go
```

Equivalent Makefile commands are:

```bash
make server
# or
make api
```

The server binds to `0.0.0.0:${APP_PORT}`.

To build and run the API and PostgreSQL services with Docker Compose:

```bash
make run
# or
docker compose up -d --build
```

Stop the Compose services with `make down` or `docker compose down`. The Compose configuration expects external Docker networks named `jahitin-network` and `main-network`.

## Known Limitations

- The project is still in active development and API contracts may change.
- Database migrations are not executed automatically by the application.
- The `cmd/seed` package currently does not compile against the generated repository types.

## Roadmap

The repository does not contain an explicit roadmap. Based on the current verified gaps, the next maintenance priorities are:

- [ ] Repair and verify the database seed command.
- [ ] Verify and document the API contract with the Flutter client.

## Related Repository

Mobile app: [safelyn_mobile](https://github.com/mirza27/safelyn_mobile)
