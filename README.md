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
* Docker support

---

# Tech Stack

## Backend

* Go
* PostgreSQL
* REST API

## Infrastructure

* Docker
* docker-compose

## Libraries

* gorilla/mux
* gorilla/sessions
* bcrypt
* lib/pq

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
→ Repository
→ PostgreSQL
```

Application lifecycle flow:

```text
Startup
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
Handler → Service → Repository → Database
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
* runtime orchestration

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
Handler → Service → Repository → Database
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
* startup/shutdown events

Logs are written to:

* stdout
* logs/app.log

Example:

```text
REQUEST request_id=1779452438945297900 method=GET path=/auth status=200 duration=0s
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
* middleware execution flow

---

# Project Structure

```text
cmd/
    voar/

internal/
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
```

## 2. Start PostgreSQL

Make sure PostgreSQL is running locally.

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
* lifecycle-aware infrastructure
* production-oriented backend engineering
