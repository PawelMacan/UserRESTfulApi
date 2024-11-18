# UserRESTfulApi

A high-performance RESTful API implementation in Go following clean architecture principles, with comprehensive testing and performance optimization.

## Project Structure

```
.
├── cmd/
│   └── server/           # Main application entry point
├── internal/
│   ├── domain/          # Business logic and entities
│   ├── handlers/        # HTTP handlers
│   ├── router.go        # HTTP router setup
│   ├── repository/      # Data access layer
│   └── service/         # Business logic implementation
├── pkg/                 # Shared packages
├── api/                 # API documentation
└── tests/              # Test suites
    ├── load/           # Load testing scripts
    └── integration/    # Integration tests
```

## Prerequisites

- Go 1.22.1 or higher
- PostgreSQL 14+
- Docker (optional)
- k6 for load testing
- Node.js (for running certain tests)

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db_name
SERVER_PORT=8080
```

## Running the Application

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Start PostgreSQL:
```bash
# If using Docker
docker run --name postgres -e POSTGRES_PASSWORD=your_password -p 5432:5432 -d postgres
```

3. Run the application:
```bash
go run cmd/server/main.go
```

### Using Docker

1. Build the Docker image:
```bash
docker build -t user-restful-api .
```

2. Run the container:
```bash
docker run -p 8080:8080 --env-file .env user-restful-api
```

### Docker Compose (Recommended for Development)

```bash
docker-compose up -d
```

## Testing

### Unit Tests

Run all unit tests:
```bash
go test ./... -v
```

### Integration Tests

Run integration tests:
```bash
go test ./tests/integration/... -v
```

### Load Testing

The application includes comprehensive load testing using k6:

1. Install k6:
```bash
brew install k6
```

2. Run the basic load test:
```bash
k6 run tests/load/load_test.js
```

3. Run extended load test (5-minute duration with 500 VUs):
```bash
k6 run tests/load/load_test.js --vus 500 --duration 5m
```

#### Load Test Profiles

- **Smoke Test**: Quick test with minimal load
  ```bash
  k6 run --vus 1 --duration 30s tests/load/load_test.js
  ```

- **Stress Test**: High load for extended period
  ```bash
  k6 run --vus 500 --duration 5m tests/load/load_test.js
  ```

- **Spike Test**: Sudden spike in users
  ```bash
  k6 run --vus 1000 --duration 1m tests/load/load_test.js
  ```

## API Endpoints

- `POST /api/users` - Create a new user
- `GET /api/users/{id}` - Get user by ID
- `GET /api/users` - List all users
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

## Performance Metrics

Current performance baseline (with 500 concurrent users):
- Request Rate: ~1,000 requests/second
- Average Response Time: ~125ms
- Success Rate: 99.99%
- P95 Response Time: ~560ms

## Monitoring

The application exposes metrics endpoints for monitoring:
- `/metrics` - Prometheus metrics
- `/health` - Health check endpoint

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
