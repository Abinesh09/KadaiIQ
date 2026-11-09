package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/config"
	"github.com/Abinesh09/KadaiIQ/backend/database"
	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/internal/routes"
	"github.com/Abinesh09/KadaiIQ/backend/internal/utils"
	"github.com/Abinesh09/KadaiIQ/backend/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func setupTestApp(t *testing.T) (*fiber.App, *config.Config, *database.DBManager) {
	db := setupTestDB(t)
	if err := migrations.RunMigrations(db); err != nil {
		t.Fatalf("Failed to execute migrations: %v", err)
	}

	cfg := &config.Config{
		Port:        "8080",
		AppVersion:  "1.0.0",
		Environment: "test",
		DatabaseURL: "test_jwt_secret_key_1234567890",
	}

	dbManager, err := database.NewDBManagerFromDB(db, cfg)
	if err != nil {
		t.Fatalf("Failed to initialize test DBManager: %v", err)
	}

	app := fiber.New()
	routes.SetupRoutes(app, cfg, dbManager)

	return app, cfg, dbManager
}

func TestAuthFlowAndTenantOnboarding(t *testing.T) {
	app, _, dbManager := setupTestApp(t)

	// 1. Test Super Admin login with seeded credentials
	loginPayload := dto.LoginRequest{
		Email:    "superadmin@kadaiiq.com",
		Password: "Admin@KadaiIQ2026",
	}
	body, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 1000)
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK for super admin login, got %d", resp.StatusCode)
	}

	var loginResp struct {
		Success bool                  `json:"success"`
		Data    dto.AuthTokenResponse `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResp)
	superAdminToken := loginResp.Data.AccessToken
	refreshToken := loginResp.Data.RefreshToken

	if superAdminToken == "" || refreshToken == "" {
		t.Fatal("Expected access and refresh tokens to be generated")
	}

	// 2. Test Invalid Credentials
	invalidReq := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader([]byte(`{"email":"superadmin@kadaiiq.com","password":"WrongPassword"}`)))
	invalidReq.Header.Set("Content-Type", "application/json")
	respInv, _ := app.Test(invalidReq, 1000)
	if respInv.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong password, got %d", respInv.StatusCode)
	}

	// 3. Test Token Refresh & Rotation
	refreshReq := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader([]byte(`{"refresh_token":"`+refreshToken+`"}`)))
	refreshReq.Header.Set("Content-Type", "application/json")
	respRef, _ := app.Test(refreshReq, 1000)
	if respRef.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 OK for token refresh, got %d", respRef.StatusCode)
	}

	// 4. Test Single-Use: Reusing old refresh token must be rejected
	respReused, _ := app.Test(refreshReq, 1000)
	if respReused.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized for reused/rotated refresh token, got %d", respReused.StatusCode)
	}

	// 5. Test Super Admin creating an organisation
	newTenantPayload := dto.CreateTenantAdminRequest{
		Name:          "Kovai Retail Mart",
		Slug:          "kovai-retail",
		BusinessName:  "Kovai Mart Retailers Pvt Ltd",
		BusinessType:  "Footwear & Apparel",
		GSTIN:         "33AAAAA0000A1Z5",
		Email:         "contact@kovaimart.com",
		OwnerName:     "Karthik Selvan",
		OwnerEmail:    "karthik@kovaimart.com",
		OwnerPassword: "Password@123",
	}
	tenantBody, _ := json.Marshal(newTenantPayload)
	createReq := httptest.NewRequest("POST", "/api/v1/super-admin/tenants", bytes.NewReader(tenantBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+superAdminToken)
	respTenant, err := app.Test(createReq, 1000)
	if err != nil || respTenant.StatusCode != http.StatusCreated {
		t.Fatalf("Super Admin create tenant failed with code %d", respTenant.StatusCode)
	}

	// 6. Test Tenant Owner login
	ownerLoginReq := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader([]byte(`{"email":"karthik@kovaimart.com","password":"Password@123"}`)))
	ownerLoginReq.Header.Set("Content-Type", "application/json")
	respOwner, _ := app.Test(ownerLoginReq, 1000)
	if respOwner.StatusCode != http.StatusOK {
		t.Fatalf("Tenant Owner login failed with code %d", respOwner.StatusCode)
	}

	var ownerResp struct {
		Data dto.AuthTokenResponse `json:"data"`
	}
	json.NewDecoder(respOwner.Body).Decode(&ownerResp)
	ownerToken := ownerResp.Data.AccessToken

	// 7. Security Check: Tenant Owner cannot access Super Admin endpoint
	unauthReq := httptest.NewRequest("GET", "/api/v1/super-admin/tenants", nil)
	unauthReq.Header.Set("Authorization", "Bearer "+ownerToken)
	respUnauth, _ := app.Test(unauthReq, 1000)
	if respUnauth.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden when regular tenant accesses super-admin route, got %d", respUnauth.StatusCode)
	}

	// 8. Test Tenant Member Invite
	var staffRole models.Role
	dbManager.DB.Where("name = ?", "Staff").First(&staffRole)

	invitePayload := dto.InviteUserRequest{
		Name:   "Priya Cashier",
		Email:  "priya@kovaimart.com",
		RoleID: staffRole.ID,
	}
	inviteBody, _ := json.Marshal(invitePayload)
	inviteReq := httptest.NewRequest("POST", "/api/v1/tenant/members", bytes.NewReader(inviteBody))
	inviteReq.Header.Set("Content-Type", "application/json")
	inviteReq.Header.Set("Authorization", "Bearer "+ownerToken)
	respInvite, _ := app.Test(inviteReq, 1000)
	if respInvite.StatusCode != http.StatusCreated {
		t.Fatalf("Tenant member invite failed with code %d", respInvite.StatusCode)
	}

	// 9. Test Password Reset Flow (expiring & single-use)
	forgotReq := httptest.NewRequest("POST", "/api/v1/auth/forgot-password", bytes.NewReader([]byte(`{"email":"karthik@kovaimart.com"}`)))
	forgotReq.Header.Set("Content-Type", "application/json")
	respForgot, _ := app.Test(forgotReq, 1000)
	if respForgot.StatusCode != http.StatusOK {
		t.Fatalf("Forgot password request failed: %d", respForgot.StatusCode)
	}

	var forgotResp struct {
		Data struct {
			ResetToken string `json:"reset_token"`
		} `json:"data"`
	}
	json.NewDecoder(respForgot.Body).Decode(&forgotResp)
	resetToken := forgotResp.Data.ResetToken

	if resetToken == "" {
		t.Fatal("Expected reset token to be generated")
	}

	// Execute password reset
	resetPayload := dto.ResetPasswordRequest{
		Token:       resetToken,
		NewPassword: "NewSecretPassword@456",
	}
	resetBody, _ := json.Marshal(resetPayload)
	resetReq := httptest.NewRequest("POST", "/api/v1/auth/reset-password", bytes.NewReader(resetBody))
	resetReq.Header.Set("Content-Type", "application/json")
	respReset, _ := app.Test(resetReq, 1000)
	if respReset.StatusCode != http.StatusOK {
		t.Fatalf("Password reset failed: %d", respReset.StatusCode)
	}

	// Test Single-Use of reset token
	respResetReuse, _ := app.Test(resetReq, 1000)
	if respResetReuse.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request when attempting to reuse password reset token, got %d", respResetReuse.StatusCode)
	}

	// Verify login with new password
	newLoginReq := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader([]byte(`{"email":"karthik@kovaimart.com","password":"NewSecretPassword@456"}`)))
	newLoginReq.Header.Set("Content-Type", "application/json")
	respNewLogin, _ := app.Test(newLoginReq, 1000)
	if respNewLogin.StatusCode != http.StatusOK {
		t.Fatalf("Login with new password failed: %d", respNewLogin.StatusCode)
	}
}

func TestTokenExpiration(t *testing.T) {
	secret := "test_secret_expired"
	// Generate token expired 1 minute ago
	expiredToken, err := utils.GenerateAccessToken(
		uuid.New(),
		"test@test.com",
		uuid.New(),
		uuid.New(),
		"Staff",
		false,
		[]string{"products:view"},
		secret,
		-1*time.Minute,
	)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	_, validateErr := utils.ValidateAccessToken(expiredToken, secret)
	if validateErr == nil {
		t.Error("Expected validation error for expired token, but got nil")
	}
}
