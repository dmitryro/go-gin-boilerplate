# Go & Gin Boilerplate Project

Boilerplate project using Go and Gin – a clean, simple, and lightweight starting point leveraging PostgreSQL, JWT-based RBAC, Swagger, and Dockerized workflows.

[![Go](https://img.shields.io/badge/Go-1.23-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green)](./LICENSE)
[![Status](https://img.shields.io/badge/Status-Stable-yellow)](https://github.com/)

**Version**: 1.0.0  
**License**: [MIT License](./LICENSE)

---

## 📑 Contents
- [🚀 Overview](#-overview)
- [🛠 Setup Instructions](#-setup-instructions)
  - [1. Clone Project](#1-clone-project)
  - [2. Start Project](#2-start-project)
    - [2.1 Standalone (Go)](#21-standalone-go)
    - [2.2 Docker & Docker Compose](#22-docker--docker-compose)
  - [3. Swagger Initialization](#3-swagger-initialization)
  - [4. .env Setup](#4-env-setup)
- [📁 Project Structure](#-project-structure)
- [🔐 RBAC System](#-rbac-system)
- [🧩 Middleware](#-middleware)
- [🧪 Handlers](#-handlers)
- [🌐 Routes](#-routes)
- [🧰 Services](#-services)
- [📦 Models](#-models)
- [🏁 Main Entry](#-main-entry)
- [📦 Go Modules](#-go-modules)
- [✅ Unit Testing](#-unit-testing)

---

## 🚀 Overview

This boilerplate is perfect for building a RESTful API in Go using the Gin framework. It includes:
- PostgreSQL DB integration
- JWT-based authentication & role-based access control (RBAC)
- Swagger API documentation
- Dockerfile and Docker Compose support
- Modular structure with services, handlers, and middleware

---

## 🛠 Setup Instructions

### 1. Clone Project

```bash
git clone git@github.com:dmitryro/go-gin-boilerplate.git gin-project
cd gin-project
```

### 2. Start Project

#### 2.1 Standalone (Go)

Make sure Go 1.23+ is installed.

```bash
export GO_PATH=$HOME/go
export GO_ROOT=/usr/local/go
export PATH=$PATH:$GO_ROOT/bin:$GO_PATH/bin
go mod tidy
go run src/cmd/api/main.go
```

#### 2.2 Docker & Docker Compose

Ensure Docker & Docker Compose are installed.

```bash
cp .env_template .env  # Set your values in .env
docker-compose up --build
```

### 3. Swagger Initialization

To regenerate Swagger docs after changes:

```bash
swag init -g ./src/cmd/api/main.go -o ./docs --parseDependency --parseInternal
```

#### 3.1 Dockerfile Swagger Support

Ensure this exists in Dockerfile:

```Dockerfile
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g src/cmd/api/main.go
```

### 4. .env Setup

Copy template and customize:

```bash
cp .env_template .env
```

Edit `.env`:

```env
PG_DATABASE=yourdb
PG_PORT=5432
PG_HOST=127.0.0.1
APP_PORT=8081
PG_USER=yourusername
PG_PASSWORD=yourpassword
JWT_KEY=your-super-secure-key
SWAGGER_YAML_DIR=./docs/swagger.yaml
SWAGGER_JSON_DIR=./docs/swagger.json
```

---

## 📁 Project Structure

```
.
├── Dockerfile
├── docker-compose.yaml
├── .env_template
├── src/
│   ├── cmd/api/main.go
│   └── internal/
│       ├── handlers/
│       │   ├── login_handler.go
│       │   ├── register_handler.go
│       │   ├── role_handler.go
│       │   └── user_handler.go
│       ├── middleware/
│       │   └── jwt_middleware.go
│       ├── models/
│       ├── routes/routes.go
│       └── services/
│           ├── login_service.go
│           ├── register_service.go
│           ├── role_service.go
│           └── user_service.go
├── docs/ (swagger)
├── sql/init_db.sql
├── tests/
└── README.md
```

---

## 🔐 RBAC System

Example roles stored in DB:

```json
[
  {
    "name": "admin",
    "permissions": ["create", "read", "update", "delete", "superuser"]
  },
  {
    "name": "guest",
    "permissions": ["read"]
  }
]
```

To add roles, insert into the `roles` table with proper permission arrays.

---

## 🧩 Middleware

Located in `src/internal/middleware/jwt_middleware.go`.

- Validates JWT tokens
- Applies role/permission checks

To add new middleware:
- Create new file in middleware/
- Register in `main.go` or `routes.go`

---

## 🧪 Handlers

Each handler corresponds to a feature:

| Handler | Description |
|---------|-------------|
| `login_handler.go` | Handles authentication |
| `register_handler.go` | User registration |
| `role_handler.go` | Role CRUD |
| `user_handler.go` | User management, search, password change |

Swagger annotations (`@Summary`, `@Param`, etc.) are included. Modify in handler comments and regenerate using `swag init`.

---

## 🌐 Routes

Routes are defined in `routes.go`.

To add a new route:
- Define in handler
- Register in `routes.go`
- Annotate for Swagger

---

## 🧰 Services

Service layer handles business logic and DB operations.

- `user_service.go`: users, roles, passwords
- `role_service.go`: role data
- `login_service.go`: auth token
- `register_service.go`: signup logic

Create new services by following similar structure and injecting via handler constructors.

---

## 📦 Models

Located in `src/internal/models`. Represents DB schema and request/response types.

Includes:
- `user.go`
- `role.go`
- `login.go`
- `login_request.go`
- `token_response.go`
- `error_response.go`

---

## 🏁 Main Entry

`src/cmd/api/main.go` handles:
- Middleware and handler registration
- Swagger setup
- Database migration and role constraints
- API group and route wiring

To refresh Swagger:
```bash
swag init -g ./src/cmd/api/main.go -o ./docs --parseDependency --parseInternal
```

---

## 📦 Go Modules

Project uses Go Modules (`go.mod`, `go.sum`).

Best practices:
- Use `go mod tidy` after adding packages
- Avoid committing compiled binaries
- Integrate with Docker `go install` for tools

---

## ✅ Unit Testing

Tests go under `/tests`. To run tests:

```bash
go test ./...
```

To run tests inside Docker:
- Add `RUN go test ./...` in Dockerfile or
- Use a separate test service in docker-compose

Best practices:
- Cover each handler and service
- Mock DB where necessary

