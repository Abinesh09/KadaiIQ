# KadaiIQ — Multi-Tenant Retail Management & Advanced Sales Analytics SaaS

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![React Version](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react)](https://react.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5%2B-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://www.docker.com)

**GitHub Repository:** [https://github.com/Abinesh09/KadaiIQ](https://github.com/Abinesh09/KadaiIQ)

---

## 1. Why KadaiIQ is Useful (The Business Value)

Retail businesses (footwear stores, apparel shops, fashion boutiques, and multi-branch retail chains) face major operational challenges:

1. **Fragmented POS Systems:** Sales files are scattered across legacy offline billing counters (e.g. ZenBiller, CSVs, Excel sheets). Business owners cannot see daily, monthly, or yearly progress in one central place.
2. **Stock Chaos by Size & Colour:** Selling footwear is complex. A shop may have 50 pairs of shoes, but if sizes 7 and 8 are sold out, customers walk away. Most generic ERPs do not track variant-level size curves or colour distributions.
3. **Mystery Margins & Profit Blindspots:** High sales volume doesn't guarantee high profit. When discounts, GST, and customer returns are mixed together, owners don't know their true Cost of Goods Sold (COGS) or net take-home margins.
4. **Dead Capital in Zero-Sales Inventory:** Money gets locked in slow-moving or zero-sales stock without early warning signals.
5. **Tax Filing Headaches:** Reconciling CGST, SGST, IGST, and HSN-wise sales for Indian GSTR-1 and GSTR-3B monthly filings takes hours of manual accountant work.

### How KadaiIQ Solves This

- **Zero-Manual Data Entry:** Drag-and-drop POS sales files directly into KadaiIQ. The system automatically reconciles transactions, checks deduplication hashes, and adjusts stock quantities in seconds.
- **Variant-Level Matrix:** Tracks footwear by size (e.g. 6, 7, 8, 9, 10, 11) and colour curves so you always stock the sizes that actually sell.
- **Accounting-Grade Financial Visibility:** Distinct tracking of Gross Sales, Discounts, Customer Returns, Taxes, Net Sales, COGS, and Gross Profit Margins with exact decimal arithmetic (`shopspring/decimal`).
- **Dead-Stock Watchdog:** Proactively highlights zero-sales and slow-moving inventory so store managers can discount or liquidate stale stock before it depreciates.
- **One-Click Indian GST Summaries:** Generates ready-to-file GSTR-1 (Table 7 & Table 12 HSN) and GSTR-3B Table 3.1 summaries.
- **Enterprise Multi-Tenancy:** A retail owner or franchise group can run multiple shops or brands under one account with strict tenant data isolation.

---

## 2. What You Have (Complete Inventory of Implemented Features)

KadaiIQ is a fully functional, enterprise-grade full-stack platform. Here is everything built into the codebase:

### 1. Executive Business Overview Dashboard
- **Hero Net Sales Card:** High-contrast revenue tracking with period-over-period percentage growth.
- **Financial Summary Ribbon:** 6-column accounting breakdown (Gross Sales, Discounts Applied, Returns Subtracted, Net Realized Sales, Cash/UPI/Card Collections, and True Gross Profit).
- **Interactive Revenue Trajectory:** Recharts dual-curve area chart comparing Gross Sales vs Net Sales over time.
- **Quick Horizons:** One-click toggles for `Today`, `Last 7 days`, `Last 30 days`, `This month`, and `This year`.
- **Secondary Intelligence Grid:** Real-time widgets for Top 5 Best Sellers, Zero-Sales Watchdog, Category Share, Tender Breakdown, and 5 Latest POS Bills.

### 2. Product Catalogue & Footwear Variant Matrix
- Multi-variant product catalog: Product Name, Code, Brand, Category, HSN Code, Base Price, and Cost Price.
- Granular variant matrix tracking: Shoe Sizes (e.g., 6, 7, 8, 9, 10, 11), Colours (Black, Tan, Brown, White), SKUs, and Barcodes.
- Barcode scanning & quick lookup for fast retail POS identification.

### 3. Inventory & Immutable Stock Movement Ledger
- Stock-on-hand tracking for every product variant.
- **Immutable Movement Ledger:** Every inventory change creates an auditable `StockMovement` entry:
  - `OPENING_STOCK` — Initial balance setup
  - `SALE` — Automated POS sales deduction
  - `SALES_RETURN` — Stock replenishment from customer returns
  - `PURCHASE_RESTOCK` — Inward supplier shipments
  - `MANUAL_ADJUSTMENT` — Physical audit corrections
  - `TRANSFER` — Inter-branch store transfers
- Zero inventory changes without an audit trail.

### 4. POS Sales & Import Engine
- **ZenBiller POS Support:** Drag-and-drop batch upload for ZenBiller CSV and Excel (.xlsx) export files.
- **Instant File Preview:** Validates row count, gross revenue, items, tax calculations, and error warnings before writing to the database.
- **Deduplication Shield:** Composite unique constraint (`tenant_id`, `source_system`, `source_bill_id`) prevents duplicate sales imports and stops double inventory deductions.
- **Import Audit History:** Complete logs showing file names, timestamps, total rows processed, duplicate bills blocked, and execution status.

### 5. Sales Performance & Advanced Revenue Analytics
- **Sales Trajectory View:** Visual area charts tracking daily and weekly bill volumes and net sales.
- **Daily Breakdown Ledger:** Day-by-day table showing invoice counts, units moved, gross sales, discounts, returns, taxes, and net revenue.
- **Monthly Breakdown Table:** Month-over-month performance progression.
- **Yearly Comparison Table:** Multi-year growth metrics.
- **Product Performance Intelligence:** Filter products into *Top-Selling*, *Slow-Moving*, and *Zero-Sales Alerts* with margins and remaining stock on hand.
- **Footwear Size & Colour Curves:** Bar charts showing which shoe sizes generate the highest revenue and donut charts of colour preferences.
- **Payment Tender Breakdown:** Distribution across Cash, UPI (GPay/PhonePe/Paytm), Credit/Debit Cards, and Net Banking.
- **Multi-Sheet Excel (.xlsx) Export:** One-click download containing all financial metrics, trends, and product data in clean spreadsheet sheets.

### 6. Indian GST Compliance & Reports
- Tenant-specific GSTIN configuration (with automated state code resolution).
- Intra-state (CGST + SGST) vs Inter-state (IGST) automated tax split.
- Configurable GST slabs per footwear category (e.g., 5% vs 12% vs 18%).
- **GSTR-1 Table 7:** B2C taxable turnover and tax breakup.
- **GSTR-1 Table 12:** HSN-wise summary with quantity, taxable value, and tax rates.
- **GSTR-3B Table 3.1:** Outward taxable supplies and net tax liability calculation.

### 7. Financial Reports & RFC-4180 CSV Exports
- **Sales Ledger Report:** Invoices, bill numbers, customer details, payment methods, and tax values.
- **Footwear Valuation Report:** Current stock quantities, purchase cost, selling price, and total inventory asset value.
- **Stock Movement Audit Trail:** Complete tamper-evident log of all adjustments and restocks.
- All reports downloadable as clean RFC-4180 CSV files.

### 8. Multi-Tenant Architecture & RBAC Security
- Strict multi-tenant isolation at the database repository layer.
- Tenant context resolved strictly from authenticated user JWTs (`X-Tenant-ID` header).
- Role-Based Access Control (RBAC):
  - **Super Admin:** Platform-wide oversight and tenant creation.
  - **Store Admin:** Full store control, pricing, imports, and user management.
  - **Store Manager:** Day-to-day operations, inventory, and sales analytics.
  - **Cashier:** POS sales creation and bill viewing.
- OWASP security headers (`nosniff`, `DENY`, `mode=block`) and non-root Docker runtime.

---

## 3. Technology Stack

| Layer | Technologies Used |
|:---|:---|
| **Backend API** | Go 1.24, Fiber v2, GORM, `pgx`, `shopspring/decimal`, `excelize/v2`, JWT |
| **Frontend Web** | React 19, TypeScript (Strict Mode), Vite 8, Tailwind CSS v4, shadcn/ui, TanStack Query v5, Recharts, Lucide React |
| **Database** | PostgreSQL 16 with versioned idempotent migrations (`000001` - `000005`) |
| **Infrastructure** | Docker, Docker Compose, NGINX reverse proxy with zero-cache headers |

---

## 4. How to Run KadaiIQ

### Method A: Running with Docker Compose (Recommended)

Running via Docker Compose brings up the frontend, backend API, and PostgreSQL database with a single command.

#### Step 1: Clone the repository
```bash
git clone https://github.com/Abinesh09/KadaiIQ.git
cd KadaiIQ
```

#### Step 2: Start all containers
```bash
docker compose up --build -d
```

#### Step 3: Access the platform in your browser
- **Frontend Dashboard:** [http://localhost:7070](http://localhost:7070)
- **Backend API Health:** [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)
- **Database Connection Health:** [http://localhost:8080/api/v1/health/db](http://localhost:8080/api/v1/health/db)
- **PostgreSQL Database:** `localhost:5432` (database: `kadai_iq`, user: `kadai_user`)

#### Step 4: Login Credentials
When the database initializes for the first time, a default Super Administrator account is seeded:
- **Email:** `superadmin@kadaiiq.com`
- **Password:** `Admin@KadaiIQ2026`

*(You can also register a brand new retail store account directly on the login screen at [http://localhost:7070/register](http://localhost:7070/register).)*

#### To Stop the Application:
```bash
docker compose down
```
*(Do not add the `-v` flag unless you want to wipe your database volume).*

---

### Method B: Running in Local Development Mode (Without Docker)

#### Prerequisites
- **Go:** Version 1.24 or higher ([Download](https://golang.org/dl/))
- **Node.js:** Version 20 or higher & npm ([Download](https://nodejs.org/))
- **PostgreSQL:** Version 16 running locally on port `5432`

#### Step 1: Setup the Database
Create the database in PostgreSQL:
```sql
CREATE DATABASE kadai_iq;
CREATE USER kadai_user WITH PASSWORD 'kadai_password';
GRANT ALL PRIVILEGES ON DATABASE kadai_iq TO kadai_user;
```

#### Step 2: Run the Go Backend
```bash
cd backend
cp .env.example .env
# Verify your DATABASE_URL in .env
go run cmd/api/main.go
```
The Go Fiber server starts on **`http://localhost:8080`**. Migrations run automatically on startup.

#### Step 3: Run the React Frontend
In a new terminal:
```bash
cd frontend
npm install
npm run dev
```
The Vite development server will start on **`http://localhost:7070`**.

---

## 5. Automated Tests & Verification

### Run Backend Tests
KadaiIQ includes 12 automated test suites covering authentication, migrations, inventory deduction idempotency, sales import deduplication, analytics calculations, Indian GST logic, and multi-tenant isolation:
```bash
cd backend
go test -v ./...
```

### Run Frontend Production Build & Type-Check
Verify strict TypeScript compilation with zero errors:
```bash
cd frontend
npm run build
```

---

## 6. Project Directory Structure

```text
KadaiIQ/
├── AGENTS.md                  # Development rules & architectural constraints
├── README.md                  # Project documentation & runbook
├── docker-compose.yml         # Container orchestration (Web: 7070, API: 8080, DB: 5432)
├── .env.example               # Environment variable templates
│
├── backend/                   # Go backend (MVC Architecture)
│   ├── cmd/api/main.go        # Server entrypoint with graceful shutdown
│   ├── config/                # Environment variable configuration
│   ├── database/              # PostgreSQL connection pool & health checks
│   ├── migrations/            # Versioned migrations (000001 to 000005)
│   ├── internal/
│   │   ├── controllers/       # HTTP request handlers
│   │   ├── services/          # Pure business logic & financial math
│   │   ├── repositories/      # Database queries with strict tenant isolation
│   │   ├── models/            # GORM database entities
│   │   ├── dto/               # Strongly-typed request/response models
│   │   ├── routes/            # Route registration & CORS setup
│   │   └── middleware/        # JWT auth, tenant context, RBAC
│   ├── tests/                 # Unit & integration test suites
│   └── Dockerfile             # Multi-stage Alpine container (non-root)
│
└── frontend/                  # React client (Vite + TypeScript)
    ├── src/
    │   ├── components/ui/     # Reusable UI primitives
    │   ├── features/
    │   │   ├── dashboard/     # Executive real-time KPI overview
    │   │   ├── products/      # Footwear product catalog & barcode lookup
    │   │   ├── inventory/     # Stock ledger & valuation
    │   │   ├── sales/         # POS data import & ledger
    │   │   ├── analytics/     # Revenue trajectory curves & size matrix
    │   │   ├── gst/           # Indian GST settings & GSTR reports
    │   │   └── reports/       # CSV/Excel export generators
    │   ├── services/api.ts    # Centralized API service client
    │   └── App.tsx            # Routes & query provider
    ├── nginx.conf             # Production NGINX server with reverse proxy
    └── Dockerfile             # Multi-stage builder + NGINX runtime
```

---

## 7. License

Proprietary SaaS — Copyright © 2026 KadaiIQ. All rights reserved.
