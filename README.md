# VoAr

VoAr is a backend web application written in Go, focused on scalable application architecture, infrastructure design, caching, observability, authentication, and runtime coordination.

### Key Technologies

Go • PostgreSQL • Redis • Docker • Database Migrations • REST API • Context Propagation • Structured Logging

---

The project combines application functionality with backend infrastructure concerns, including:

* layered architecture
* middleware orchestration
* request lifecycle management
* context propagation
* observability
* graceful shutdown
* transaction management
* database migrations
* caching strategies
* runtime coordination

---

# Features

## Architecture Overview

```text
Request
→ Middleware
→ Handler
→ Service
→ Cache
→ Repository
→ PostgreSQL
```

Redis is used through a Cache-Aside pattern with automatic fallback to a NoOp implementation when Redis is unavailable.

## Application

* User registration and authentication
* Session-based authorization
* User profiles
* Article publishing system
* Comment system
* Notification system
* Server-side HTML rendering
* REST API endpoints
* Google OAuth integration (partial)
* Redis-backed post caching
* Automatic cache fallback

## Backend Infrastructure

* Request-scoped context propagation
* Structured request logging
* Request correlation with request IDs
* Timeout middleware
* Graceful shutdown
* Transaction lifecycle logging
* Context-aware database operations
* Concurrent request handling
* Cache hit/miss/error monitoring
* Docker support
* Database migrations

---

# Tech Stack

## Backend

* Go
* PostgreSQL
* Redis
* REST API

## Infrastructure

* Docker
* Docker Compose
* golang-migrate

## Libraries

* gorilla/mux
* gorilla/sessions
* bcrypt
* lib/pq
* github.com/redis/go-redis/v9

---

# Runtime Architecture

Request lifecycle flow:

```text
Request
→ RequestIDMiddleware
→ TimeoutMiddleware
→ LoggingMiddleware
→ Handler
→ Service
→ Cache (Redis or NoOp)
→ Repository
→ PostgreSQL
```

Application lifecycle flow:

```text
Startup
→ Database initialization
→ Redis connection or NoOp fallback
→ Runtime execution
→ Concurrent request handling
→ Graceful shutdown signal
→ Request draining
→ Cleanup
→ Shutdown
```

---

# Layered Architecture

```text
Handler → Service → Cache → Repository → Database
```

## Handler Layer

Responsible for:

* HTTP request handling
* request validation
* response generation
* session management
* HTTP status management

## Service Layer

Responsible for:

* business logic
* authorization
* validation
* transaction coordination
* notification orchestration
* caching coordination
* runtime orchestration

## Cache Layer

Responsible for:

* caching frequently accessed posts
* Redis integration
* fallback to NoOp cache
* cache lifecycle logging
* TTL management
* serialization and deserialization

## Repository Layer

Responsible for:

* database access
* SQL execution
* context-aware queries
* transaction participation through DBTX

---

# Database

## Schema

The application currently consists of the following entities:

* users
* articles
* comments
* notifications

Relationships:

```text
User
 ├─ Articles
 ├─ Comments
 └─ Notifications

Article
 ├─ Author
 └─ Comments

Comment
 ├─ Author
 └─ Article
```

Database changes are managed through versioned migrations using golang-migrate.

---

# Runtime Infrastructure

## Context Propagation

Context flows through the entire request lifecycle:

```text
Handler → Service → Cache → Repository → Database
```

Used for:

* request cancellation
* timeout propagation
* graceful shutdown
* request metadata propagation

The project uses:

* QueryContext
* QueryRowContext
* ExecContext

---

## Structured Logging

The application includes centralized structured logging with:

* request IDs
* HTTP method and path
* response status
* execution duration
* transaction lifecycle events
* cache events
* startup and shutdown events

Logs are written to:

* stdout
* logs/app.log

Example:

```text
REQUEST request_id=1780602449045974300 method=GET path=/show/1 status=200 duration=4.85ms
INFO request_id=1780602449045974300 cache_hit key=post:1
INFO request_id=1780602029563122200 cache_miss key=post:1
INFO request_id=1780602029563122200 cache_set key=post:1
```

---

## Timeout Middleware

The timeout middleware creates bounded request execution using Go contexts.

This protects the application from:

* uncontrolled execution
* hanging requests
* resource exhaustion

---

## Graceful Shutdown

VoAr supports coordinated runtime shutdown with:

* SIGINT / SIGTERM handling
* request draining
* bounded shutdown timeout
* context cancellation propagation

---

# Testing

Current test coverage includes:

* user service authentication
* middleware behavior
* request ID propagation
* timeout injection
* context cancellation

Run all tests:

```bash
go test ./...
```

---

# Observability

The application provides visibility into:

* request lifecycle
* timeout cancellations
* startup and shutdown events
* request duration
* transaction lifecycle
* cache usage
* middleware execution flow

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

# Run Locally

## 1. Create .env

```env
SESSION_SECRET=your_secret

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=voar
```

## 2. Apply migrations

```bash
migrate -path db/migrations \
-database "postgres://postgres:password@localhost:5432/voar?sslmode=disable" up
```

## 3. Start PostgreSQL and Redis

Ensure PostgreSQL and Redis are running locally.

## 4. Run application

```bash
go run ./cmd/voar
```

Application runs on:

```text
http://localhost:8080
```

---

# Docker

Build and run:

```bash
docker compose up --build
```

---

# Backend Engineering Concepts

* Layered architecture
* Dependency injection
* Repository pattern
* Middleware architecture
* Request lifecycle management
* Context propagation
* Structured logging
* Runtime observability
* Timeout middleware
* Graceful shutdown
* Transaction handling
* Redis cache-aside pattern
* Session management
* PostgreSQL integration
* Database migrations
* Docker containerization
* Entity relationships
* Service orchestration
* Notification workflows

```
```
