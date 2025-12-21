# Architecture Overview

## System Architecture

```
┌─────────────────┐
│  API Gateway    │  (Port: TBD)
│  (Partial)      │
└────────┬────────┘
         │
    ┌────┴────┬──────────┬──────────┬──────────┐
    │         │          │          │          │
┌───▼───┐ ┌──▼───┐  ┌───▼───┐  ┌───▼───┐  ┌───▼───┐
│ User  │ │Order │  │Inventory│ │Payment│ │Notify │
│Service│ │Service│  │ Service │ │Service│ │Service│
│ :8081 │ │ :8082│  │  :TBD  │  │ :TBD  │  │ :TBD  │
└───┬───┘ └──┬───┘  └───┬───┘  └───┬───┘  └───┬───┘
    │        │          │          │          │
    │        │          │          │          │
┌───▼───┐ ┌──▼───┐  ┌───▼───┐  ┌───▼───┐  ┌───▼───┐
│userdb │ │orderdb│  │inventory│ │payment│ │notify │
│ :5432 │ │ :5433│  │  :5434 │  │ :5435 │  │ :5436 │
└───────┘ └──────┘  └────────┘  └───────┘  └───────┘

         ┌──────────────┐
         │   RabbitMQ   │
         │  :5672/15672 │
         └──────────────┘
```

## Communication Patterns

### Synchronous (HTTP)
- Client → API Gateway → Services
- Service-to-service direct calls (when needed)

### Asynchronous (Events)
- Services publish events to RabbitMQ
- Other services consume events they're interested in
- Event-driven workflow for order processing

## Clean Architecture Layers

Each service follows this structure:

```
┌─────────────────────────────────────┐
│     Presentation Layer (HTTP)        │  ← Handlers, Routers
├─────────────────────────────────────┤
│     Application Layer (Use Cases)   │  ← Business Logic, DTOs
├─────────────────────────────────────┤
│     Domain Layer (Entities)          │  ← Core Business Rules
├─────────────────────────────────────┤
│  Infrastructure Layer (External)    │  ← DB, Messaging, Config
└─────────────────────────────────────┘
```

## Event Flow Architecture

```
Order Created
     │
     ▼
┌─────────────────┐
│  Order Service  │  Publishes: order.created
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ RabbitMQ Topic  │  Exchange: "ecommerce-events"
│    Exchange     │
└────────┬────────┘
         │
    ┌────┴────┬──────────┬──────────┐
    │         │          │          │
    ▼         ▼          ▼          ▼
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│Inventory│ │Payment │ │Order   │ │Notify  │
│Service  │ │Service │ │Service │ │Service │
│(Consumer)│ │(Consumer)│ │(Consumer)│ │(Consumer)│
└────────┘ └────────┘ └────────┘ └────────┘
```

## Database Architecture

- **Database per Service**: Each microservice has its own PostgreSQL database
- **No Shared Database**: Services cannot directly access each other's databases
- **Data Consistency**: Achieved through events and eventual consistency
- **Transaction Boundaries**: Limited to single service (distributed transactions via Saga)

## Deployment Architecture

- **Containerized**: Each service has a Dockerfile
- **Orchestrated**: Docker Compose for local development
- **Network**: All services on `microservices-network` bridge network
- **Port Mapping**: Each service exposed on different ports

## Security Architecture

- **Authentication**: Planned in API Gateway (JWT tokens)
- **Password Security**: bcrypt hashing in User Service
- **Network Isolation**: Services communicate via internal network
- **Data Encryption**: At rest (PostgreSQL) and in transit (HTTPS planned)

## Scalability Considerations

- **Horizontal Scaling**: Each service can scale independently
- **Stateless Services**: HTTP services are stateless (except database state)
- **Message Queue**: RabbitMQ handles load distribution
- **Database Scaling**: Each database can be scaled independently

## Resilience Patterns

- **Circuit Breaker**: (Planned for API Gateway)
- **Retry Logic**: (Planned for event consumers)
- **Dead Letter Queue**: (Planned for failed messages)
- **Health Checks**: Each service has `/health` endpoint
- **Graceful Degradation**: Services can operate partially if dependencies fail

