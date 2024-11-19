#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
YELLOW='\033[1;33m'

# Error handling
set -e
trap 'cleanup' ERR INT TERM

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}Cleaning up...${NC}"
    docker-compose down 2>/dev/null || true
    exit 1
}

# Version requirements
REQUIRED_GO_VERSION="1.22.1"
REQUIRED_POSTGRES_VERSION="14"

echo -e "${YELLOW}Starting User REST API Application...${NC}"

# Check if Docker is installed and running
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check Docker version
DOCKER_VERSION=$(docker --version | cut -d" " -f3 | cut -d"." -f1)
if [ "$DOCKER_VERSION" -lt "20" ]; then
    echo -e "${RED}Docker version 20 or higher is required${NC}"
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}docker-compose is not installed. Please install docker-compose first.${NC}"
    exit 1
fi

# Check if Docker daemon is running
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}Docker daemon is not running. Please start Docker first.${NC}"
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env file...${NC}"
    cat > .env << EOL
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=userapi
DB_SSLMODE=disable
SERVER_PORT=8080
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=userapi
EOL
    echo -e "${GREEN}Created .env file with default values${NC}"
fi

# Function to check if a container is running and healthy
check_container() {
    local container_name=$1
    local max_retries=${2:-30}
    local retry_count=0

    while [ $retry_count -lt $max_retries ]; do
        local status=$(docker ps --filter "name=$container_name" --format "{{.Status}}")
        if [ -n "$status" ]; then
            if [[ $status == *"healthy"* ]] || [[ $status == *"starting"* ]]; then
                echo -e "${GREEN}✓ $container_name is running and healthy${NC}"
                return 0
            fi
        fi
        echo -n "."
        sleep 1
        ((retry_count++))
    done
    echo -e "\n${RED}✗ $container_name failed to start properly${NC}"
    return 1
}

# Stop existing containers if they exist
echo -e "${YELLOW}Stopping existing containers...${NC}"
docker-compose down 2>/dev/null || true

# Remove old containers and volumes
echo -e "${YELLOW}Cleaning up old containers and volumes...${NC}"
docker-compose rm -f 2>/dev/null || true
docker volume prune -f 2>/dev/null || true

# Start all services
echo -e "${YELLOW}Starting services...${NC}"
docker-compose up -d

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
for i in {1..30}; do
    if docker-compose exec postgres pg_isready -h localhost -U postgres &> /dev/null; then
        echo -e "${GREEN}✓ PostgreSQL is ready${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}Failed to connect to PostgreSQL after 30 seconds${NC}"
        cleanup
    fi
    echo -n "."
    sleep 1
done

# Check all services with health checks
echo -e "\n${YELLOW}Checking services status...${NC}"
check_container "userrestfulapi-postgres-1" 30 && \
check_container "userrestfulapi-app-1" 30 && \
check_container "userrestfulapi-nginx-1" 30

if [ $? -eq 0 ]; then
    echo -e "\n${GREEN}All services are running successfully!${NC}"
    echo -e "${YELLOW}API is available at: http://localhost:8080${NC}"
    echo -e "${YELLOW}Health check: http://localhost:8080/health${NC}"
    echo -e "${YELLOW}Metrics: http://localhost:8080/metrics${NC}"
    
    # Verify API health
    echo -e "\n${YELLOW}Verifying API health...${NC}"
    if curl -s http://localhost:8080/health | grep -q "ok"; then
        echo -e "${GREEN}✓ API is healthy${NC}"
    else
        echo -e "${RED}✗ API health check failed${NC}"
        cleanup
    fi
    
    # Show logs
    echo -e "\n${YELLOW}Showing application logs (Ctrl+C to exit):${NC}"
    docker-compose logs -f app
else
    echo -e "\n${RED}Some services failed to start. Showing logs:${NC}"
    docker-compose logs
    cleanup
fi
