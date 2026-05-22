````markdown
# VoAr

VoAr is a production-oriented backend web application written in Go.

The project was built from scratch as a personal backend engineering project focused on:

- layered backend architecture
- runtime infrastructure
- request lifecycle management
- observability
- context propagation
- middleware coordination
- graceful shutdown
- production-oriented backend engineering concepts

The primary goal of the project is practical backend engineering experience, runtime systems understanding, and infrastructure-oriented portfolio development.

---

# Features

## Application Features

- User registration and authentication
- Session-based authorization
- User profile system
- Article creation and viewing
- REST API endpoints
- HTML template rendering
- Google OAuth integration (partially implemented)

## Backend Infrastructure Features

- Layered backend architecture
- Middleware-based runtime coordination
- Request-scoped context propagation
- Structured request logging
- Request correlation with request IDs
- Timeout middleware with bounded execution
- Graceful shutdown lifecycle coordination
- Runtime lifecycle observability
- Concurrent request execution handling
- PostgreSQL integration with cancellable queries
- Docker and docker-compose support

---

# Tech Stack

## Backend

- Go
- PostgreSQL
- REST API

## Infrastructure

- Docker
- docker-compose

## Libraries

- gorilla/mux
- gorilla/sessions
- bcrypt
- lib/pq

## Runtime Infrastructure

- Context propagation
- Middleware architecture
- Structured logging
- Request correlation
- Timeout middleware
- Graceful shutdown

---

# Runtime Architecture

VoAr uses a layered backend architecture combined with middleware-driven runtime infrastructure.

## Request Runtime Flow

```text
request
→ RequestIDMiddleware
→ TimeoutMiddleware
→ LoggingMiddleware
→ Handler
→ Service
→ Repository
→ PostgreSQL
````

## Application Lifecycle Flow

```text
application startup
→ runtime execution
→ request handling
→ concurrent request processing
→ graceful shutdown signal
→ request draining
→ cleanup
→ runtime termination
```

---

# Layered Architecture

The application uses a classic backend layering model:

```text
Handler → Service → Repository → Database
```

## Handler Layer

Responsible for:

* HTTP request handling
* Input extraction
* Response generation
* HTTP status management
* Request context extraction
* Request lifecycle entrypoints

## Service Layer

Responsible for:

* Business logic
* Validation
* Authorization rules
* Error coordination
* Request flow orchestration

## Repository Layer

Responsible for:

* Database access
* SQL execution
* PostgreSQL interaction
* QueryContext / ExecContext usage
* Context-aware database operations

## Middleware Infrastructure Layer

Responsible for:

* Request correlation
* Structured logging
* Runtime instrumentation
* Timeout propagation
* Request lifecycle coordination
* Runtime observability
* Execution boundaries

---

# Runtime Infrastructure

VoAr includes production-style runtime infrastructure patterns.

## Request Correlation

Each request receives a unique request ID through middleware.

Request IDs allow:

* correlated request logs
* execution flow reconstruction
* runtime observability
* operational debugging

## Structured Logging

The application implements centralized structured logging.

Structured logs include:

* request ID
* HTTP method
* request path
* response status
* execution duration
* runtime lifecycle events

Example:

```text
REQUEST request_id=1779452438945297900 method=GET path=/slow status=408 duration=5.0002472s
```

## Timeout Middleware

VoAr implements timeout-based bounded execution protection.

The timeout middleware:

* creates derived timeout contexts
* propagates deadlines across layers
* supports cancellation-aware execution
* protects runtime resources from uncontrolled execution

This improves:

* runtime predictability
* infrastructure stability
* resource management

## Graceful Shutdown

The application supports graceful runtime termination.

Shutdown lifecycle includes:

* SIGINT / SIGTERM handling
* request draining
* bounded shutdown timeout
* coordinated runtime termination
* lifecycle-aware infrastructure cleanup

This prevents:

* abrupt request interruption
* uncontrolled runtime termination
* inconsistent shutdown behavior

---

# Context Propagation

VoAr uses Go contexts as centralized request lifecycle propagation infrastructure.

Request context flows through:

```text
Handler → Service → Repository → Database
```

The propagated context carries:

* request-scoped metadata
* request IDs
* cancellation signals
* timeout deadlines
* lifecycle coordination information

The repository layer uses:

* QueryContext
* QueryRowContext
* ExecContext

This enables:

* request cancellation
* timeout propagation
* graceful shutdown coordination
* bounded execution
* context-aware database operations

---

# Observability

VoAr includes middleware-driven runtime observability infrastructure.

The application provides visibility into:

* request execution lifecycle
* timeout cancellations
* runtime startup/shutdown events
* request duration metrics
* correlated execution flow
* infrastructure lifecycle transitions

The backend runtime is observable rather than operating as a black box system.

---

# Concurrency and Runtime Coordination

The application runtime uses concurrent lifecycle coordination.

The HTTP server executes in a dedicated goroutine while the main goroutine coordinates:

* application lifecycle management
* graceful shutdown
* signal handling
* runtime termination flow

This separation improves:

* runtime coordination
* lifecycle management
* infrastructure predictability

---

# Project Structure

```text
cmd/
    voar/
        main.go

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

## 1. Clone the repository

```bash
git clone <repository-url>
```

## 2. Create .env file

Example:

```env
SESSION_SECRET=your_secret
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_CALLBACK_URL=

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=voar
```

## 3. Run PostgreSQL

Make sure PostgreSQL is running locally.

## 4. Start the application

```bash
go run ./cmd/voar
```

Application runs on:

```text
http://localhost:8080
```

---

# Docker

The project includes:

* Dockerfile
* docker-compose.yml

Run with Docker:

```bash
docker-compose up --build
```

Services:

* app
* postgres

The application communicates with PostgreSQL through the Docker network using the postgres service name.

---

# Authentication

VoAr uses session-based authentication.

Features include:

* cookie sessions
* middleware authorization
* protected routes
* login/logout flow
* bcrypt password hashing

---

# Database

PostgreSQL is used as the primary database.

The project uses:

* SQL queries
* QueryRowContext
* QueryContext
* ExecContext
* context-aware database execution

---

# Backend Engineering Concepts Implemented

* Layered architecture
* Dependency injection
* Middleware architecture
* Request lifecycle management
* Context propagation
* Request correlation
* Structured logging
* Runtime observability
* Timeout middleware
* Graceful shutdown
* Bounded execution
* Runtime lifecycle coordination
* Session management
* Repository pattern
* REST API fundamentals
* PostgreSQL integration
* Docker containerization
* Error handling
* Pagination
* Concurrent runtime coordination

---

# Project Goal

VoAr was created primarily as a backend engineering and runtime infrastructure learning project.

The focus of development is:

* backend architecture
* runtime systems understanding
* middleware coordination
* lifecycle-aware infrastructure
* production-oriented backend engineering
* observability
* request lifecycle management
* infrastructure reliability concepts

---

# Interview Summary

VoAr is a production-oriented Go backend application implementing layered architecture, middleware-driven runtime infrastructure, PostgreSQL integration, request lifecycle coordination, structured logging, timeout middleware, and graceful shutdown support.

The project uses request-scoped context propagation across handlers, services, repositories, and database operations to support cancellation propagation, bounded execution, runtime observability, and coordinated lifecycle management.

The backend runtime includes request correlation, structured telemetry, graceful runtime termination, concurrent request execution handling, and middleware-based runtime instrumentation.

```
```
