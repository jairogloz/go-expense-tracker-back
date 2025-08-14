# Go Expense Tracker Backend

A backend service for parsing natural language expense descriptions into structured transactions using OpenAI and storing them in PostgreSQL.

## Architecture

This project follows hexagonal (clean) architecture principles with clear separation of concerns:

- **Presentation Layer** (`handlers/`): HTTP handlers using Gin framework
- **Application Layer** (`app/`): Use cases that orchestrate business logic
- **Domain Layer** (`domain/`): Core business models and interfaces
- **Service Layer** (`services/`): Business logic implementations
- **Infrastructure Layer** (`infra/`): External adapters (PostgreSQL, OpenAI)

## Features

- Parse natural language expense descriptions using OpenAI
- Store structured transactions in PostgreSQL
- RESTful API with JSON responses
- Graceful shutdown
- Environment-based configuration
- Clean architecture with dependency injection

## Project Structure

```
go-expense-tracker-back/
├── cmd/
│   └── server/                 # Application entry point
│       └── main.go
├── internal/
│   ├── app/                    # Use cases
│   │   └── parse_input_usecase.go
│   ├── domain/                 # Core models and interfaces
│   │   ├── transaction.go
│   │   └── ports.go
│   ├── handlers/               # HTTP handlers
│   │   └── transaction_handler.go
│   ├── services/               # Business logic
│   │   └── transaction_service.go
│   └── infra/                  # External adapters
│       ├── database.go
│       ├── openai_service.go
│       └── postgres_repository.go
├── config/
│   └── config.go               # Configuration management
├── .env                        # Environment variables
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- OpenAI API key

## Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd go-expense-tracker-back
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**

   ```sql
   CREATE DATABASE expense_tracker;
   CREATE USER expense_user WITH PASSWORD 'expense_password';
   GRANT ALL PRIVILEGES ON DATABASE expense_tracker TO expense_user;
   ```

4. **Configure environment variables**
   Copy `.env` file and update the values:

   ```bash
   # Database configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=expense_user
   DB_PASSWORD=expense_password
   DB_NAME=expense_tracker
   DB_SSLMODE=disable

   # OpenAI configuration
   OPENAI_API_KEY=your_openai_api_key_here

   # Server configuration
   PORT=8080
   ```

5. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

### Health Check

```
GET /health
```

Returns server health status.

### Parse Input

```
POST /parse
Content-Type: application/json

{
  "text": "I spent $25.50 at Starbucks for coffee this morning"
}
```

Response:

```json
{
  "transactions": [
    {
      "id": 1,
      "amount": 25.5,
      "currency": "USD",
      "category": "food",
      "type": "expense",
      "date": "2024-01-15T12:00:00Z",
      "vendor": "Starbucks",
      "description": "coffee this morning"
    }
  ],
  "message": "Successfully parsed and saved transactions"
}
```

### Get Transaction

```
GET /transactions/{id}
```

### Get Transactions

```
GET /transactions?limit=10&offset=0
```

## Supported Categories

### Expense Categories

- `food` - Restaurants, groceries, etc.
- `transport` - Gas, public transit, rideshare, etc.
- `utilities` - Electricity, water, internet, etc.
- `shopping` - Clothing, electronics, etc.
- `health` - Medical, pharmacy, fitness, etc.
- `education` - Books, courses, tuition, etc.
- `entertainment` - Movies, games, streaming, etc.
- `other` - Miscellaneous expenses

### Income Categories

- `salary` - Regular employment income
- `freelance` - Contract/freelance work
- `investments` - Dividends, capital gains, etc.
- `bonus` - Performance bonuses, gifts, etc.

## Example Usage

### Parsing Multiple Transactions

```bash
curl -X POST http://localhost:8080/parse \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Today I bought groceries for $85.30 at Whole Foods and paid $45 for gas at Shell. Also received my freelance payment of $500."
  }'
```

### Getting Transactions

```bash
# Get specific transaction
curl http://localhost:8080/transactions/1

# Get transactions with pagination
curl "http://localhost:8080/transactions?limit=20&offset=0"
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o expense-tracker cmd/server/main.go
```

## License

This project is licensed under the MIT License.
