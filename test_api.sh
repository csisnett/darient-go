#!/bin/bash

# API Testing Script
# Make sure your server is running on localhost:8080 before running this script

BASE_URL="http://localhost:8080"

echo "=== API Testing Script ==="
echo "Testing endpoints for the backend API"

# Test Health Check
echo "1. Testing Health Check..."
curl -X GET "$BASE_URL/health" -H "Content-Type: application/json"
echo -e "\n"

# Test Items
echo "2. Testing Items API..."
echo "Creating an item..."
ITEM_RESPONSE=$(curl -s -X POST "$BASE_URL/api/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Item",
    "description": "This is a test item"
  }')
echo $ITEM_RESPONSE
ITEM_ID=$(echo $ITEM_RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
echo -e "\n"

echo "Getting all items..."
curl -X GET "$BASE_URL/api/items" -H "Content-Type: application/json"
echo -e "\n"

if [ ! -z "$ITEM_ID" ]; then
  echo "Getting item by ID ($ITEM_ID)..."
  curl -X GET "$BASE_URL/api/items/$ITEM_ID" -H "Content-Type: application/json"
  echo -e "\n"
fi

# Test Clients
echo "3. Testing Clients API..."
echo "Creating a client..."
CLIENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/clients" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "birth_date": "1990-01-01T00:00:00Z",
    "country": "USA"
  }')
echo $CLIENT_RESPONSE
CLIENT_ID=$(echo $CLIENT_RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
echo -e "\n"

if [ ! -z "$CLIENT_ID" ]; then
  echo "Getting client by ID ($CLIENT_ID)..."
  curl -X GET "$BASE_URL/api/clients/$CLIENT_ID" -H "Content-Type: application/json"
  echo -e "\n"

  echo "Updating client..."
  curl -X PUT "$BASE_URL/api/clients/$CLIENT_ID" \
    -H "Content-Type: application/json" \
    -d '{
      "full_name": "Jane Doe",
      "email": "jane.doe@example.com",
      "birth_date": "1990-01-01T00:00:00Z",
      "country": "Canada"
    }'
  echo -e "\n"
fi

# Test Banks
echo "4. Testing Banks API..."
echo "Creating a bank..."
BANK_RESPONSE=$(curl -s -X POST "$BASE_URL/api/banks" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Bank",
    "type": "PRIVATE"
  }')
echo $BANK_RESPONSE
BANK_ID=$(echo $BANK_RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
echo -e "\n"

echo "Getting all banks..."
curl -X GET "$BASE_URL/api/banks" -H "Content-Type: application/json"
echo -e "\n"

if [ ! -z "$BANK_ID" ]; then
  echo "Getting bank by ID ($BANK_ID)..."
  curl -X GET "$BASE_URL/api/banks/$BANK_ID" -H "Content-Type: application/json"
  echo -e "\n"

  echo "Updating bank..."
  curl -X PUT "$BASE_URL/api/banks/$BANK_ID" \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Updated Test Bank",
      "type": "GOVERNMENT"
    }'
  echo -e "\n"
fi

# Test Credits
echo "5. Testing Credits API..."
if [ ! -z "$CLIENT_ID" ] && [ ! -z "$BANK_ID" ]; then
  echo "Creating a credit..."
  CREDIT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/credits" \
    -H "Content-Type: application/json" \
    -d "{
      \"client_id\": $CLIENT_ID,
      \"bank_id\": $BANK_ID,
      \"min_payment\": 100.0,
      \"max_payment\": 1000.0,
      \"term_months\": 12,
      \"credit_type\": \"AUTO\",
      \"status\": \"PENDING\"
    }")
  echo $CREDIT_RESPONSE
  CREDIT_ID=$(echo $CREDIT_RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
  echo -e "\n"

  echo "Getting all credits..."
  curl -X GET "$BASE_URL/api/credits" -H "Content-Type: application/json"
  echo -e "\n"

  if [ ! -z "$CREDIT_ID" ]; then
    echo "Getting credit by ID ($CREDIT_ID)..."
    curl -X GET "$BASE_URL/api/credits/$CREDIT_ID" -H "Content-Type: application/json"
    echo -e "\n"

    echo "Getting credits by client ID ($CLIENT_ID)..."
    curl -X GET "$BASE_URL/api/clients/$CLIENT_ID/credits" -H "Content-Type: application/json"
    echo -e "\n"

    echo "Getting credits by bank ID ($BANK_ID)..."
    curl -X GET "$BASE_URL/api/banks/$BANK_ID/credits" -H "Content-Type: application/json"
    echo -e "\n"

    echo "Updating credit status..."
    curl -X PUT "$BASE_URL/api/credits/$CREDIT_ID" \
      -H "Content-Type: application/json" \
      -d "{
        \"client_id\": $CLIENT_ID,
        \"bank_id\": $BANK_ID,
        \"min_payment\": 150.0,
        \"max_payment\": 1500.0,
        \"term_months\": 24,
        \"credit_type\": \"MORTGAGE\",
        \"status\": \"APPROVED\"
      }"
    echo -e "\n"
  fi
fi

# Test Error Cases
echo "6. Testing Error Cases..."
echo "Testing invalid item ID..."
curl -X GET "$BASE_URL/api/items/999999" -H "Content-Type: application/json"
echo -e "\n"

echo "Testing invalid client ID..."
curl -X GET "$BASE_URL/api/clients/999999" -H "Content-Type: application/json"
echo -e "\n"

echo "Testing invalid bank ID..."
curl -X GET "$BASE_URL/api/banks/999999" -H "Content-Type: application/json"
echo -e "\n"

echo "Testing invalid credit ID..."
curl -X GET "$BASE_URL/api/credits/999999" -H "Content-Type: application/json"
echo -e "\n"

echo "Testing invalid JSON..."
curl -X POST "$BASE_URL/api/items" \
  -H "Content-Type: application/json" \
  -d '{"invalid": json}'
echo -e "\n"

echo "Testing invalid bank type..."
curl -X POST "$BASE_URL/api/banks" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid Bank",
    "type": "INVALID_TYPE"
  }'
echo -e "\n"

echo "Testing invalid credit type..."
curl -X POST "$BASE_URL/api/credits" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "bank_id": 1,
    "min_payment": 100.0,
    "max_payment": 1000.0,
    "term_months": 12,
    "credit_type": "INVALID_TYPE",
    "status": "PENDING"
  }'
echo -e "\n"

# Cleanup (optional - uncomment if you want to delete test data)
echo "7. Cleanup (optional)..."
echo "Uncomment the lines below if you want to delete the test data"

# if [ ! -z "$CREDIT_ID" ]; then
#   echo "Deleting credit..."
#   curl -X DELETE "$BASE_URL/api/credits/$CREDIT_ID"
#   echo -e "\n"
# fi

# if [ ! -z "$CLIENT_ID" ]; then
#   echo "Deleting client..."
#   curl -X DELETE "$BASE_URL/api/clients/$CLIENT_ID"
#   echo -e "\n"
# fi

# if [ ! -z "$BANK_ID" ]; then
#   echo "Deleting bank..."
#   curl -X DELETE "$BASE_URL/api/banks/$BANK_ID"
#   echo -e "\n"
# fi

echo "=== API Testing Complete ==="