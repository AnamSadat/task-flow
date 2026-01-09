# Task Flow API

REST API untuk task management menggunakan Go standard library.

## Tech Stack

- Go 1.25+
- MySQL
- JWT Authentication

## Getting Started

### Prerequisites

- Go 1.25+
- MySQL
- Air (hot reload) - optional

### Installation

```bash
# Clone repository
git clone <repo-url>
cd task-flow

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env sesuai konfigurasi database kamu
```

### Setup Database

```bash
# Install migrate CLI
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "mysql://user:pass@tcp(localhost:3306)/task_flow" up
```

### Run Server

```bash
# Development (dengan hot reload)
air

# Atau tanpa hot reload
go run ./cmd/api/main.go
```

Server berjalan di `http://localhost:3000` (atau sesuai `APP_ADDR` di .env)

## API Endpoints

### Auth

```
POST /auth/register    - Register user baru
POST /auth/login       - Login, dapat access & refresh token
POST /auth/refresh     - Refresh access token
POST /auth/logout      - Logout, revoke refresh token
```

### Users (Protected)

```
GET /users/me          - Get current user info
```

### Tasks

```
GET  /tasks            - Get all tasks
POST /tasks            - Create new task
```

### Health Check

```
GET /health            - Server health check
```

## Project Structure

```
task-flow/
├── cmd/
│   └── api/
│       └── main.go           # Entry point
├── internal/
│   ├── config/               # Configuration
│   │   ├── config.go         # App config
│   │   └── database.go       # Database connection
│   ├── handler/              # HTTP handlers
│   │   ├── auth.go
│   │   ├── task.go
│   │   └── users.go
│   ├── service/              # Business logic
│   │   ├── auth/
│   │   │   └── auth.go
│   │   └── task.go
│   ├── repository/           # Data access layer
│   │   ├── user.go           # Interface
│   │   ├── task.go           # Interface
│   │   ├── refresh_token.go  # Interface
│   │   └── mysql/            # MySQL implementation
│   │       ├── user.go
│   │       ├── task.go
│   │       └── refresh_token.go
│   ├── model/                # Data models
│   │   ├── user.go
│   │   └── task.go
│   ├── middleware/           # HTTP middleware
│   │   ├── auth.go
│   │   └── logger.go
│   ├── pkg/                  # Shared packages
│   │   ├── jwt/
│   │   └── cookie/
│   ├── httpx/                # HTTP helpers
│   │   └── response.go
│   ├── router/               # Route registration
│   │   └── router.go
│   └── utils/                # Utilities
│       └── generate_id.go
├── migrations/               # Database migrations
├── docs/                     # Documentation
└── tests/                    # Integration tests
```

## Environment Variables

```env
APP_ADDR=:3000
DATABASE_URL=mysql://user:pass@localhost:3306/task_flow
JWT_SECRET=your-secret-key
ACCESS_TTL=15m
REFRESH_TTL=168h
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose
go test -v ./...

# Run specific package
go test -v ./internal/service/auth/...

# Run with coverage
go test -cover ./...
```

## Development

### Hot Reload dengan Air

```bash
# Install air
go install github.com/air-verse/air@latest

# Run
air
```

### Database Migration

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq <migration_name>

# Run migrations
migrate -path migrations -database "mysql://..." up

# Rollback
migrate -path migrations -database "mysql://..." down 1
```

## Architecture

```
Request → Router → Middleware → Handler → Service → Repository → Database
                                              ↓
                                           Model
```

- **Handler**: Parse HTTP request, call service, format response
- **Service**: Business logic, validation
- **Repository**: Database operations
- **Model**: Data structures
