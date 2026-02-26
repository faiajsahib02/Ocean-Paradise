# Oasis — Distributed Hotel Operations System

A full-stack hotel management system built with **Go** and **React + TypeScript**, designed around **Hexagonal Architecture**, **Domain-Driven Design**, and **SOLID principles**. The codebase is intentionally structured to demonstrate clean software engineering practices: loose coupling, dependency injection, and clear separation of concerns across every layer.

---

## Motivation

Building a hotel management system means juggling a lot: guests checking in and out, rooms changing status in real time, housekeeping tasks being assigned and tracked, laundry and restaurant orders overlapping, and all of it eventually rolling up into a single invoice. Most hobby projects shy away from that kind of complexity — I wanted to lean into it.

I built Oasis specifically to work with a domain that has **multiple interacting bounded contexts**, which made it the perfect playground for applying Hexagonal Architecture and Domain-Driven Design in a realistic setting. I also wanted to explore **RAG (Retrieval-Augmented Generation)** in a context where it actually makes sense: a hotel concierge that answers questions grounded in real hotel data, not generic LLM hallucinations.

The result is a system where every architectural decision — from the strict interface boundaries to the dependency injection wiring — has a purpose rooted in solving a real design problem.

---

## Architecture

The backend follows **Hexagonal Architecture (Ports & Adapters)**, keeping the domain core completely isolated from external concerns like HTTP, databases, or third-party APIs. Each layer depends only on abstractions — never on concrete implementations.

```
React + TypeScript (Frontend)
        |
        | HTTP / WebSocket
        v
REST Handlers  +  Middleware (JWT, CORS)
        |
        | calls via interface (Port)
        v
Service Layer — business logic, one bounded context per module
        |
        | depends on repository interface, not implementation
        v
Repository Layer — data access adapters
        |
        | sqlx
        v
PostgreSQL
```

---

## Design Principles

### Domain-Driven Design (DDD)

The backend is divided into **8 bounded contexts**, each owning its domain model, service, and repository. Domain entities live in `domain/` and are plain Go structs with no framework dependencies — they model the business, nothing else.

```
domain/
  guest.go          # Guest entity — check-in, check-out, profile
  room.go           # Room entity — status: vacant / occupied / cleaning
  invoice.go        # Invoice entity — aggregates charges across services
  housekeeping.go   # Housekeeping task entity
  laundry.go        # Laundry request + items
  restaurant.go     # Menu items and orders
  staff.go          # Staff profiles and roles
```

Naming reflects the business language directly — `CheckInDate`, `RoomStatus`, `InvoicePreview` — not database column names or HTTP fields.

### Hexagonal Architecture — Ports & Adapters

Each module exposes a **Port** (interface) and provides an **Adapter** (implementation). HTTP handlers and repositories both connect to the domain through these interfaces, making every component independently replaceable and testable.

```go
// port.go — defines what the module can do (the contract)
type Service interface {
    Create(room domain.Room) (*domain.Room, error)
    Find(roomNumber string) (*domain.Room, error)
    GetAll(status string) ([]domain.Room, error)
}

// service.go — implements the contract, depends on a repository interface
type service struct {
    rmRepo RoomRepo
}

func NewService(rmRepo RoomRepo) *service {
    return &service{rmRepo: rmRepo}
}
```

The HTTP handler depends on `room.Service` (the interface), not the concrete `service` struct. The repository is also behind an interface. Neither the handler nor the service ever imports a concrete type from the other layer.

### Dependency Injection

All dependencies are wired once at startup in `cmd/serve.go` using **constructor-based injection**. No service creates its own dependencies. No global state. The dependency graph is explicit and fully controlled.

```go
// cmd/serve.go — full DI composition root
db         := infra.NewConnection(cfg)
roomRepo   := repository.NewRoomRepo(db)
roomSvc    := room.NewService(roomRepo)
roomHandler := roomhandler.NewHandler(roomSvc)
```

