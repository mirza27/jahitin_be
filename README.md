# Jahitin Backend

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Gin Framework](https://img.shields.io/badge/Web_Framework-Gin-008ECF?logo=gin&logoColor=white)](https://gin-gonic.com)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL_16-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![SQLC](https://img.shields.io/badge/Database_Access-SQLC-FF6B6B)](https://sqlc.dev)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://www.docker.com)
[![Mobile Client](https://img.shields.io/badge/Mobile_Client-jahitin__mobile-6200EE?logo=flutter&logoColor=white)](https://github.com/mirza27/jahitin_mobile)

REST API backend for Jahitin, a mobile application designed for home-based tailors to manage orders, customer records, and garment measurements.

This project was built as a portfolio project to demonstrate backend API development with Go, relational database modeling with PostgreSQL, and integration with a Flutter mobile client.

Companion mobile app: [github.com/mirza27/jahitin_mobile](https://github.com/mirza27/jahitin_mobile)

---

## Project Overview

Home tailors often handle orders through paper notes or chat messages, making it difficult to track garment measurements, service types (custom tailoring vs. alterations), deadlines, and order statuses. 

Jahitin Backend provides the HTTP endpoints to back this workflow, handling:
- Device-based and account authentication.
- Customer directory and contact management.
- Multi-item order creation with custom line items, measurements, notes, and pricing.
- Status tracking across the order lifecycle (pending, in progress, completed, picked up).
- Master data catalogs for garment categories and tailoring services.

---

## Architecture

The codebase follows a standard layered architecture:

```text
HTTP Request
     │
     ▼
Gin Router & Middleware (Auth, CORS)
     │
     ▼
Handlers (Request validation & HTTP response formatting)
     │
     ▼
Services (Business logic & transaction coordination)
     │
     ▼
Repository (Type-safe SQL queries generated via SQLC)
     │
     ▼
PostgreSQL Database
```

### Key Technical Decisions

- **Go & Gin**: Chosen for simple routing, low resource usage, and clean standard library integration.
- **SQLC over traditional ORM**: SQL queries are written in raw SQL (`database/query/`) and compiled into type-safe Go structs using SQLC. This avoids ORM reflection overhead while keeping full control over query design and database indexes.
- **Database Transactions**: Order creation and multi-item updates run inside explicit database transactions to ensure consistency across `orders` and `order_items` tables.
- **HMAC-SHA256 Token Auth**: Session tokens are signed and verified with HMAC-SHA256 for lightweight, stateless authentication without requiring an external OAuth provider.
- **Schema Migrations**: Database versioning is handled using `golang-migrate` SQL files in `database/migration/`.

---

## Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Gin
- **Database**: PostgreSQL 16
- **Database Driver**: `pgx/v5` via `database/sql`
- **Query Generator**: SQLC
- **Configuration**: Viper (`.env` file loader)
- **Containerization**: Docker & Docker Compose
- **Migrations**: golang-migrate

---

## Project Structure

```text
jahitin_be/
├── api/                 # Server setup, route registration, and auth middleware
├── cmd/                 # CLI entrypoints (e.g. database seed tools)
├── config/              # Viper environment variable loader
├── database/
│   ├── migration/       # Up/Down SQL migration files
│   ├── query/           # Raw SQL queries compiled by SQLC
│   └── repository/      # Generated Go database models and query methods
├── internal/
│   ├── handler/         # HTTP request handlers
│   └── service/         # Business logic layer
├── token/               # HMAC token creation and verification
├── docker-compose.yml   # Multi-container setup for local development
├── Dockerfile           # Multi-stage Docker build
├── Makefile             # Common dev shortcuts
├── sqlc.yaml            # SQLC configuration
└── main.go              # Application entrypoint
```

---

## API Endpoints

Protected endpoints require the `Authorization: Bearer <token>` header.

### Authentication

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| POST | `/auth/login/local` | No | Login using registered device identifier |
| POST | `/auth/login/account` | No | Login using email/username and password |
| POST | `/user/register/local` | No | Register a local device user |
| POST | `/user/register/account` | No | Register an account-based user |
| GET | `/auth/session` | Yes | Verify active session and return current user profile |

### Orders & Items

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| POST | `/order/create` | Yes | Create an order with customer data and item list |
| GET | `/order/list` | Yes | List orders with optional status and search filters |
| GET | `/order/detail/:order_id` | Yes | Get full order details and individual items |
| PUT | `/order/update/:order_id` | Yes | Update order name and deadline |
| PUT | `/order/update/:order_id/items` | Yes | Replace/update line items in an order |
| PUT | `/order/status/update` | Yes | Update order progress status |
| DELETE | `/order/delete` | Yes | Delete an order |

### Master Data & Customers

| Method | Endpoint | Auth | Description |
| :--- | :--- | :--- | :--- |
| GET | `/category/list` | Yes | List clothes categories (e.g., Kemeja, Celana, Gamis) |
| GET | `/service/list` | Yes | List service types (e.g., Jahit Baru, Permak) |
| GET | `/customer/list` | Yes | List tailor's customer records |
| GET | `/customer/detail` | Yes | Get customer detail by ID |

---

## Sample Request Payload

### Create Order (`POST /order/create`)

```json
{
  "name": "Pesanan Baru Bu Siti",
  "deadline": "2026-10-15T17:00:00+07:00",
  "customer": {
    "customer_name": "Bu Siti",
    "customer_phone": "085123456789"
  },
  "order_items": [
    {
      "clothes_for": "Bu Siti",
      "clothes_category_id": 2,
      "service_type_id": 1,
      "price": 150000,
      "notes": "Ukuran LD 100cm, panjang gamis 135cm",
      "is_save_customer_notes": true
    },
    {
      "clothes_for": "Anak",
      "clothes_category_id": 3,
      "custom_service_name": "Kecilkan Pinggang",
      "price": 35000,
      "notes": "Kecilkan 2 cm",
      "is_save_customer_notes": false
    }
  ]
}
```

---

## Getting Started

### Prerequisites

- Go 1.25 or newer
- PostgreSQL 16
- Docker & Docker Compose (optional)

### Local Setup

1. **Clone repository:**
   ```bash
   git clone https://github.com/mirza27/jahitin_be.git
   cd jahitin_be
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   ```
   Update `.env` with your PostgreSQL database credentials and a secure token secret key.

3. **Run database migrations:**
   ```bash
   make migrate
   ```

4. **Start API server:**
   ```bash
   go run main.go
   # or
   make server
   ```

   The server runs on `http://localhost:8001` by default.

### Docker Setup

To run both the PostgreSQL database and API service inside containers:

```bash
docker compose up -d --build
```

---

## Scope & Future Improvements

This project is built to demonstrate practical full-stack mobile backend integration. Areas that can be expanded in the future include:
- Automated integration testing with `testcontainers-go`.
- Role-based access control and multi-tenant isolation.
- PDF invoice generation and automated customer notification webhooks.
