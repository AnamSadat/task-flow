# Task Flow using Golang Standar Library

## Structure Project

```
internal/
├── config/              # Database config
│   └── database.go
├── handler/             # HTTP handlers ONLY
│   ├── auth.go          # Login, Refresh, Logout handlers
│   └── task.go          # Task handlers
├── service/             # Business logic
│   └── auth/
│       └── auth.go      # Auth service (Login, Refresh, Logout logic)
├── repository/          # Database interfaces
│   ├── user.go          # UserRepo interface
│   └── refresh_token.go # RefreshTokenRepo interface
├── model/               # Entities
│   └── user.go          # User struct
├── pkg/                 # Shared utilities
│   └── jwt/
│       └── jwt.go       # JWT sign/verify
├── middleware/
│   └── logger.go
├── httpx/
│   └── response.go
└── router/
    └── router.go        # Route registration
```
