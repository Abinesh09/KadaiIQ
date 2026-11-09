package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Abinesh09/KadaiIQ/backend/config"
	"github.com/Abinesh09/KadaiIQ/backend/database"
	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func TestHealthCheck(t *testing.T) {
	cfg := &config.Config{
		Port:        "8080",
		AppVersion:  "1.0.0",
		Environment: "test",
	}

	app := fiber.New()
	routes.SetupRoutes(app, cfg, nil)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	resp, err := app.Test(req, 1000)
	if err != nil {
		t.Fatalf("Failed to execute test request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	var apiResp dto.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if !apiResp.Success {
		t.Errorf("Expected success to be true, got %v", apiResp.Success)
	}

	if apiResp.Message != "Service is healthy" {
		t.Errorf("Unexpected message: %s", apiResp.Message)
	}
}

func TestHealthCheckDBUnavailable(t *testing.T) {
	cfg := &config.Config{
		Port:        "8080",
		AppVersion:  "1.0.0",
		Environment: "test",
	}

	app := fiber.New()
	routes.SetupRoutes(app, cfg, nil)

	req := httptest.NewRequest("GET", "/api/v1/health/db", nil)
	resp, err := app.Test(req, 1000)
	if err != nil {
		t.Fatalf("Failed to execute test request: %v", err)
	}

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("Expected status code 503, got %d", resp.StatusCode)
	}

	var apiResp dto.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if apiResp.Success {
		t.Errorf("Expected success to be false when DB is unavailable, got %v", apiResp.Success)
	}
}

func TestHealthCheckDBHealthy(t *testing.T) {
	cfg := &config.Config{
		Port:        "8080",
		AppVersion:  "1.0.0",
		Environment: "test",
	}

	db := setupTestDB(t)
	dbManager, err := database.NewDBManagerFromDB(db, cfg)
	if err != nil {
		t.Fatalf("Failed to create dbManager: %v", err)
	}

	app := fiber.New()
	routes.SetupRoutes(app, cfg, dbManager)

	req := httptest.NewRequest("GET", "/api/v1/health/db", nil)
	resp, err := app.Test(req, 1000)
	if err != nil {
		t.Fatalf("Failed to execute test request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	var apiResp dto.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if !apiResp.Success {
		t.Errorf("Expected success to be true when DB is healthy, got %v", apiResp.Success)
	}
	if apiResp.Message != "Database is healthy" {
		t.Errorf("Expected message 'Database is healthy', got '%s'", apiResp.Message)
	}
}