This means every component can be tested in isolation by injecting a mock — the business logic is never coupled to a real database or HTTP layer.

### Loose Coupling

- Handlers depend on service **interfaces**, not implementations
- Services depend on repository **interfaces**, not concrete types
- Domain entities have zero dependencies on any external package
- No circular imports — the dependency flow is strictly one-directional

Swapping PostgreSQL for another database, or replacing the REST layer with gRPC, would require changes only at the adapter level — the domain and service logic remain untouched.

### SOLID Principles

| Principle                 | Application                                                                                                                                            |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Single Responsibility** | Each service owns exactly one bounded context. Handlers translate HTTP, services apply business rules, repositories access data — each does one thing. |
| **Open / Closed**         | New behaviour is added by implementing existing interfaces, not by modifying existing services.                                                        |
| **Liskov Substitution**   | Any repository implementation can replace another without breaking the service that depends on it.                                                     |
| **Interface Segregation** | Each module defines a small, focused interface exposing only what its consumers actually need.                                                         |
| **Dependency Inversion**  | High-level modules (services) depend on abstractions (interfaces). Low-level modules (repositories) implement those abstractions.                      |

---

## AI Concierge — Retrieval-Augmented Generation (RAG)

One of the more technically interesting features of this system is the AI-powered guest concierge, built using a **RAG (Retrieval-Augmented Generation)** pipeline.

Rather than relying on a general-purpose LLM that hallucinates hotel-specific information, the system embeds hotel documents (policies, menus, room details) into a **pgvector** store in PostgreSQL. When a guest asks a question, the relevant context is retrieved via **vector similarity search** and injected into the prompt before being sent to **OpenAI GPT** — so responses are grounded in actual hotel data.

```
Guest query
    |
    v
Embed query → vector search against pgvector (PostgreSQL)
    |
    | top-k matching document chunks
    v
Construct prompt: [retrieved context] + [user question]
    |
    v
OpenAI GPT → contextually accurate, hotel-specific response
```

This is implemented as its own bounded context (`rag/`) following the same Hexagonal Architecture pattern as every other module — the OpenAI and pgvector dependencies are adapters behind an interface, keeping the core retrieval logic testable and decoupled.

---

## Tech Stack

| Layer     | Technology                                                                   |
| --------- | ---------------------------------------------------------------------------- |
| Backend   | Go 1.22, Gorilla Mux, Gorilla WebSocket, sqlx, sql-migrate, godotenv         |
| Auth      | Custom JWT — HMAC-SHA256 signing, role claims, middleware enforcement        |
| AI        | OpenAI API — RAG-based concierge with vector search (pgvector)               |
| Frontend  | React 18, TypeScript, Vite, Tailwind CSS, React Router, Axios, Framer Motion |
| Database  | PostgreSQL 14+, 11 versioned migrations                                      |
| Real-time | WebSocket hub with goroutine-per-client and broadcast channel                |

---

## Modules

| Module         | Responsibility                                                                                                   |
| -------------- | ---------------------------------------------------------------------------------------------------------------- |
| `guest`        | Registration, authentication, check-in/out, profiles                                                             |
| `room`         | Inventory, availability, status lifecycle                                                                        |
| `staff`        | Staff accounts, roles, RBAC                                                                                      |
| `housekeeping` | Task assignment, real-time status via WebSocket                                                                  |
| `laundry`      | Requests, item-level tracking, status management                                                                 |
| `restaurant`   | Menu management (with soft delete), orders                                                                       |
| `invoice`      | Aggregates charges from multiple services into a single bill                                                     |
| `rag`          | AI concierge — RAG pipeline: pgvector similarity search + OpenAI GPT for context-aware, hotel-grounded responses |

---

## Project Structure

