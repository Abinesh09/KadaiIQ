# KadaiIQ — Multi-Tenant Retail Management & Advanced Sales Analytics SaaS

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react)](https://react.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5%2B-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://www.docker.com)

**GitHub Repository:** [https://github.com/Abinesh09/KadaiIQ](https://github.com/Abinesh09/KadaiIQ)

KadaiIQ is an enterprise-grade, multi-tenant Retail Management and Advanced Sales Analytics SaaS platform designed specifically for retail organizations (with specialized intelligence for footwear, apparel, and consumer goods).

---

## 1. System Architecture

KadaiIQ adheres to clean architecture principles with complete separation between the Go backend service and React frontend client:

```text
KadaiIQ/
├── AGENTS.md                  # Development guidelines & architectural constraints
├── README.md                  # System documentation & runbook
├── docker-compose.yml         # Container orchestration (Web: 7070, API: 8080, DB: 5432)
├── .env.example               # Environment variables specification
├── .dockerignore              # Container build exclusions
│
├── backend/                   # Go + Fiber + GORM + PostgreSQL (MVC Architecture)
│   ├── cmd/api/main.go        # Server entrypoint with graceful shutdown
│   ├── config/                # Environment configuration loader
│   ├── database/              # Connection pooling & health probes
│   ├── migrations/            # Versioned migrations runner (000001 - 000005)
│   ├── internal/
│   │   ├── controllers/       # HTTP request/response handlers
│   │   ├── services/          # Pure business logic & financial calculations
│   │   ├── repositories/      # Database queries with strict tenant isolation
│   │   ├── models/            # GORM database entities & relationships
│   │   ├── dto/               # Strongly-typed request/response payloads
│   │   ├── routes/            # Route registration, CORS, security headers
│   │   ├── middleware/        # JWT auth, tenant context, RBAC authorization
│   │   └── utils/             # JWT, cryptographic password hashing
│   ├── tests/                 # Comprehensive unit & integration test suites
│   ├── go.mod                 # Go dependencies
│   └── Dockerfile             # Multi-stage Alpine container (non-root user)
│
└── frontend/                  # React + Vite + TypeScript + Tailwind CSS + shadcn/ui
    ├── src/
    │   ├── components/ui/     # Reusable shadcn/ui primitives (Button, Card, Badge)
    │   ├── features/
    │   │   ├── auth/          # Login, Register, Forgot Password
    │   │   ├── dashboard/     # Executive real-time KPI overview
    │   │   ├── tenants/       # Tenant onboarding & settings
    │   │   ├── users/         # Team member invitations & RBAC roles
    │   │   ├── products/      # Footwear product catalog & barcode lookup
    │   │   ├── inventory/     # Stock ledger, adjustments & valuation
    │   │   ├── sales/         # POS data import (ZenBiller), preview & invoice viewer
    │   │   ├── analytics/     # Recharts trend curves, margins & footwear curves
    │   │   ├── gst/           # Indian GST settings, GSTR-1 & GSTR-3B summaries
    │   │   └── reports/       # Sales, valuation, audit reports & CSV downloads
    │   ├── services/api.ts    # Centralized API service client
    │   ├── context/           # Authentication state context
    │   ├── App.tsx            # QueryClient & React Router setup
    │   └── main.tsx           # Application root
    ├── nginx.conf             # Production NGINX server on port 7070 with API reverse-proxy
    ├── package.json           # Dependencies (TanStack Query, Recharts, Lucide)
    └── Dockerfile             # Multi-stage Node builder + NGINX runtime
```

---

## 2. Technology Stack

### Backend
- **Language:** Go 1.24+
- **HTTP Framework:** Fiber v2
- **ORM:** GORM with PostgreSQL driver (`pgx`)
- **Database:** PostgreSQL 16
- **Arithmetic:** `shopspring/decimal` for exact financial calculations (zero floating-point inaccuracies)
- **Data Parser:** `xuri/excelize/v2` for Excel (.xlsx) and Go standard `encoding/csv`
- **Security:** OWASP secure headers, Bcrypt password hashing, JWT access & refresh tokens

### Frontend
- **Framework:** React 19 + TypeScript (Strict Mode)
- **Build Tool:** Vite 8
- **Styling:** Tailwind CSS v4 + shadcn/ui design language
- **State Management:** TanStack React Query v5
- **Charts & Visualizations:** Recharts (Area, Bar, Donut curves)
- **Icons:** Lucide React

---

## 3. Implemented Roadmap Modules

| Phase | Module | Status | Highlights |
|:---|:---|:---:|:---|
| **Phase 1** | **Project Foundation** | ✅ Complete | Go Fiber setup, React Vite setup, Docker compose, health check endpoint. |
| **Phase 2** | **Database & Multi-Tenant** | ✅ Complete | GORM PostgreSQL pooling, versioned migration engine, tenant isolation boundaries. |
| **Phase 3** | **Authentication & RBAC** | ✅ Complete | JWT tokens, password reset flow, Super Admin vs Tenant Admin scoping. |
| **Phase 4** | **Products & Inventory** | ✅ Complete | Footwear variants (sizes, colours, barcodes), immutable stock ledger (`StockMovement`). |
| **Phase 5** | **POS Sales Data Import** | ✅ Complete | ZenBiller/CSV/Excel parser, composite deduplication index, idempotent stock deduction. |
| **Phase 6** | **Advanced Analytics** | ✅ Complete | Recharts revenue curves, COGS & margin tracking, footwear size & colour intelligence. |
| **Phase 7** | **Reports & Indian GST** | ✅ Complete | GSTIN validation, GSTR-1 (Table 7 & 12 HSN), GSTR-3B Table 3.1, downloadable CSV reports. |
| **Phase 8** | **Security & Deployment** | ✅ Complete | OWASP headers, non-root Docker builds, port 7070 web container, end-to-end tests. |

---

## 4. Quick Start & Deployment

### Running with Docker Compose (Recommended)

1. Ensure **Docker Desktop** is running.
2. Clone and start all services:
   ```bash
   git clone https://github.com/Abinesh09/KadaiIQ.git
   cd KadaiIQ
   docker compose up --build -d
   ```
3. Access KadaiIQ:
   - **Frontend Web UI:** [http://localhost:7070](http://localhost:7070)
   - **Backend API:** [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)
   - **Database Health Check:** [http://localhost:8080/api/v1/health/db](http://localhost:8080/api/v1/health/db)
   - **PostgreSQL Database:** `localhost:5432` (database: `kadai_iq`, user: `kadai_user`)

To stop all services:
```bash
docker compose down
```

---

### Running in Local Development Mode

#### Prerequisites
- Go 1.24+
- Node.js 20+ and npm
- PostgreSQL 16 running on `localhost:5432`

#### 1. Setup Backend
```bash
cd backend
cp .env.example .env   # Configure DB credentials
go run cmd/api/main.go
```
The Go API starts on `http://localhost:8080`.

#### 2. Setup Frontend
```bash
cd frontend
npm install
npm run dev
```
The Vite development server starts on `http://localhost:7070`.

---

## 5. Security & Multi-Tenant Model

1. **Mandatory Tenant Scoping:** Every business record (`Product`, `ProductVariant`, `StockMovement`, `Sale`, `SaleItem`, `ImportJob`, `TenantGSTSetting`) has a foreign key to `tenant_id`.
2. **Context Resolution:** The backend extracts the authenticated tenant from the verified JWT claims and header `X-Tenant-ID`. It is impossible for Tenant A to query or mutate Tenant B's data.
3. **OWASP HTTP Security Headers:** Enforces `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `X-XSS-Protection: 1; mode=block`, and `Referrer-Policy: strict-origin-when-cross-origin`.
4. **Non-Root Docker Containers:** Production backend container runs under unprivileged user `kadaiuser` (UID 10001).
5. **Deduplication Shield:** Composite unique index `idx_tenant_source_bill (tenant_id, source_system, source_bill_id)` prevents POS data re-imports from double counting sales or double deducting stock.

---

## 6. Testing & Quality Assurance

### Run Backend Test Suites
```bash
cd backend
go test -v ./...
```
Executes 12 test suites verifying:
- Authentication & tenant onboarding
- JWT token expiration
- Migration idempotency
- Product catalogue & barcode lookup
- Stock movements & inventory deduction idempotency
- Sales import parser & deduplication
- Advanced analytics & margin intelligence
- Indian GST validation, GSTR-1 & GSTR-3B calculations
- Report generation & CSV streaming exports
- OWASP security headers & role-based access control
- Strict multi-tenant isolation across all modules

### Run Frontend Production Build & Typecheck
```bash
cd frontend
npm run build
```
Ensures 100% strict TypeScript compliance with zero warnings or errors.

---

## 7. License

Proprietary SaaS — Copyright © 2026 KadaiIQ. All rights reserved.
