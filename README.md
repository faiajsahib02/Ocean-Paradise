# Oasis — Hotel Management System

Full-stack hotel management system built with **Go** and **React + TypeScript**.
Manages guests, rooms, staff, housekeeping, laundry, restaurant, invoicing, and an AI concierge — with real-time WebSocket updates and role-based access control.

## Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go, Gorilla Mux/WebSocket, sqlx, sql-migrate, custom JWT (HMAC-SHA256), OpenAI API |
| Frontend | React 18, TypeScript, Vite, Tailwind CSS, React Router, Axios, Framer Motion |
| Database | PostgreSQL (11 versioned migrations) |

## Architecture

Hexagonal / Clean Architecture with Domain-Driven Design.

```
Frontend (React + TS)
  ↕ HTTP / WebSocket
REST Handlers → Middleware (JWT auth, CORS)
  ↓
Service Layer (8 bounded contexts)
  ↓ interfaces (ports)
Repository Layer (adapters)
  ↓ sqlx
PostgreSQL
```

Each domain module (`guest/`, `room/`, `staff/`, `laundry/`, `restaurant/`, `housekeeping/`, `invoice/`, `rag/`) contains:
- **port.go** — service interface
- **service.go** — implementation, receives repository via constructor injection

## Features

| Module | Highlights |
|--------|------------|
| Guest | Registration, auth, check-in/out, profile, booking history |
| Room | Inventory, status tracking (vacant/occupied/cleaning), availability |
| Staff | Auth, roles, task assignment, RBAC |
| Housekeeping | Task assignment & tracking, real-time WebSocket updates |
| Laundry | Requests, item tracking, status management |
| Restaurant | Menu CRUD (soft-delete), orders |
| Invoice | Multi-service charge aggregation, billing |
| RAG | AI concierge powered by OpenAI GPT |

Cross-cutting: JWT authentication, role-based access control (guest / staff / admin), CORS, WebSocket real-time hub.

## Project Structure

```
Backend/
├── cmd/serve.go            # App bootstrap & dependency injection
├── config/                 # Env-based configuration
├── domain/                 # Pure business entities (DDD)
├── guest/ room/ staff/     # Service modules (port.go + service.go)
│   laundry/ restaurant/
│   housekeeping/ invoice/ rag/
├── repository/             # Database implementations
├── rest/
│   ├── server.go           # Router setup
│   ├── handlers/           # HTTP handlers per module
│   └── middlewares/        # Auth, CORS
├── infra/db/               # Connection & migration runner
├── migrations/             # SQL migration files
├── ws/                     # WebSocket hub & client
└── util/                   # JWT helpers, response formatting

Frontend/src/
├── App.tsx                 # Root component & routing
├── context/AuthContext.tsx  # Global auth state
├── pages/                  # Guest, staff, admin pages
├── components/             # Shared UI (Button, Card, Navbar, ProtectedRoute)
├── layouts/                # Guest / Staff / Public layouts
├── services/api.ts         # Axios HTTP client
└── types/index.ts          # TypeScript interfaces
```

## Getting Started

### Prerequisites

Go 1.22+, Node.js 18+, PostgreSQL 14+

### Backend

```bash
cd Backend
go mod download

# create Backend/.env
# DB_HOST=localhost DB_PORT=5432 DB_USER=... DB_PASSWORD=... DB_NAME=oasis
# JWT_SECRET=... PORT=8080 OPENAI_API_KEY=...

go run main.go migrate
go run main.go          # → http://localhost:8080
```

### Frontend

```bash
cd Frontend
npm install

# create Frontend/.env
# VITE_API_URL=http://localhost:8080
# VITE_WS_URL=ws://localhost:8080/ws

npm run dev             # → http://localhost:5173
```

## License

MIT

