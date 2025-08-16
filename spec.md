# Go Expense Tracker Backend API Specification

## Base URL

```
http://localhost:8080
```

## Endpoints

### 1. Health Check

**GET /health**

**Description:** Check if the API is running

**Request:** No body required

**Response:**

```json
{
  "status": "healthy",
  "time": "2024-08-14T15:30:00Z"
}
```

**Status Codes:** 200

---

### 2. Parse Natural Language Input

**POST /parse**

**Description:** Parse natural language text into structured transactions using AI

**Request Body:**

```json
{
  "text": "I spent 50 pesos at the grocery store for food today"
}
```

**Response:**

```json
{
  "transactions": [
    {
      "id": 0,
      "amount": 50.0,
      "currency": "MXN",
      "category": "food",
      "type": "expense",
      "date": "2024-08-14T15:30:00Z",
      "description": "Grocery store purchase"
    }
  ],
  "message": "Successfully parsed and saved transactions"
}
```

**Status Codes:**

- 200: Success
- 400: Invalid request body
- 500: Internal server error

---

### 3. Get All Transactions

**GET /transactions**

**Description:** Retrieve all transactions with pagination

**Query Parameters:**

- `limit` (optional): Number of transactions to return (default: 10)
- `offset` (optional): Number of transactions to skip (default: 0)

**Request:** No body required

**Response:**

```json
{
  "transactions": [
    {
      "id": 1,
      "amount": 50.0,
      "currency": "MXN",
      "category": "food",
      "type": "expense",
      "date": "2024-08-14T15:30:00Z",
      "description": "Grocery store purchase"
    },
    {
      "id": 2,
      "amount": 1500.0,
      "currency": "MXN",
      "category": "salary",
      "type": "income",
      "date": "2024-08-14T09:00:00Z",
      "description": "Monthly salary payment"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

**Status Codes:**

- 200: Success
- 500: Internal server error

---

### 4. Get Single Transaction

**GET /transactions/{id}**

**Description:** Retrieve a specific transaction by ID

**Path Parameters:**

- `id`: Transaction ID (integer)

**Request:** No body required

**Response:**

```json
{
  "id": 1,
  "amount": 50.0,
  "currency": "MXN",
  "category": "food",
  "type": "expense",
  "date": "2024-08-14T15:30:00Z",
  "description": "Grocery store purchase"
}
```

**Status Codes:**

- 200: Success
- 400: Invalid transaction ID
- 404: Transaction not found
- 500: Internal server error

---

### 5. Update Transaction

**PUT /transactions/{id}**

**Description:** Update an existing transaction

**Path Parameters:**

- `id`: Transaction ID (integer)

**Request Body:**

```json
{
  "amount": 75.5,
  "currency": "MXN",
  "category": "shopping",
  "type": "expense",
  "date": "2024-08-14T15:30:00Z",
  "description": "Updated: Shopping at the mall"
}
```

**Response:**

```json
{
  "id": 1,
  "amount": 75.5,
  "currency": "MXN",
  "category": "shopping",
  "type": "expense",
  "date": "2024-08-14T15:30:00Z",
  "description": "Updated: Shopping at the mall"
}
```

**Status Codes:**

- 200: Success
- 400: Invalid request body or transaction ID
- 404: Transaction not found
- 500: Internal server error

---

### 6. Delete Transaction

**DELETE /transactions/{id}**

**Description:** Delete a transaction

**Path Parameters:**

- `id`: Transaction ID (integer)

**Request:** No body required

**Response:**

```json
{
  "message": "Transaction deleted successfully"
}
```

**Status Codes:**

- 200: Success
- 400: Invalid transaction ID
- 404: Transaction not found
- 500: Internal server error

---

## Data Models

### Transaction

```json
{
  "id": 1,
  "amount": 50.0,
  "currency": "MXN",
  "category": "food",
  "type": "expense",
  "date": "2024-08-14T15:30:00Z",
  "description": "Transaction description"
}
```

### Available Categories

**Expense Categories:**

- food
- transport
- utilities
- shopping
- health
- education
- entertainment
- other

**Income Categories:**

- salary
- freelance
- investments
- bonus

### Transaction Types

- expense
- income

### Currency

- Default: MXN
- Format: 3-letter ISO code (USD, MXN, EUR, etc.)

### Date Format

- ISO 8601 format: "2024-08-14T15:30:00Z"

## Error Response Format

```json
{
  "error": "Error message",
  "details": "Detailed error information (optional)"
}
```

## CORS Support

The API includes CORS headers for cross-origin requests:

- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Accept, Authorization, Content-Type, X-CSRF-Token`

## Example Usage

### Creating transactions from natural language:

```bash
curl -X POST http://localhost:8080/parse \
  -H "Content-Type: application/json" \
  -d '{"text": "I spent 50 pesos on groceries and received 1500 pesos salary"}'
```

### Getting all transactions:

```bash
curl http://localhost:8080/transactions?limit=5&offset=0
```

### Updating a transaction:

```bash
curl -X PUT http://localhost:8080/transactions/1 \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 75.50,
    "currency": "MXN",
    "category": "food",
    "type": "expense",
    "date": "2024-08-14T15:30:00Z",
    "description": "Updated grocery purchase"
  }'
```

### Deleting a transaction:

```bash
curl -X DELETE http://localhost:8080/transactions/1
```
