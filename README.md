# Foodie Service

A Go-based microservice template with basic server functionality using the Fiber framework.

## Prerequisites

- Go 1.21 or higher
- Git

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd foodie-service
```

2. Install dependencies:
```bash
go mod download
```

## Running the Service

To start the service:

```bash
go run .
```

The server will start on port 3000 by default (or the port specified in your .env file).

## API Endpoints

### Authentication
- `POST /auth/signup` - User registration
  - Request Body: `{"email": "string", "password": "string", "name": "string"}`

- `POST /auth/login` - User login
  - Request Body: `{"email": "string", "password": "string"}`
  - Returns: `{"token": "string", "user": {...}}`
  - The returned token must be included in the Authorization header for protected routes
  

### Public Routes
- `GET /health` - Health check endpoint
  - Returns: `{"status": "healthy"}`
- `GET /products` - Get all products
- `GET /products/:id` - Get product by ID
- `POST /products` - Bulk insert products
- `GET /coupons` - Get available coupons

### Protected Routes
All protected routes require a valid JWT token in the header api-key:
```
api-key: <token>
```

Protected endpoints:
- `POST /orders` - Place a new order
- `GET /orders` - Get user's old orders

## Authentication
To access protected routes:
1. First, authenticate using the login endpoint
2. Copy the token from the login response
3. Include the token in the header for all subsequent requests to protected routes
4. Format: `api-key: <your-token>`

## Development

To add new endpoints, modify the `index.go` file in routes folder and add new routes to the Fiber app. Fiber provides a simple and intuitive API similar to Express.js.

## Project Structure

```
.
├── main.go          # Main server implementation using Fiber
├── go.mod           # Go module file
├── go.sum           # Go module checksum file
└── .env            # Environment variables (optional)
```

## Framework

This service uses [Fiber](https://github.com/gofiber/fiber), a fast, express-inspired web framework for Go. It's designed to be:
- Zero memory allocation
- Express-like API
- High performance
- Low memory footprint 

