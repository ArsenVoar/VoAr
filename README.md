# VoAr

![Go](https://img.shields.io/badge/Go-1.25-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue)
![Redis](https://img.shields.io/badge/Redis-7-red)
![License](https://img.shields.io/badge/License-MIT-green)

VoAr is a backend web application written in Go that demonstrates modern backend engineering practices, including layered architecture, PostgreSQL integration, Redis Cache-Aside caching, request lifecycle management, structured logging, and graceful shutdown.

---

# Features

- User registration and authentication
- Session-based authorization
- User profiles
- Post publishing
- Comment system
- Notification system
- Server-side HTML rendering
- Redis Cache-Aside caching
- Structured logging
- Context propagation
- Graceful shutdown
- Database migrations

---

# Architecture

```text
HTTP Request
        │
        ▼
 Middleware
        │
        ▼
   Handlers
        │
        ▼
   Services
   ├── Cache
   └── Repository
            │
            ▼
      PostgreSQL
```

Redis is integrated using the Cache-Aside pattern with automatic fallback to a NoOp cache implementation when Redis is unavailable.

---

# Tech Stack

### Backend

- Go
- PostgreSQL
- Redis

### Infrastructure

- Docker
- Docker Compose
- golang-migrate

### Libraries

- gorilla/mux
- gorilla/sessions
- bcrypt
- lib/pq
- go-redis/v9

---

# Project Structure

```text
cmd/
└── voar/

db/
└── migrations/

internal/
├── cache/
├── contextkeys/
├── database/
├── handler/
├── logger/
├── middleware/
├── models/
├── repository/
└── service/

pkg/
└── google/

web/
├── css/
└── templates/
```

---

# Database

Current schema:

- users
- posts
- comments
- notifications

Relationships:

```text
User
├── Posts
├── Comments
└── Notifications

Post
└── Comments
```

Database schema is managed through versioned SQL migrations.

---

# Running Locally

## 1. Clone the repository

```bash
git clone https://github.com/ArsenVoar/VoAr.git
cd VoAr
```

## 2. Configure environment

Create a `.env` file:

```env
SESSION_SECRET=your_secret

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=voar
```

## 3. Start PostgreSQL and Redis

Either run them manually or with Docker.

## 4. Apply database migrations

```bash
migrate -path db/migrations -database "postgres://postgres:your_password@localhost:5432/voar?sslmode=disable" up
```

## 5. Run the application

```bash
go run ./cmd/voar
```

Application:

```
http://localhost:8080
```

---

# Docker

Build and start all services:

```bash
docker compose up --build
```

---

# Testing

Run tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

---


# Backend Concepts Demonstrated

- Layered architecture
- Dependency Injection
- Repository pattern
- Middleware composition
- Context propagation
- Structured logging
- Graceful shutdown
- Database transactions
- Redis Cache-Aside
- Session-based authentication
- Request lifecycle management
