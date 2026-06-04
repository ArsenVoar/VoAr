# VoAr

VoAr is a production-oriented backend web application written in Go.

The project focuses on backend engineering concepts including:

* layered architecture
* runtime infrastructure
* middleware coordination
* request lifecycle management
* context propagation
* observability
* graceful shutdown
* transaction handling
* concurrent runtime coordination
* caching with Redis or NoOp fallback

---

# Features

## Application

* User registration and authentication
* Session-based authorization
* User profiles
* Article and post system
* REST API endpoints
* HTML template rendering
* Google OAuth integration (partial)
* Caching of frequently accessed posts using Redis
* Automatic fallback to NoOp cache when Redis is unavailable

## Backend Infrastructure

* Request-scoped context propagation
* Structured request logging
* Request correlation with request IDs
* Timeout middleware
* Graceful shutdown
* Transaction lifecycle logging
* Context-aware database queries
* Concurrent request handling
* Middleware-driven runtime coordination
* Redis caching for improved read performance
* Cache hit/miss/error logging
* Docker support

---

# Tech Stack

## Backend

* Go
* PostgreSQL
* Redis for caching
* REST API

## Infrastructure

* Docker
* docker-compose

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
→ Runtime execution
→ Redis connection or NoOp fallback
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
* HTTP status management

## Service Layer

Responsible for:

* business logic
* authorization
* validation
* transaction coordination
* caching coordination
* runtime orchestration

## Cache Layer

Responsible for:

* caching frequently accessed data (posts)
* Redis integration
* fallback to NoOp cache if Redis is unavailable
* logging cache hits, misses, sets, deletes, errors
* TTL management for cache entries
* JSON serialization/deserialization for cached data

## Repository Layer

Responsible for:

* database access
* SQL execution
* context-aware queries

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
* HTTP method/path
* response status
* execution duration
* transaction lifecycle events
* cache events: hit, miss, set, delete, error
* startup/shutdown events

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

This protects the runtime from:

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

The project includes tests for:

* login service
* middleware behavior
* request ID propagation
* timeout injection
* context cancellation
* cache hits, misses, and fallback behavior

Run all tests:

```bash
go test ./...
```

---

# Observability

The backend runtime is observable rather than operating as a black box system.

The project provides visibility into:

* request lifecycle
* timeout cancellations
* runtime startup/shutdown
* request duration
* transaction lifecycle
* cache usage: hits, misses, sets, deletes, errors
* middleware execution flow

---

# Project Structure

```text
cmd/
    voar/

internal/
    cache/
    contextkeys/
    database/
    handler/
    logger/
    middleware/
    models/
    repository/
    service/

pkg/
    google/

web/
    templates/
    css/
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
REDIS_ADDR=localhost:6379
CACHE_TTL=5m
```

## 2. Start PostgreSQL and Redis

Make sure PostgreSQL and Redis are running locally.

## 3. Run application

```bash
go run ./cmd/voar
```

Application runs on:

```text
http://localhost:8080
```

---

# Docker

Run with Docker:

```bash
docker-compose up --build
```

---

# Backend Engineering Concepts

* Layered architecture
* Dependency injection
* Middleware architecture
* Request lifecycle management
* Context propagation
* Structured logging
* Runtime observability
* Timeout middleware
* Graceful shutdown
* Transaction handling
* Redis caching
* Concurrent runtime coordination
* Repository pattern
* Session management
* PostgreSQL integration
* Docker containerization

---

# Project Goal

VoAr was created as a backend engineering and runtime infrastructure learning project.

The main focus areas are:

* backend architecture
* runtime systems
* middleware coordination
* observability
* caching strategies with Redis and NoOp fallback
* lifecycle-aware infrastructure
* production-oriented backend engineering
