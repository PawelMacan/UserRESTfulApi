# UserRESTfulApi

A RESTful API implementation in Go following clean architecture principles.

## Project Structure

```
.
├── cmd/
│   └── server/           
├── internal/
│   ├── domain/        # Business logic and entities
│   ├── handlers/      # HTTP handlers
│   ├────router.go   # HTTP router/
│   ├── repository/    # Data access layer
│   └── service/       # Business logic implementation
├── pkg/               # Shared packages
└── api/              # API documentation
```
