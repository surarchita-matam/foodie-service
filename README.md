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

3. (Optional) Create a `.env` file:
```bash
PORT=8080
```

## Running the Service

To start the service:

```bash
go run main.go
```

The server will start on port 8080 by default (or the port specified in your .env file).

## API Endpoints

- `GET /health` - Health check endpoint
  - Returns: `{"status": "healthy"}`

## Development

To add new endpoints, modify the `main.go` file and add new routes to the Fiber app. Fiber provides a simple and intuitive API similar to Express.js.

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