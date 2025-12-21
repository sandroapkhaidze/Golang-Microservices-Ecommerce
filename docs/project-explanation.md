# Golang Microservices E-Commerce Project - AI Prompt Documentation

## Project Overview

This is a **production-grade event-driven microservices e-commerce system** built in Go. The project demonstrates modern microservices architecture patterns including Clean Architecture, event-driven communication via RabbitMQ, and distributed transaction management using the Saga pattern.

## Architecture & Design Patterns

### Core Architecture Principles

1. **Clean Architecture (Hexagonal Architecture)**
   - Each service follows a layered architecture with clear separation of concerns:
     - **Domain Layer**: Core business entities and repository interfaces
     - **Application Layer**: Use cases (business logic) and DTOs
     - **Infrastructure Layer**: Database implementations, messaging, external services
     - **Presentation Layer**: HTTP handlers and routers

2. **Event-Driven Architecture**
   - Services communicate asynchronously via RabbitMQ message broker
   - Uses topic exchange pattern for flexible event routing
   - Events include correlation IDs for distributed tracing

3. **Database per Service Pattern**
   - Each microservice has its own PostgreSQL database
   - Complete data isolation between services
   - Services communicate only through events/APIs

4. **Saga Pattern (Planned)**
   - For managing distributed transactions across services
   - Saga directory exists but implementation is pending

## Technology Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 15 (separate database per service)
- **Message Broker**: RabbitMQ 3 (with management UI)
- **Code Generation**: sqlc (type-safe SQL queries)
- **Containerization**: Docker & Docker Compose
- **Password Hashing**: bcrypt (via golang.org/x/crypto)

## Infrastructure Setup

### Docker Compose Services

The `docker-compose.yml` defines:
- **6 PostgreSQL databases** (one per service):
  - `postgres-user` (port 5432)
  - `postgres-order` (port 5433)
  - `postgres-inventory` (port 5434)
  - `postgres-payment` (port 5435)
  - `postgres-notification` (port 5436)
  - `postgres-gateway` (port 5437)
- **RabbitMQ** (ports 5672 for AMQP, 15672 for management UI)

### Shared Components

Located in `/shared` directory:
- **Events Package**: Common event definitions for inter-service communication
  - `base_event.go`: Base event structure with correlation IDs
  - `order_events.go`: Order lifecycle events
  - `inventory_events.go`: Inventory reservation events
  - `payment_events.go`: Payment processing events
- **Messaging Package**: RabbitMQ connection and pub/sub abstractions
  - `connection.go`: RabbitMQ connection management
  - `publisher.go`: Event publishing functionality
  - `consumer.go`: Event consumption with handler pattern

## Service Structure & Implementation Status

### 1. User Service ✅ **FULLY IMPLEMENTED**

**Port**: 8081  
**Database**: userdb (port 5432)

#### Structure:
```
user-service/
├── cmd/main.go                    # Service entry point
├── internal/
│   ├── domain/
│   │   ├── entity/user.go         # User domain entity with password hashing
│   │   └── repository/            # Repository interface
│   ├── application/
│   │   ├── dto/user_dto.go        # Request/Response DTOs
│   │   └── usecase/
│   │       ├── register_user.go   # User registration logic
│   │       └── login_user.go      # User authentication logic
│   ├── infrastructure/
│   │   ├── config/database.go     # Database configuration
│   │   └── persistence/
│   │       ├── postgres_user_repository.go  # PostgreSQL implementation
│   │       └── sqlc/              # Generated SQL code
│   └── presentation/
│       └── http/
│           ├── user_handler.go    # HTTP handlers
│           └── router.go          # Route definitions
└── migrations/
    └── 001_create_users_table.sql # Database schema
```

#### Implemented Features:
- ✅ User registration with email validation
- ✅ User login with password authentication
- ✅ Password hashing using bcrypt
- ✅ User CRUD operations (Create, Read, Update, Delete)
- ✅ Email uniqueness validation
- ✅ Role-based access (customer/admin)
- ✅ Soft delete (IsActive flag)

#### API Endpoints:
- `POST /api/v1/users/register` - Register new user
- `POST /api/v1/users/login` - User authentication
- `GET /health` - Health check

