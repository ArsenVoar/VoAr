# VoAr

VoAr is a backend-focused web application written in Go.

The project was built from scratch as a personal backend engineering project with a strong focus on architecture, request lifecycle management, authentication, PostgreSQL integration, and layered application design.

The primary goal of the project is not commercial deployment, but practical backend development experience, infrastructure understanding, and portfolio presentation.

---

# Features

* User registration and authentication
* Session-based authorization
* User profile system
* Article creation and viewing
* REST API endpoints
* PostgreSQL integration
* Layered backend architecture
* Middleware-based authorization
* Context propagation through application layers
* Docker and docker-compose support
* HTML template rendering
* Google OAuth integration (partially implemented)

---

# Tech Stack

* Go
* PostgreSQL
* Docker
* gorilla/mux
* gorilla/sessions
* bcrypt
* HTML templates
* PostgreSQL driver (lib/pq)

---

# Architecture

The application uses a layered backend architecture:

Handler → Service → Repository → Database

## Handler Layer

Responsible for:

* HTTP request handling
* Input extraction
* Response rendering
* HTTP status codes
* Context extraction from requests

## Service Layer

Responsible for:

* Business logic
* Validation
* Authorization logic
* Error mapping
* Request flow orchestration

## Repository Layer

Responsible for:

* Database access
* SQL queries
* PostgreSQL interaction
* QueryContext / ExecContext usage

## Infrastructure

The project includes:

* Dependency injection through main.go
* Session middleware
* Environment-based configuration
* Context propagation across layers
* Request lifecycle control

---

# Context Propagation

VoAr uses Go context propagation across backend layers.

HTTP request context flows through:

Handler → Service → Repository → Database

The repository layer uses:

* QueryContext
* QueryRowContext
* ExecContext

This allows:

* request cancellation
* timeout support
* request-scoped execution flow
* proper backend resource management

---

# Project Structure

```text
cmd/
    voar/
        main.go

internal/
    database/
    handler/
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

VoAr uses session-based authentication.

Features:

* cookie sessions
* middleware authorization
* protected profile routes
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
* layered repository access

---

# Backend Concepts Implemented

* Layered architecture
* Dependency injection
* Context propagation
* Request lifecycle handling
* Middleware authorization
* Repository pattern
* Session management
* REST API basics
* PostgreSQL integration
* Docker containerization
* Error handling
* Pagination

---

# Project Goal

The project was created primarily as a backend engineering project and portfolio application.

The focus of development is:

* backend architecture
* infrastructure understanding
* request lifecycle management
* practical Go backend engineering
* production-oriented thinking

---

# Interview Summary

VoAr is a Go backend application using layered architecture with handlers, services, repositories, and PostgreSQL.

The project implements authentication, session management, middleware authorization, REST API endpoints, Docker support, and context propagation through all backend layers.

The application uses Go contexts with QueryContext/ExecContext to support request lifecycle management and cancellable database operations.

The main goal of the project was to practice practical backend engineering concepts rather than build a feature-heavy commercial application.
