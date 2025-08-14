#!/bin/bash

# Example API usage script for Go Expense Tracker Backend
# Make sure the server is running on localhost:8080

API_BASE="http://localhost:8080"

echo "🚀 Testing Go Expense Tracker API"
echo "=================================="

# Health check
echo "1. Health Check:"
curl -s "$API_BASE/health" | jq . || echo "Response: $(curl -s $API_BASE/health)"
echo ""

# Parse a simple expense
echo "2. Parsing a simple expense:"
curl -s -X POST "$API_BASE/parse" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "I spent $25.50 at Starbucks for coffee this morning"
  }' | jq . || echo "Response: $(curl -s -X POST $API_BASE/parse -H 'Content-Type: application/json' -d '{"text": "I spent $25.50 at Starbucks for coffee this morning"}')"
echo ""

# Parse multiple transactions
echo "3. Parsing multiple transactions:"
curl -s -X POST "$API_BASE/parse" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Today I bought groceries for $85.30 at Whole Foods, paid $45 for gas at Shell, and received my freelance payment of $500 from ABC Company."
  }' | jq . || echo "Response: $(curl -s -X POST $API_BASE/parse -H 'Content-Type: application/json' -d '{"text": "Today I bought groceries for $85.30 at Whole Foods, paid $45 for gas at Shell, and received my freelance payment of $500 from ABC Company."}')"
echo ""

# Get transactions
echo "4. Getting transactions:"
curl -s "$API_BASE/transactions?limit=5" | jq . || echo "Response: $(curl -s $API_BASE/transactions?limit=5)"
echo ""

# Get specific transaction (ID 1)
echo "5. Getting specific transaction (ID 1):"
curl -s "$API_BASE/transactions/1" | jq . || echo "Response: $(curl -s $API_BASE/transactions/1)"
echo ""

echo "✅ API testing complete!"
