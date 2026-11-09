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

	// If DB is initialized, mount Auth and Tenant management routes
	if dbManager != nil && dbManager.DB != nil {
		authRepo := repositories.NewAuthRepository(dbManager.DB)
		tenantRepo := repositories.NewTenantRepository(dbManager.DB)

		authService := services.NewAuthService(authRepo, tenantRepo, cfg)
		tenantService := services.NewTenantService(tenantRepo, authRepo)

		authCtrl := controllers.NewAuthController(authService)
		tenantCtrl := controllers.NewTenantController(tenantService)

		jwtAuth := middleware.JWTMiddleware(cfg.DatabaseURL)

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

		// Tenant-Scoped Management Routes
		tenantScope := api.Group("/tenant", jwtAuth, middleware.RequireTenantContext())
		tenantScope.Get("/members", tenantCtrl.GetMembers)
		tenantScope.Post("/members", tenantCtrl.InviteMember)
		tenantScope.Put("/members/:id/role", tenantCtrl.UpdateMemberRole)
	}
}
