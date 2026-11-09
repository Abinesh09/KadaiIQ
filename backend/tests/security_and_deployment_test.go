package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
)

func TestSecurityAndDeploymentHardening(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// TEST 1: OWASP Security Headers & CORS
	t.Run("OWASP Security Headers and CORS Enforcement", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/health", nil)
		req.Header.Set("Origin", "http://localhost:7070")
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Health request failed: %v", err)
		}

		// Verify OWASP Security Headers
		if resp.Header.Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("Expected X-Content-Type-Options: nosniff, got %s", resp.Header.Get("X-Content-Type-Options"))
		}
		if resp.Header.Get("X-Frame-Options") != "DENY" {
			t.Errorf("Expected X-Frame-Options: DENY, got %s", resp.Header.Get("X-Frame-Options"))
		}
		if resp.Header.Get("X-XSS-Protection") != "1; mode=block" {
			t.Errorf("Expected X-XSS-Protection: 1; mode=block, got %s", resp.Header.Get("X-XSS-Protection"))
		}
		if resp.Header.Get("Referrer-Policy") != "strict-origin-when-cross-origin" {
			t.Errorf("Expected Referrer-Policy: strict-origin-when-cross-origin, got %s", resp.Header.Get("Referrer-Policy"))
		}

		// Verify CORS for Frontend Port 7070
		if resp.Header.Get("Access-Control-Allow-Origin") != "http://localhost:7070" {
			t.Errorf("Expected CORS origin http://localhost:7070, got %s", resp.Header.Get("Access-Control-Allow-Origin"))
		}
	})

	// TEST 2: Protected Endpoints Deny Access Without Token (401 Unauthorized)
	t.Run("Protected Endpoints Require Authentication", func(t *testing.T) {
		endpoints := []string{
			"/api/v1/products",
			"/api/v1/inventory/summary",
			"/api/v1/sales",
			"/api/v1/analytics/overview",
			"/api/v1/gst/settings",
			"/api/v1/reports/sales",
		}

		for _, ep := range endpoints {
			req := httptest.NewRequest("GET", ep, nil)
			resp, err := app.Test(req, 1000)
			if err != nil {
				t.Fatalf("Request to %s failed: %v", ep, err)
			}
			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("Expected 401 for unauthenticated %s, got %d", ep, resp.StatusCode)
			}
		}
	})

	// TEST 3: Super Admin Access Control (Tenant Owner cannot access platform management)
	t.Run("Super Admin Route Role Protection", func(t *testing.T) {
		// Register a normal Tenant Owner
		regPayload := dto.RegisterTenantRequest{
			OrganizationName: "Security Test Store",
			Slug:             "security-test-store",
			OwnerName:        "Store Admin",
			Email:            "admin@securitystore.com",
			Password:         "Password@123",
		}
		body, _ := json.Marshal(regPayload)
		regReq := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
		regReq.Header.Set("Content-Type", "application/json")
		regResp, err := app.Test(regReq, 1000)
		if err != nil || regResp.StatusCode != http.StatusCreated {
			t.Fatalf("Tenant registration failed: %v", err)
		}

		var authResp struct {
			Data dto.AuthTokenResponse `json:"data"`
		}
		json.NewDecoder(regResp.Body).Decode(&authResp)
		tenantToken := authResp.Data.AccessToken

		// Tenant Owner attempts to access Super Admin platform route -> Must be 403 Forbidden!
		saReq := httptest.NewRequest("GET", "/api/v1/super-admin/tenants", nil)
		saReq.Header.Set("Authorization", "Bearer "+tenantToken)
		saResp, err := app.Test(saReq, 1000)
		if err != nil {
			t.Fatalf("Super Admin request failed: %v", err)
		}
		if saResp.StatusCode != http.StatusForbidden {
			t.Errorf("Expected 403 Forbidden for non-superadmin user, got %d", saResp.StatusCode)
		}
	})

	// TEST 4: Health Probes for Docker / Kubernetes Readiness and Liveness
	t.Run("Dual Health Check Probes", func(t *testing.T) {
		// Liveness probe
		reqLive := httptest.NewRequest("GET", "/api/v1/health", nil)
		respLive, err := app.Test(reqLive, 1000)
		if err != nil || respLive.StatusCode != http.StatusOK {
			t.Errorf("Liveness probe failed: %v", err)
		}

		// Readiness probe (Database connectivity & latency)
		reqReady := httptest.NewRequest("GET", "/api/v1/health/db", nil)
		respReady, err := app.Test(reqReady, 1000)
		if err != nil || respReady.StatusCode != http.StatusOK {
			t.Errorf("Readiness probe failed: %v", err)
		}
	})
}
