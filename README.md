# UserRESTfulApi

A high-performance RESTful API implementation in Go following clean architecture principles, featuring load balancing, containerization, and comprehensive testing.

## Features

- Clean Architecture design
- Load balanced with Nginx
- Multiple application instances
- PostgreSQL database
- Docker containerization
- Comprehensive health checks
- Detailed logging and monitoring
- Input validation and error handling

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
├── scripts/            # Utility scripts
└── tests/              # Test suites
    ├── load/           # Load testing scripts
    └── integration/    # Integration tests
```

## Prerequisites

- Go 1.22.1 or higher
- Docker 20+ and docker-compose
- PostgreSQL 14+
- k6 for load testing (optional)

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=userapi
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080

# PostgreSQL Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=userapi
```

## Quick Start

1. Clone the repository:
```bash
git clone <your-repository-url>
cd UserRESTfulApi
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your desired configuration
```

3. Start the application stack:
```bash
docker-compose up -d
```

4. Verify the application is running:
```bash
curl http://localhost:8080/health
```

## Running the Application

### Using Docker Compose (Recommended)

1. Start the application stack:
```bash
docker-compose up -d
```

This will start:
- 3 application instances for load balancing
- PostgreSQL database
- Nginx reverse proxy

2. Verify the services are running:
```bash
docker-compose ps
```

3. Check the application health:
```bash
curl http://localhost:8080/health
```

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Start PostgreSQL:
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:16-alpine
```

3. Run the application:
```bash
go run cmd/server/main.go
```

## API Endpoints

### User Management
- `POST /api/users` - Create a new user
  - Required fields: name, email, password
  - Password requirements:
    * Minimum length
    * Uppercase letter
    * Special character

- `GET /api/users/{id}` - Get user by ID
- `GET /api/users` - List all users
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

### System
- `/health` - Health check endpoint
- `/metrics` - Prometheus metrics (if configured)

## Testing

### Unit Tests

Run the unit tests:
```bash
go test ./internal/... ./pkg/... -v
```

### Integration Tests

The project includes comprehensive integration tests that verify the entire API functionality in a containerized environment. The tests cover:
- User CRUD operations
- Input validation
- Error handling
- Data persistence
- API response codes

To run integration tests:

1. Make sure Docker and docker-compose are installed and running.

2. Run the integration test suite:
```bash
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```

This will:
- Start a dedicated test PostgreSQL database
- Build and run the test application container
- Execute all integration tests
- Automatically clean up containers after completion

The test environment uses separate configuration defined in `tests/integration/.env.test` to avoid conflicts with the development environment.

#### Test Configuration

The integration tests use their own environment variables defined in `tests/integration/.env.test`:
```env
# Database Configuration
DB_HOST=test-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_api_test
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8081
```

#### Adding New Integration Tests

To add new integration tests:
1. Create a new test file in `tests/integration/`
2. Use the provided test utilities and helpers
3. Follow the existing patterns for consistency
4. Ensure proper cleanup after test execution

## Load Testing

The application includes k6 load testing scripts:

1. Install k6:
```bash
brew install k6
```

2. Run load tests:
```bash
# Smoke test (1 VU, 30s)
k6 run --vus 1 --duration 30s tests/load/load_test.js

# Load test (100 VUs, 5m)
k6 run --vus 100 --duration 5m tests/load/load_test.js

# Stress test (500 VUs, 5m)
k6 run --vus 500 --duration 5m tests/load/load_test.js
```

## Performance

Current performance metrics (with load balancing):
- 3 application instances
- Nginx load balancing (least connections algorithm)
- Average response time: < 5ms for most requests
- Successful health checks across all instances

## Security Considerations

Development configuration:
- Default credentials (change in production)
- SSL mode disabled
- Debug mode enabled

Production recommendations:
- Replace default credentials
- Enable SSL/TLS
- Implement proper secret management
- Configure trusted proxies
- Enable release mode

## Monitoring

The application provides:
- Health check endpoint (`/health`)
- Detailed logging in debug mode
- Request/response timing information
- Error tracking and reporting

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