#### Domain Model:
- **User Entity**: ID, Email, Password (hashed), FirstName, LastName, Role, IsActive, timestamps
- **Methods**: `HashPassword()`, `CheckPassword()`, `IsAdmin()`

---

### 2. Order Service ✅ **FULLY IMPLEMENTED**

**Port**: 8082  
**Database**: orderdb (port 5433)

#### Structure:
```
order-service/
├── cmd/main.go                    # Service entry point with RabbitMQ setup
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── order.go           # Order entity with status management
│   │   │   └── order_item.go      # Order item entity
│   │   └── repository/            # Repository interface
│   ├── application/
│   │   ├── dto/order_dto.go      # Request/Response DTOs
│   │   ├── saga/                 # Saga pattern (directory exists, empty)
│   │   └── usecase/
│   │       └── create_order.go   # Order creation with event publishing
│   ├── infrastructure/
│   │   ├── config/database.go    # Database configuration
│   │   └── persistence/
│   │       ├── postgres_order_repository.go  # PostgreSQL implementation
│   │       └── sqlc/             # Generated SQL code
│   └── presentation/
│       └── http/
│           ├── order_handler.go  # HTTP handlers
│           └── router.go        # Route definitions
└── migrations/
    └── 001_create_orders_tables.sql  # Database schema
```

#### Implemented Features:
- ✅ Order creation with items
- ✅ Order status management (pending, processing, completed, failed, cancelled)
- ✅ Total amount calculation
- ✅ Event publishing to RabbitMQ (`OrderCreatedEvent`)
- ✅ Correlation ID for distributed tracing
- ✅ Transaction support for order + items creation

#### API Endpoints:
- `POST /api/v1/orders` - Create new order
- `GET /health` - Health check

#### Domain Model:
- **Order Entity**: ID, UserID, Status, TotalAmount, Items[], CorrelationID, timestamps
- **Order Statuses**: pending, processing, completed, failed, cancelled
- **Methods**: `CalculateTotal()`, `MarkAsProcessing()`, `MarkAsCompleted()`, `MarkAsFailed()`, `MarkAsCancelled()`, `CanBeCancelled()`

#### Events Published:
- `order.created` - When a new order is created
- (Planned: `order.completed`, `order.failed`)

---

### 3. Inventory Service ⚠️ **NOT IMPLEMENTED**

**Database**: inventorydb (port 5434)

#### Status:
- Directory structure exists
- No Go files implemented
- Database configured in docker-compose
- Event definitions exist in shared/events

#### Planned Features (based on events):
- Inventory reservation for orders
- Stock management
- Product availability checking
- Event consumption: `order.created`
- Event publishing: `inventory.reserved`, `inventory.reservation_failed`

---

### 4. Payment Service ⚠️ **NOT IMPLEMENTED**

**Database**: paymentdb (port 5435)

#### Status:
- Directory structure exists
- Only empty `main.go` file
- Database configured in docker-compose
- Event definitions exist in shared/events

#### Planned Features (based on events):
- Payment processing
- Payment method management
- Transaction recording
- Event consumption: `inventory.reserved`
- Event publishing: `payment.processed`, `payment.failed`

---

### 5. Notification Service ⚠️ **NOT IMPLEMENTED**

**Database**: notificationdb (port 5436)

#### Status:
- Directory structure exists
- Only empty `main.go` file
- Database configured in docker-compose

#### Planned Features:
- Email notifications
- SMS notifications
- Push notifications
- Event consumption: Various order/payment events
- Notification history tracking

---

### 6. API Gateway ⚠️ **PARTIALLY IMPLEMENTED**

**Database**: gatewaydb (port 5437)

#### Status:
- Directory structure exists
- Internal structure: config, handler, middleware, proxy
- `main.go` exists but minimal implementation

#### Planned Features:
- Request routing to backend services
- Authentication/Authorization middleware
- Rate limiting
- Request aggregation
- Load balancing
- API versioning

## Event Flow (Planned Saga Pattern)

### Order Processing Flow:

1. **Order Service** receives order creation request
   - Creates order with status "pending"
   - Publishes `order.created` event

2. **Inventory Service** (consumer)
   - Consumes `order.created` event
   - Reserves inventory
   - Publishes `inventory.reserved` or `inventory.reservation_failed`

