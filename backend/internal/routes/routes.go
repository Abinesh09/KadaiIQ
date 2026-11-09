package routes

import (
	"github.com/Abinesh09/KadaiIQ/backend/config"
	"github.com/Abinesh09/KadaiIQ/backend/database"
	"github.com/Abinesh09/KadaiIQ/backend/internal/controllers"
	"github.com/Abinesh09/KadaiIQ/backend/internal/middleware"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/Abinesh09/KadaiIQ/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, dbManager *database.DBManager) {
	// Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:7070,http://127.0.0.1:7070,http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Tenant-ID",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// API v1 group
	api := app.Group("/api/v1")

	// Health routes
	healthController := controllers.NewHealthController(cfg.AppVersion, dbManager)
	api.Get("/health", healthController.Check)
	api.Get("/health/db", healthController.CheckDB)

	// If DB is initialized, mount Auth, Tenant, Product, Inventory, and Sales routes
	if dbManager != nil && dbManager.DB != nil {
		authRepo := repositories.NewAuthRepository(dbManager.DB)
		tenantRepo := repositories.NewTenantRepository(dbManager.DB)
		productRepo := repositories.NewProductRepository(dbManager.DB)
		saleRepo := repositories.NewSaleRepository(dbManager.DB)
		analyticsRepo := repositories.NewAnalyticsRepository(dbManager.DB)

		authService := services.NewAuthService(authRepo, tenantRepo, cfg)
		tenantService := services.NewTenantService(tenantRepo, authRepo)
		productService := services.NewProductService(productRepo)
		saleService := services.NewSaleService(saleRepo, productRepo)
		analyticsService := services.NewAnalyticsService(analyticsRepo)

		authCtrl := controllers.NewAuthController(authService)
		tenantCtrl := controllers.NewTenantController(tenantService)
		productCtrl := controllers.NewProductController(productService)
		saleCtrl := controllers.NewSaleController(saleService)
		analyticsCtrl := controllers.NewAnalyticsController(analyticsService)

		jwtAuth := middleware.JWTMiddleware(cfg.DatabaseURL)
		tenantScope := middleware.RequireTenantContext()

		// Public Auth Routes
		auth := api.Group("/auth")
		auth.Post("/login", authCtrl.Login)
		auth.Post("/register", authCtrl.Register)
		auth.Post("/refresh", authCtrl.RefreshToken)
		auth.Post("/logout", authCtrl.Logout)
		auth.Post("/forgot-password", authCtrl.ForgotPassword)
		auth.Post("/reset-password", authCtrl.ResetPassword)

		// Protected Auth Routes
		auth.Get("/me", jwtAuth, authCtrl.Me)
		auth.Post("/switch-tenant", jwtAuth, authCtrl.SwitchTenant)

		// Super Admin Management Routes
		superAdmin := api.Group("/super-admin", jwtAuth, middleware.RequireSuperAdmin())
		superAdmin.Get("/tenants", tenantCtrl.ListTenants)
		superAdmin.Post("/tenants", tenantCtrl.CreateTenant)
		superAdmin.Put("/tenants/:id", tenantCtrl.UpdateTenant)
		superAdmin.Patch("/tenants/:id/status", tenantCtrl.SetStatus)
		superAdmin.Patch("/tenants/:id/features", tenantCtrl.ToggleFeature)

		// Tenant User Management Routes
		tenantUserScope := api.Group("/tenant", jwtAuth, tenantScope)
		tenantUserScope.Get("/members", tenantCtrl.GetMembers)
		tenantUserScope.Post("/members", tenantCtrl.InviteMember)
		tenantUserScope.Put("/members/:id/role", tenantCtrl.UpdateMemberRole)

		// Products Management Routes (Tenant Scoped)
		prodRoutes := api.Group("/products", jwtAuth, tenantScope)
		prodRoutes.Get("/", productCtrl.ListProducts)
		prodRoutes.Post("/", productCtrl.CreateProduct)
		prodRoutes.Get("/meta/brands", productCtrl.GetBrands)
		prodRoutes.Get("/meta/categories", productCtrl.GetCategories)
		prodRoutes.Get("/barcode/:barcode", productCtrl.LookupBarcode)
		prodRoutes.Get("/:id", productCtrl.GetProduct)
		prodRoutes.Delete("/:id", productCtrl.ArchiveProduct)

		// Inventory Management Routes (Tenant Scoped)
		invRoutes := api.Group("/inventory", jwtAuth, tenantScope)
		invRoutes.Get("/summary", productCtrl.GetInventorySummary)
		invRoutes.Get("/movements", productCtrl.GetStockMovements)
		invRoutes.Post("/adjustments", productCtrl.RecordStockAdjustment)

		// Sales & POS Data Import Routes (Tenant Scoped)
		salesRoutes := api.Group("/sales", jwtAuth, tenantScope)
		salesRoutes.Get("/", saleCtrl.ListSales)
		salesRoutes.Get("/import/history", saleCtrl.GetImportHistory)
		salesRoutes.Post("/preview", saleCtrl.PreviewImport)
		salesRoutes.Post("/import", saleCtrl.ExecuteImport)
		salesRoutes.Get("/:id", saleCtrl.GetSale)

		// Analytics Routes (Tenant Scoped)
		analyticsRoutes := api.Group("/analytics", jwtAuth, tenantScope)
		analyticsRoutes.Get("/overview", analyticsCtrl.GetOverview)
		analyticsRoutes.Get("/trends", analyticsCtrl.GetTrends)
		analyticsRoutes.Get("/payment-methods", analyticsCtrl.GetPaymentMethods)
		analyticsRoutes.Get("/products", analyticsCtrl.GetTopProducts)
		analyticsRoutes.Get("/categories", analyticsCtrl.GetCategorySales)
		analyticsRoutes.Get("/footwear", analyticsCtrl.GetFootwearIntelligence)
	}
}
