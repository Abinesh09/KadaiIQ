# AGENTS.md

# KadaiIQ – AI Coding Agent Instructions

## 1. Project Identity

Project: KadaiIQ

Description: Multi-Tenant Retail Management and Sales Analytics SaaS.

Repository: [https://github.com/Abinesh09/KadaiIQ](https://github.com/Abinesh09/KadaiIQ)

KadaiIQ is a multi-tenant platform for retail organisations to manage products, inventory, sales and business analytics.

The main product priority is advanced sales analytics.

## 2. Mandatory Architecture

The repository must maintain two independent application directories:

- `backend/` – All Go backend code.
- `frontend/` – All React frontend code.

Never mix frontend application code with backend application code.

### Backend stack

- Go
- Fiber
- GORM
- PostgreSQL
- REST APIs
- JWT authentication

### Frontend stack

- React
- TypeScript
- Vite
- Tailwind CSS
- shadcn/ui
- React Router
- TanStack Query
- Recharts

## 3. Backend Coding Rules

Follow MVC architecture with service and repository layers.

Required separation:

- Models: Database entities and relationships.
- Controllers: HTTP request and response handling.
- Services: Business logic.
- Repositories: Database operations.
- Routes: Endpoint registration.
- Middleware: Authentication, tenant context and authorisation.
- DTO: Request and response structures.
- Validators: Input validation.

Rules:

1. Controllers must not contain business logic.
2. Services must not depend directly on HTTP request objects.
3. Database queries must be handled through repositories.
4. Use GORM consistently.
5. Use database transactions for operations that modify multiple related records.
6. Avoid duplicate business logic.
7. Return consistent API responses.
8. Handle errors explicitly.
9. Use context.Context where appropriate.
10. Do not expose database models directly in API responses.

## 4. Frontend Coding Rules

Organise the frontend by features.

Rules:

1. Use TypeScript strict mode.
2. Prefer functional React components.
3. Use reusable components.
4. Keep API calls outside UI components.
5. Use TanStack Query for server state.
6. Use React Hook Form and Zod for forms and validation.
7. Use shadcn/ui as the base UI component library.
8. Use Tailwind CSS for styling.
9. Use Lucide React icons.
10. Avoid large components with unrelated responsibilities.
11. Keep types and interfaces organised.
12. Do not use `any` unless absolutely necessary and documented.
13. Provide loading, empty, success and error states.
14. Ensure responsive layouts.
15. Avoid unnecessary inline styles and duplicated CSS.

## 5. Multi-Tenant Security

Tenant isolation is mandatory.

1. Every tenant-owned record must be associated with a tenant.
2. Resolve tenant access from authenticated user membership.
3. Never trust frontend-supplied tenant IDs.
4. Every tenant-scoped repository query must enforce tenant filtering.
5. Prevent cross-tenant access in all API endpoints.
6. Apply tenant isolation to exports and analytics.
7. Super Admin access must be explicitly authorised.
8. Test tenant isolation for every important feature.
9. Never bypass tenant middleware for convenience.
10. Keep tenant-specific settings separate.

Security must be enforced in the backend, not only in the frontend.

## 6. Database Rules

Use PostgreSQL and GORM.

- Use migrations for schema changes.
- Add appropriate indexes.
- Define foreign keys and constraints.
- Use numeric/decimal types for money.
- Store timestamps consistently in UTC.
- Use unique constraints where required.
- Avoid unnecessary database duplication.
- Use transactions for sales and stock updates.
- Preserve historical transaction values.

Do not automatically delete historical sales when a product is deactivated.

## 7. Inventory Rules

Inventory must be tracked at product-variant level.

For example, footwear variants may have different:

- Sizes
- Colours
- Barcodes
- SKUs
- Stock quantities

All inventory changes must create stock movement records.

Supported movement types should include:

- Opening stock
- Purchase/restock
- Sale
- Sales return
- Manual adjustment
- Transfer

Never change stock quantities without maintaining movement history.

## 8. Sales and Analytics Rules

- Keep original source transaction IDs when importing POS data.
- Prevent duplicate sales imports.
- Handle cancelled bills and returns correctly.
- Keep gross sales, discounts, tax and net sales distinct.
- Use reliable purchase cost data before calculating profit.
- Apply tenant and date filters to all analytics.
- Ensure report totals match the underlying sales records.
- Use appropriate date/time boundaries and tenant time zones.

Barcode scans alone must never be considered completed sales.

## 9. GST Rules

- GST settings must be tenant-specific.
- Applicable GST rates must be configurable per product.
- Preserve tax details at transaction time.
- Support applicable CGST, SGST and IGST calculations.
- Validate rates against current Indian GST rules.
- Do not assume all footwear products have identical GST rates.
- Do not silently recalculate historical tax values.

## 10. API Rules

Use REST APIs under:

`/api/v1`

Keep endpoints organised by module.

- Authentication
- Super Admin
- Tenants
- Users
- Products
- Inventory
- Sales
- Analytics
- Reports
- GST

Use consistent:

- HTTP status codes.
- Response formats.
- Error handling.
- Pagination.
- Filtering.
- Sorting.
- Validation.

Document important APIs.

## 11. Environment and Configuration

- Use environment variables.
- Maintain `.env.example`.
- Never commit `.env` or credentials.
- Never hardcode database passwords, JWT secrets or API keys.
- Keep development and production configurations separate.
- Use Docker Compose for local services.

## 12. Code Quality

Before submitting changes:

Backend:

- Run gofmt.
- Run go vet where applicable.
- Run `go test ./...`.
- Check error handling.

Frontend:

- Run lint.
- Run TypeScript checks.
- Run production build.
- Fix all introduced errors.

Do not leave broken builds or failing tests.

## 13. Git Rules

- Use clear commit messages.
- Keep changes focused.
- Avoid committing generated files unnecessarily.
- Do not overwrite existing working code without inspection.
- Never commit secrets.
- Do not force-push or delete branches without explicit permission.
- Do not merge pull requests without authorisation.

## 14. Agent Workflow

For every task:

1. Read this AGENTS.md.
2. Inspect the existing implementation.
3. Identify the relevant backend and frontend modules.
4. Explain the implementation plan.
5. Make the smallest maintainable changes.
6. Run relevant tests and checks.
7. Report files changed.
8. Report commands executed and their results.
9. Mention unresolved issues honestly.

Do not modify unrelated modules.

Do not implement multiple future phases without approval.

Ask before making destructive database or repository changes.

## 15. Development Priority

Follow this order unless explicitly instructed otherwise:

1. Repository foundation.
2. Backend and frontend setup.
3. Database configuration.
4. Authentication.
5. Super Admin and tenant management.
6. Tenant isolation and permissions.
7. Product management.
8. Inventory management.
9. Sales import.
10. Analytics.
11. Reports.
12. GST.
13. Testing and deployment.

## 16. Core Principle

Build KadaiIQ as a real, maintainable, secure multi-tenant SaaS product.

Prefer clean architecture, reusable components, clear responsibilities, strong tenant isolation and accurate business calculations over shortcuts.