```
Backend/
├── cmd/serve.go          # Composition root — wires all dependencies
├── config/               # Environment-based configuration
├── domain/               # Business entities — no framework imports
├── {module}/
│   ├── port.go           # Service interface (the contract)
│   └── service.go        # Business logic implementation
├── repository/           # Database adapters (implement service interfaces)
├── rest/
│   ├── handlers/         # HTTP adapters — translate requests to service calls
│   └── middlewares/      # JWT auth, CORS, role enforcement
├── infra/db/             # Database connection and migration runner
├── migrations/           # SQL migration files (versioned)
├── ws/                   # WebSocket hub and client management
└── util/                 # JWT generation/parsing, response helpers

Frontend/src/
├── context/AuthContext   # Global auth state via Context API
├── pages/                # Role-scoped pages: guest, staff, admin
├── components/           # Reusable UI components
├── layouts/              # Guest, Staff, Public layout wrappers
├── services/api.ts       # Axios instance and API layer
└── types/index.ts        # Shared TypeScript interfaces
```

---

## Quick Start

**Prerequisites:** Go 1.22+, Node.js 18+, PostgreSQL 14+

### 1. Clone the repository

```bash
git clone https://github.com/your-username/oasis.git
cd oasis
```

### 2. Set up and run the backend

```bash
cd Backend
go mod download

# Backend/.env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=oasis
JWT_SECRET=your_secret
PORT=8080
OPENAI_API_KEY=your_key

go run main.go migrate
go run main.go
# API available at http://localhost:8080
```

### 3. Set up and run the frontend

```bash
# In a new terminal
cd Frontend
npm install
```

Create `Frontend/.env`:

```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
```

```bash
npm run dev
# App available at http://localhost:5173
```

---

## Usage

Once the app is running, navigate to `http://localhost:5173`.

### Guest Portal

- **Register / log in** as a guest to access the guest dashboard
- **View available rooms** and browse the hotel's room inventory
- **Submit service requests** — laundry pick-up, restaurant orders, housekeeping calls
- **Chat with the AI Concierge** — ask questions about the hotel and get answers grounded in real hotel data via the RAG pipeline
- **View your invoice** — all charges from all services are aggregated into a single bill

### Staff Portal

- **Log in** with a staff account to access the staff dashboard
- **Manage room statuses** — mark rooms as vacant, occupied, or in-cleaning
- **Handle housekeeping tasks** in real time — task status updates are pushed via WebSocket
- **Process laundry requests** — view items, update status per request
- **Manage restaurant orders** and the menu (add, update, soft-delete items)
- **Generate guest invoices** aggregating all service charges

### Admin

- Full access to all staff functionality
- **Manage staff accounts** — create and assign roles (staff / admin)
- Configure hotel documents used by the AI Concierge RAG pipeline (menus, policies, room details)

### API Endpoints

The REST API is available at `http://localhost:8080`. Key endpoint groups:

| Prefix             | Description                          |
| ------------------ | ------------------------------------ |
| `POST /guest/register` | Guest registration                 |
| `POST /guest/login`    | Guest & staff authentication       |
| `GET  /rooms`          | Room availability                  |
| `POST /laundry`        | Submit a laundry request           |
| `POST /restaurant/order` | Place a restaurant order         |
| `GET  /invoice/:id`    | Fetch aggregated invoice for guest |
| `POST /rag/ask`        | Query the AI concierge             |
| `WS   /ws`             | WebSocket connection for real-time updates |

All protected routes require a **JWT Bearer token** in the `Authorization` header.

---

## Contributing

### Clone the repo

```bash
git clone https://github.com/your-username/oasis.git
cd oasis
```

### Run the backend

```bash
cd Backend
go mod download
go run main.go migrate   # run database migrations
go run main.go           # start the API server
```

### Run the frontend

```bash
cd Frontend
npm install
npm run dev
```

### Run the backend tests

```bash
cd Backend
go test ./...
```

### Submit a pull request

If you'd like to contribute, please **fork the repository** and open a pull request to the `main` branch. Make sure your code follows the existing hexagonal structure — new features should be added as their own bounded context with a `port.go` and `service.go`, wired up in `cmd/serve.go`.

---

## License

MIT
