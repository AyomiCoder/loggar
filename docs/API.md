# API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication

All API endpoints except `/auth/login` require JWT authentication via the `Authorization` header.

### Headers
```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

---

## Endpoints

### 1. Health Check

**GET** `/health`

Check if the API server is running.

**Response:**
```json
{
  "status": "ok"
}
```

---

### 2. User Login

**POST** `/auth/login`

Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "email": "test@loggar.dev",
  "password": "test123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "email": "test@loggar.dev"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Invalid credentials
- `500 Internal Server Error` - Database error

**Example with curl:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@loggar.dev",
    "password": "test123"
  }'
```

---

### 3. Analyze Logs

**POST** `/api/analyze`

Analyze log content using AI to identify issues and get recommendations.

**Headers:**
```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "logs": "2026-01-15 10:23:45 ERROR [database] Connection pool exhausted\n2026-01-15 10:23:46 ERROR [auth] Timeout waiting for database connection\n2026-01-15 10:23:47 ERROR [payment] Failed to process payment - database unavailable"
}
```

**Response:**
```json
{
  "primary_issue": "Database connection pool exhaustion",
  "secondary_effects": [
    "Auth service timeouts",
    "Payment retries failing"
  ],
  "first_seen": "2026-01-15T10:23:45Z",
  "likely_causes": [
    {
      "cause": "Unreleased DB connections",
      "confidence": 0.63
    },
    {
      "cause": "Traffic spike exceeded pool size",
      "confidence": 0.27
    }
  ],
  "recommended_actions": [
    "Check connection release in auth middleware",
    "Inspect pool max size vs current RPS"
  ],
  "similar_past_incidents": [
    {
      "date": "2025-11-12",
      "resolution": "Fixed middleware leak"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request` - Missing or invalid logs field
- `401 Unauthorized` - Missing or invalid JWT token
- `500 Internal Server Error` - AI analysis failed

**Example with curl:**
```bash
# First, get your token
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@loggar.dev","password":"test123"}' | jq -r '.token')

# Then analyze logs
curl -X POST http://localhost:8080/api/analyze \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "logs": "ERROR: Connection timeout\nERROR: Database unavailable"
  }'
```

---

## Database Setup

### 1. Create Database
```bash
createdb loggar
```

### 2. Run Schema
```bash
psql loggar < scripts/setup-db.sql
```

### 3. Test User
The setup script creates a test user:
- **Email:** `test@loggar.dev`
- **Password:** `test123`

---

## Running the Server

### 1. Set up environment
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 2. Start the server
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

---

## Testing the API

### Using curl

**1. Login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@loggar.dev","password":"test123"}'
```

**2. Save the token:**
```bash
export TOKEN="<paste_token_here>"
```

**3. Analyze logs:**
```bash
curl -X POST http://localhost:8080/api/analyze \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"logs":"ERROR: Database connection failed"}'
```

### Using Postman

1. **Login:**
   - Method: POST
   - URL: `http://localhost:8080/auth/login`
   - Body (JSON):
     ```json
     {
       "email": "test@loggar.dev",
       "password": "test123"
     }
     ```

2. **Copy the token from response**

3. **Analyze:**
   - Method: POST
   - URL: `http://localhost:8080/api/analyze`
   - Headers:
     - `Authorization: Bearer <your_token>`
     - `Content-Type: application/json`
   - Body (JSON):
     ```json
     {
       "logs": "Your log content here"
     }
     ```
