#!/bin/bash

# CRUD Operations Test Script for Go Expense Tracker Backend
# Make sure the server is running on localhost:8080

API_BASE="http://localhost:8080"

echo "🧪 Testing CRUD Operations for Transactions"
echo "==========================================="

# First, let's create some transactions to work with
echo "1. Creating transactions via parse:"
PARSE_RESPONSE=$(curl -s -X POST "$API_BASE/parse" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "I spent 50 pesos at the grocery store for food today"
  }')
echo "$PARSE_RESPONSE" | jq . 2>/dev/null || echo "Response: $PARSE_RESPONSE"
echo ""

# Get all transactions to find IDs
echo "2. Getting all transactions:"
ALL_TRANSACTIONS=$(curl -s "$API_BASE/transactions?limit=10")
echo "$ALL_TRANSACTIONS" | jq . 2>/dev/null || echo "Response: $ALL_TRANSACTIONS"

# Extract first transaction ID for testing
TRANSACTION_ID=$(echo "$ALL_TRANSACTIONS" | jq -r '.transactions[0].id' 2>/dev/null)
echo "Using transaction ID: $TRANSACTION_ID for testing"
echo ""

# READ - Get specific transaction
echo "3. READ - Getting transaction $TRANSACTION_ID:"
GET_RESPONSE=$(curl -s "$API_BASE/transactions/$TRANSACTION_ID")
echo "$GET_RESPONSE" | jq . 2>/dev/null || echo "Response: $GET_RESPONSE"
echo ""

# UPDATE - Update the transaction
echo "4. UPDATE - Updating transaction $TRANSACTION_ID:"
UPDATE_RESPONSE=$(curl -s -X PUT "$API_BASE/transactions/$TRANSACTION_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 75.50,
    "currency": "MXN",
    "category": "shopping",
    "type": "expense",
    "date": "2024-08-14T15:30:00Z",
    "description": "Updated: Shopping at the mall"
  }')
echo "$UPDATE_RESPONSE" | jq . 2>/dev/null || echo "Response: $UPDATE_RESPONSE"
echo ""

# READ - Verify the update
echo "5. READ - Verifying update for transaction $TRANSACTION_ID:"
VERIFY_RESPONSE=$(curl -s "$API_BASE/transactions/$TRANSACTION_ID")
echo "$VERIFY_RESPONSE" | jq . 2>/dev/null || echo "Response: $VERIFY_RESPONSE"
echo ""

# Test with invalid data
echo "6. UPDATE - Testing with invalid data (should fail):"
INVALID_UPDATE=$(curl -s -X PUT "$API_BASE/transactions/$TRANSACTION_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": -50,
    "currency": "INVALID",
    "category": "invalid_category"
  }')
echo "$INVALID_UPDATE" | jq . 2>/dev/null || echo "Response: $INVALID_UPDATE"
echo ""

# Test updating non-existent transaction
echo "7. UPDATE - Testing update of non-existent transaction (should return 404):"
NONEXISTENT_UPDATE=$(curl -s -X PUT "$API_BASE/transactions/99999" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100,
    "currency": "MXN",
    "category": "food",
    "type": "expense",
    "date": "2024-08-14T12:00:00Z",
    "description": "This should not work"
  }')
echo "$NONEXISTENT_UPDATE" | jq . 2>/dev/null || echo "Response: $NONEXISTENT_UPDATE"
echo ""

# DELETE - Ask user before deleting
echo "8. DELETE - Do you want to test deletion of transaction $TRANSACTION_ID? (y/n)"
read -r CONFIRM

if [[ $CONFIRM =~ ^[Yy]$ ]]; then
    echo "Deleting transaction $TRANSACTION_ID:"
    DELETE_RESPONSE=$(curl -s -X DELETE "$API_BASE/transactions/$TRANSACTION_ID")
    echo "$DELETE_RESPONSE" | jq . 2>/dev/null || echo "Response: $DELETE_RESPONSE"
    echo ""
    
    # Verify deletion
    echo "Verifying deletion (should return 404):"
    VERIFY_DELETE=$(curl -s "$API_BASE/transactions/$TRANSACTION_ID")
    echo "$VERIFY_DELETE" | jq . 2>/dev/null || echo "Response: $VERIFY_DELETE"
    echo ""
else
    echo "Skipping deletion test."
    echo ""
fi

# Test deleting non-existent transaction
echo "9. DELETE - Testing deletion of non-existent transaction (should return 404):"
NONEXISTENT_DELETE=$(curl -s -X DELETE "$API_BASE/transactions/99999")
echo "$NONEXISTENT_DELETE" | jq . 2>/dev/null || echo "Response: $NONEXISTENT_DELETE"
echo ""

echo "✅ CRUD operations testing complete!"
echo ""
echo "📋 Summary of endpoints tested:"
echo "  POST   /parse                - Create transactions via AI parsing"
echo "  GET    /transactions         - List all transactions"
echo "  GET    /transactions/:id     - Get specific transaction"
echo "  PUT    /transactions/:id     - Update transaction"
echo "  DELETE /transactions/:id     - Delete transaction"