3. **Payment Service** (consumer)
   - Consumes `inventory.reserved` event
   - Processes payment
   - Publishes `payment.processed` or `payment.failed`

4. **Order Service** (consumer)
   - Consumes payment/inventory events
   - Updates order status accordingly
   - Publishes `order.completed` or `order.failed`

5. **Notification Service** (consumer)
   - Consumes order completion/failure events
   - Sends notifications to users

## Code Quality & Patterns

### Repository Pattern
- All services use repository interfaces in domain layer
- PostgreSQL implementations in infrastructure layer
- Type-safe queries generated by sqlc

### Use Case Pattern
- Business logic encapsulated in use case structs
- Each use case has a single `Execute()` method
- DTOs for input/output to maintain layer boundaries

### Dependency Injection
- Dependencies passed through constructors
- No global state
- Easy to test and mock

### Error Handling
- Domain errors returned from repositories
- Use cases handle business logic errors
- HTTP layer converts errors to appropriate status codes

## Database Schema

### Users Table
- UUID primary key
- Email (unique)
- Hashed password
- First/Last name
- Role (customer/admin)
- IsActive flag
- Timestamps with auto-update trigger

### Orders Table
- UUID primary key
- User ID (foreign key concept)
- Status (enum-like)
- Total amount (decimal)
- Correlation ID (unique, for tracing)
- Timestamps with auto-update trigger

### Order Items Table
- UUID primary key
- Order ID (foreign key with CASCADE)
- Product ID
- Quantity
- Price (decimal)

## Development Status Summary

| Service | Status | Features | Events |
|---------|--------|----------|--------|
| User Service | ✅ Complete | Registration, Login, CRUD | None (synchronous) |
| Order Service | ✅ Complete | Create Order, Status Management | Publishes `order.created` |
| Inventory Service | ❌ Not Started | - | Should consume/publish inventory events |
| Payment Service | ❌ Not Started | - | Should consume/publish payment events |
| Notification Service | ❌ Not Started | - | Should consume various events |
| API Gateway | ⚠️ Partial | Structure only | - |

## Key Implementation Details

### Event Publishing (Order Service)
- Uses RabbitMQ topic exchange named "ecommerce-events"
- Events include BaseEvent with correlation ID for tracing
- JSON serialization
- Persistent message delivery

### Database Access
- Uses sqlc for type-safe SQL queries
- SQL files in `queries/` directory
- Generated Go code in `sqlc/` directory
- Repository pattern abstracts database implementation

### Configuration
- Environment variable based configuration
- Default values for local development
- Database connection pooling
- RabbitMQ connection management

## Next Steps for Completion

1. **Implement Inventory Service**
   - Product entity and repository
   - Inventory reservation logic
   - Event consumers for order events
   - Event publishers for reservation results

2. **Implement Payment Service**
   - Payment entity and repository
   - Payment processing logic
   - Event consumers for inventory events
   - Event publishers for payment results

3. **Implement Notification Service**
   - Notification entity and repository
   - Email/SMS sending logic
   - Event consumers for various events

4. **Complete API Gateway**
   - Service discovery/routing
   - Authentication middleware
   - Request/response transformation

5. **Implement Saga Pattern**
   - Saga orchestrator or choreography
   - Compensation logic for rollbacks
   - State management

6. **Add Testing**
   - Unit tests for use cases
   - Integration tests for repositories
   - E2E tests for API endpoints

7. **Add Observability**
   - Logging (structured logging)
   - Metrics (Prometheus)
   - Distributed tracing (Jaeger/Zipkin)

## Usage Instructions

1. Start infrastructure: `docker-compose up -d`
2. Run migrations for each service
3. Start services individually (each has its own main.go)
4. Services communicate via RabbitMQ events
5. API Gateway routes requests to appropriate services

## Architecture Benefits

- **Scalability**: Each service can scale independently
- **Maintainability**: Clear separation of concerns
- **Testability**: Dependency injection enables easy mocking
- **Resilience**: Event-driven architecture handles failures gracefully
- **Technology Diversity**: Each service can use different tech if needed
- **Team Autonomy**: Teams can work on services independently

