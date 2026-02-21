# synapsis-order-service

## Overview
The `synapsis-order-service` is a Go-based microservice designed to manage orders. It provides RESTful APIs for creating, retrieving, and managing orders, and integrates with an inventory service via gRPC. The service is built with scalability and maintainability in mind, leveraging modern tools and frameworks.

## Features
- **Order Management**: Create, retrieve, and cancel orders.
- **Inventory Integration**: Validate stock availability via gRPC calls to the inventory service.
- **Database Persistence**: Store and manage order data using PostgreSQL.
- **Observability**: Elastic APM tracing for monitoring and debugging.
- **Validation**: Input validation for API requests.
- **Error Handling**: Consistent and structured error responses.

## Prerequisites
- **Go**: Version 1.18 or later.
- **PostgreSQL**: Version 13 or later.
- **Protocol Buffers Compiler**: `protoc` for generating gRPC code.
- **Make**: For running predefined tasks.

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/0x5w4/synapsis-order-service.git
cd synapsis-order-service
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Configure Environment Variables
Create a `.env` file in the root directory with the following variables:
```
DATABASE_URL=postgres://user:password@localhost:5432/synapsis_order_service?sslmode=disable
GRPC_INVENTORY_SERVICE=localhost:50051
APM_SERVER_URL=http://localhost:8200
```

### 4. Run Database Migrations
```bash
make migrate-up
```

### 5. Start the Service
```bash
make run
```

### 6. Generate gRPC Code (if needed)
If you modify the `.proto` files, regenerate the gRPC code:
```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

## API Endpoints

### 1. Create Order
**POST** `/api/v1/orders`
- **Description**: Create a new order.
- **Request Body**:
```json
{
  "user_id": "string",
  "items": [
    {
      "product_id": "string",
      "quantity": 1
    }
  ]
}
```
- **Response**:
```json
{
  "order_id": "string",
  "status": "created"
}
```

### 2. List Orders
**GET** `/api/v1/orders`
- **Description**: Retrieve a list of all orders.
- **Query Parameters**:
  - `page` (optional): Page number for pagination.
  - `perPage` (optional): Number of items per page.
- **Response**:
```json
[
  {
    "order_id": "string",
    "status": "string"
  }
]
```

### 3. Get Order Details
**GET** `/api/v1/orders/:id`
- **Description**: Retrieve details of a specific order.
- **Response**:
```json
{
  "order_id": "string",
  "status": "string",
  "items": [
    {
      "product_id": "string",
      "quantity": 1
    }
  ]
}
```

### 4. Cancel Order
**POST** `/api/v1/orders/:id/cancel`
- **Description**: Cancel an existing order.
- **Response**:
```json
{
  "order_id": "string",
  "status": "cancelled"
}
```

## Testing

### Run Unit Tests
```bash
make test
```

### Run Integration Tests
```bash
make test-integration
```

### Test Coverage
To check test coverage:
```bash
go test -cover ./...
```

## Development

### Code Structure
- **cmd/**: Entry point for the application.
- **config/**: Configuration files and utilities.
- **internal/**: Core application logic, including adapters and domain services.
- **pkg/**: Shared libraries and utilities.
- **proto/**: Protocol Buffers definitions and generated gRPC code.
- **migration/**: Database migration scripts.

### Common Commands
- **Run the application**:
  ```bash
  make run
  ```
- **Run database migrations**:
  ```bash
  make migrate-up
  ```
- **Rollback migrations**:
  ```bash
  make migrate-down
  ```
- **Generate mocks**:
  ```bash
  mockery --all --output=mocks
  ```


### Code Style
- Follow the Go coding standards.
- Use `golangci-lint` for linting:
  ```bash
  golangci-lint run
  ```

## Observability

### Elastic APM
The service integrates with Elastic APM for tracing. Ensure the `APM_SERVER_URL` environment variable is set correctly.

### Logs
Logs are written to the console in JSON format. Use a log aggregator for centralized logging.

## License
This project is licensed under the MIT License.