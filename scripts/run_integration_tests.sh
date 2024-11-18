#!/bin/bash
set -e

# Clean up any existing test containers
echo "Cleaning up existing test containers..."
docker-compose -f docker-compose.test.yml down -v

# Start test containers and run tests
echo "Starting test containers and running tests..."
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Store the test exit code
TEST_EXIT_CODE=$?

# Cleanup
echo "Cleaning up test containers..."
docker-compose -f docker-compose.test.yml down -v

# Exit with the test exit code
exit $TEST_EXIT_CODE
