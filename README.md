# synapsis-order-service

## Overview
The `synapsis-order-service` is a Go-based microservice designed to manage orders. It provides RESTful APIs for creating, retrieving, and managing orders, and integrates with an inventory service via gRPC.

## Features
- Create, retrieve, and cancel orders.
- Integration with an inventory service for stock validation.
- Database persistence using PostgreSQL.
- Elastic APM tracing for observability.

## Prerequisites
- Go 1.18 or later
- PostgreSQL 13 or later
- `protoc` (Protocol Buffers Compiler) for gRPC code generation

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

## API Endpoints

### 1. Create Order
**POST** `/api/v1/orders`
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

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for discussion.

## License
This project is licensed under the MIT License.