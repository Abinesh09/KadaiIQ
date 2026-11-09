# KadaiIQ

Multi-Tenant Retail Management and Sales Analytics SaaS platform.

Repository: [https://github.com/Abinesh09/KadaiIQ](https://github.com/Abinesh09/KadaiIQ)

---

## 1. Monorepo Architecture

KadaiIQ follows a clean monorepo architecture with strict separation between backend and frontend applications:

```text
KadaiIQ/
├── AGENTS.md                  # Permanent coding agent instructions
├── README.md                  # Project overview and run guide
├── docker-compose.yml         # Container orchestration (Web: 7070, API: 8080, DB: 5432)
├── .env.example               # Root environment variables template
├── .gitignore                 # Monorepo git ignores
│
├── backend/                   # Go + Fiber + GORM + PostgreSQL
│   ├── cmd/api/main.go        # Backend API entrypoint
│   ├── config/                # Environment configuration loader
│   ├── database/              # GORM database connection & pool
│   ├── migrations/            # Schema migration scripts
│   ├── internal/              # Core MVC internal packages
│   │   ├── controllers/       # HTTP handlers (no business logic)
│   │   ├── services/          # Business logic & tenant context
│   │   ├── repositories/      # Database queries & tenant filters
│   │   ├── models/            # GORM entities & relationships
│   │   ├── dto/               # API request & response structures
│   │   ├── routes/            # Route registrations & middlewares
│   │   ├── middleware/        # JWT auth, tenant isolation
│   │   ├── validators/        # Request payload validation
│   │   └── utils/             # Helper utilities
│   ├── tests/                 # Unit & integration tests
│   ├── go.mod                 # Go module definition
│   └── Dockerfile             # Multi-stage Go alpine build
│
└── frontend/                  # React + Vite + TypeScript + Tailwind + shadcn/ui
    ├── src/
    │   ├── components/
    │   │   ├── ui/            # Reusable shadcn/ui primitives
    │   │   ├── common/        # Shared presentation components
    │   │   └── layout/        # AppLayout, Header, Sidebar
    │   ├── features/          # Domain-driven feature modules
    │   ├── services/          # External API clients
    │   ├── hooks/             # Custom React hooks
    │   ├── routes/            # App routing with React Router
    │   ├── lib/               # Utility functions (cn, etc.)
    │   ├── App.tsx            # QueryClient & Router provider setup
    │   └── main.tsx           # React DOM root
    ├── package.json
    ├── vite.config.ts         # Vite configuration with port 7070 & proxy
    └── Dockerfile             # Multi-stage NGINX build on port 7070
```

---

## 2. Quick Start

### Running with Docker Compose (Recommended)

1. Ensure **Docker Desktop** is running.
2. Start all services:
   ```bash
   docker compose up --build
   ```
3. Access the services:
   - **Frontend Web UI**: [http://localhost:7070](http://localhost:7070)
   - **Backend API**: [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)
   - **PostgreSQL Database**: `localhost:5432` (`kadai_iq`)

To stop the containers:
```bash
docker compose down
```

---

### Running Locally (Development Mode)

#### 1. Backend (Go + Fiber)
```bash
cd backend
go run cmd/api/main.go
```
The API server starts on [http://localhost:8080](http://localhost:8080).

#### 2. Frontend (React + Vite)
```bash
cd frontend
npm install
npm run dev
```
The frontend dev server starts on [http://localhost:7070](http://localhost:7070) with auto-proxying for `/api/v1` to the Go backend.

---

## 3. Verification & Tests

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Build & Lint
```bash
cd frontend
npm run build
```
