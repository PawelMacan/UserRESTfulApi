#!/bin/bash

# Base URL for the API
API_URL="http://localhost:8080"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Testing User Management API${NC}\n"

# 1. Create a new user
echo -e "${GREEN}1. Creating a new user...${NC}"
CREATE_RESPONSE=$(curl -s -X POST "${API_URL}/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "Test123!@#",
    "name": "John Doe"
  }')
echo "Response: $CREATE_RESPONSE"
USER_ID=$(echo $CREATE_RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
echo -e "Created user ID: $USER_ID\n"

# 2. Get the created user
echo -e "${GREEN}2. Getting user details...${NC}"
curl -s -X GET "${API_URL}/users/$USER_ID" \
  -H "Content-Type: application/json" | json_pp
echo -e "\n"

# 3. Update the user
echo -e "${GREEN}3. Updating user...${NC}"
curl -s -X PUT "${API_URL}/users/$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe.updated@example.com",
    "password": "NewTest123!@#",
    "name": "John Doe Updated"
  }' | json_pp
echo -e "\n"

# 4. List users with pagination
echo -e "${GREEN}4. Listing users (page 1, limit 10)...${NC}"
curl -s -X GET "${API_URL}/users?page=1&limit=10" \
  -H "Content-Type: application/json" | json_pp
echo -e "\n"

# 5. Delete the user
echo -e "${GREEN}5. Deleting user...${NC}"
curl -s -X DELETE "${API_URL}/users/$USER_ID" \
  -H "Content-Type: application/json"
echo -e "\n"

# 6. Verify deletion by trying to get the user
echo -e "${GREEN}6. Verifying deletion...${NC}"
curl -s -X GET "${API_URL}/users/$USER_ID" \
  -H "Content-Type: application/json" | json_pp
echo -e "\n"

# 7. Test invalid user creation (wrong email format)
echo -e "${GREEN}7. Testing invalid user creation (wrong email format)...${NC}"
curl -s -X POST "${API_URL}/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "Test123!@#",
    "name": "Invalid User"
  }' | json_pp
echo -e "\n"

# 8. Test invalid user creation (weak password)
echo -e "${GREEN}8. Testing invalid user creation (weak password)...${NC}"
curl -s -X POST "${API_URL}/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "weak",
    "name": "Weak Password User"
  }' | json_pp
echo -e "\n"

# Function to test password validation
test_password() {
    local password=$1
    local expected_status=$2
    local description=$3
    
    response=$(curl -s -w "%{http_code}" -X POST "$API_URL/users" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"test@example.com\",
            \"password\": \"$password\",
            \"name\": \"Test User\"
        }")
    
    status_code=${response: -3}
    response_body=${response:0:${#response}-3}
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓ Password Test ($description): PASSED${NC}"
    else
        echo -e "${RED}✗ Password Test ($description): FAILED${NC}"
        echo "Expected status: $expected_status, Got: $status_code"
        echo "Response: $response_body"
    fi
}

echo "Testing password validation..."

# Test valid password
test_password "ValidP@ssw0rd" "201" "Valid password"

# Test password too short
test_password "Short1!" "400" "Too short"

# Test password without uppercase
test_password "password123!" "400" "No uppercase"

# Test password without lowercase
test_password "PASSWORD123!" "400" "No lowercase"

# Test password without number
test_password "Password!@#" "400" "No number"

# Test password without special char
test_password "Password123" "400" "No special char"

# Test password with space
test_password "Valid Pass@123" "400" "Contains space"

# List users
echo -e "\nListing users..."
curl -s "$API_URL/users" | jq '.'
