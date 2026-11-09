package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/shopspring/decimal"
)

func TestProductAndInventoryLifecycle(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// Step 1: Register Tenant A (Chennai Shoes)
	tenantAPayload := dto.RegisterTenantRequest{
		OrganizationName: "Chennai Shoes",
		Slug:             "chennai-shoes",
		OwnerName:        "Anand Kumar",
		Email:            "anand@chennaishoes.com",
		Password:         "Password@123",
	}
	bodyA, _ := json.Marshal(tenantAPayload)
	reqA := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(bodyA))
	reqA.Header.Set("Content-Type", "application/json")
	respA, err := app.Test(reqA, 1000)
	if err != nil || respA.StatusCode != http.StatusCreated {
		t.Fatalf("Tenant A registration failed: %v", err)
	}

	var authA struct {
		Data dto.AuthTokenResponse `json:"data"`
	}
	json.NewDecoder(respA.Body).Decode(&authA)
	tokenA := authA.Data.AccessToken

	// Step 2: Register Tenant B (Madurai Footwear)
	tenantBPayload := dto.RegisterTenantRequest{
		OrganizationName: "Madurai Footwear",
		Slug:             "madurai-footwear",
		OwnerName:        "Bala Murugan",
		Email:            "bala@maduraifootwear.com",
		Password:         "Password@123",
	}
	bodyB, _ := json.Marshal(tenantBPayload)
	reqB := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(bodyB))
	reqB.Header.Set("Content-Type", "application/json")
	respB, _ := app.Test(reqB, 1000)
	var authB struct {
		Data dto.AuthTokenResponse `json:"data"`
	}
	json.NewDecoder(respB.Body).Decode(&authB)
	tokenB := authB.Data.AccessToken

	// Step 3: Tenant A creates a footwear product with variants (sizes 7 & 8) and opening stock
	createProdPayload := dto.CreateProductRequest{
		Name:        "Air Stride Pro Running Shoes",
		ProductCode: "ASP-001",
		Gender:      "Men",
		Description: "High-performance mesh running footwear",
		Variants: []dto.CreateVariantRequest{
			{
				Size:              "UK 7",
				Colour:            "Black/Red",
				SKU:               "ASP-BLK-07",
				Barcode:           "890123456701",
				PurchasePrice:     decimal.NewFromFloat(1500.00),
				SellingPrice:      decimal.NewFromFloat(2999.00),
				GSTRate:           decimal.NewFromFloat(12.00),
				OpeningStock:      20,
				MinStockThreshold: 5,
			},
			{
				Size:              "UK 8",
				Colour:            "Black/Red",
				SKU:               "ASP-BLK-08",
				Barcode:           "890123456702",
				PurchasePrice:     decimal.NewFromFloat(1500.00),
				SellingPrice:      decimal.NewFromFloat(2999.00),
				GSTRate:           decimal.NewFromFloat(12.00),
				OpeningStock:      15,
				MinStockThreshold: 5,
			},
		},
	}
	prodBody, _ := json.Marshal(createProdPayload)
	createReq := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(prodBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+tokenA)
	createResp, err := app.Test(createReq, 1000)
	if err != nil || createResp.StatusCode != http.StatusCreated {
		t.Fatalf("Product creation failed with code: %d", createResp.StatusCode)
	}

	var prodResp struct {
		Data dto.ProductResponse `json:"data"`
	}
	json.NewDecoder(createResp.Body).Decode(&prodResp)
	createdProduct := prodResp.Data

	if createdProduct.TotalStock != 35 {
		t.Errorf("Expected total initial stock to be 35, got %d", createdProduct.TotalStock)
	}
	if len(createdProduct.Variants) != 2 {
		t.Fatalf("Expected 2 variants, got %d", len(createdProduct.Variants))
	}
	variant1ID := createdProduct.Variants[0].ID

	// Step 4: Test Barcode Duplication prevention within Tenant A
	dupProdPayload := createProdPayload
	dupProdPayload.ProductCode = "ASP-002"
	dupBody, _ := json.Marshal(dupProdPayload)
	dupReq := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(dupBody))
	dupReq.Header.Set("Content-Type", "application/json")
	dupReq.Header.Set("Authorization", "Bearer "+tokenA)
	dupResp, _ := app.Test(dupReq, 1000)
	if dupResp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request when attempting to create duplicate barcode in same tenant, got %d", dupResp.StatusCode)
	}

	// Step 5: Test Barcode lookup by scanner
	lookupReq := httptest.NewRequest("GET", "/api/v1/products/barcode/890123456701", nil)
	lookupReq.Header.Set("Authorization", "Bearer "+tokenA)
	lookupResp, _ := app.Test(lookupReq, 1000)
	if lookupResp.StatusCode != http.StatusOK {
		t.Fatalf("Barcode lookup failed with code: %d", lookupResp.StatusCode)
	}

	// Step 6: Test Tenant Isolation - Tenant B scanning Tenant A's barcode must return 404 Not Found
	crossLookupReq := httptest.NewRequest("GET", "/api/v1/products/barcode/890123456701", nil)
	crossLookupReq.Header.Set("Authorization", "Bearer "+tokenB)
	crossResp, _ := app.Test(crossLookupReq, 1000)
	if crossResp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found when Tenant B looks up Tenant A barcode, got %d", crossResp.StatusCode)
	}

	// Step 7: Record a Purchase Restock of 10 units for Variant 1
	restockPayload := dto.StockAdjustmentRequest{
		VariantID:       variant1ID,
		MovementType:    "PURCHASE",
		Quantity:        10,
		UnitCost:        decimal.NewFromFloat(1450.00),
		ReferenceNumber: "PO-2026-001",
		Notes:           "Bulk supplier restock",
	}
	restockBody, _ := json.Marshal(restockPayload)
	restockReq := httptest.NewRequest("POST", "/api/v1/inventory/adjustments", bytes.NewReader(restockBody))
	restockReq.Header.Set("Content-Type", "application/json")
	restockReq.Header.Set("Authorization", "Bearer "+tokenA)
	restockResp, _ := app.Test(restockReq, 1000)
	if restockResp.StatusCode != http.StatusOK {
		t.Fatalf("Restock adjustment failed: %d", restockResp.StatusCode)
	}

	var moveResp struct {
		Data dto.StockMovementResponse `json:"data"`
	}
	json.NewDecoder(restockResp.Body).Decode(&moveResp)
	if moveResp.Data.NewStock != 30 {
		t.Errorf("Expected Variant 1 stock to be 30 after restock (20 + 10), got %d", moveResp.Data.NewStock)
	}

	// Step 8: Record Sale Deduction of 5 units
	salePayload := dto.StockAdjustmentRequest{
		VariantID:       variant1ID,
		MovementType:    "SALE",
		Quantity:        5,
		ReferenceNumber: "BILL-1001",
		Notes:           "Walk-in POS purchase",
	}
	saleBody, _ := json.Marshal(salePayload)
	saleReq := httptest.NewRequest("POST", "/api/v1/inventory/adjustments", bytes.NewReader(saleBody))
	saleReq.Header.Set("Content-Type", "application/json")
	saleReq.Header.Set("Authorization", "Bearer "+tokenA)
	saleResp, _ := app.Test(saleReq, 1000)
	if saleResp.StatusCode != http.StatusOK {
		t.Fatalf("Sale stock deduction failed: %d", saleResp.StatusCode)
	}

	// Step 9: Prevent Overselling (e.g. attempting to sell 50 units when only 25 remain)
	oversellPayload := dto.StockAdjustmentRequest{
		VariantID:    variant1ID,
		MovementType: "SALE",
		Quantity:     50,
	}
	oversellBody, _ := json.Marshal(oversellPayload)
	oversellReq := httptest.NewRequest("POST", "/api/v1/inventory/adjustments", bytes.NewReader(oversellBody))
	oversellReq.Header.Set("Content-Type", "application/json")
	oversellReq.Header.Set("Authorization", "Bearer "+tokenA)
	oversellResp, _ := app.Test(oversellReq, 1000)
	if oversellResp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request when selling exceeds available stock, got %d", oversellResp.StatusCode)
	}

	// Step 10: Check Inventory Summary
	sumReq := httptest.NewRequest("GET", "/api/v1/inventory/summary", nil)
	sumReq.Header.Set("Authorization", "Bearer "+tokenA)
	sumResp, _ := app.Test(sumReq, 1000)
	if sumResp.StatusCode != http.StatusOK {
		t.Fatalf("Inventory summary failed: %d", sumResp.StatusCode)
	}

	var sumData struct {
		Data dto.InventorySummaryResponse `json:"data"`
	}
	json.NewDecoder(sumResp.Body).Decode(&sumData)
	summary := sumData.Data

	if summary.TotalProducts != 1 {
		t.Errorf("Expected 1 total product, got %d", summary.TotalProducts)
	}
	if summary.TotalVariants != 2 {
		t.Errorf("Expected 2 total variants, got %d", summary.TotalVariants)
	}
	// Variant 1: 25 units. Variant 2: 15 units. Total: 40 units
	if summary.TotalStockUnits != 40 {
		t.Errorf("Expected 40 total units, got %d", summary.TotalStockUnits)
	}

	// Step 11: Check Stock Movements Ledger
	moveReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/inventory/movements?variant_id=%s", variant1ID), nil)
	moveReq.Header.Set("Authorization", "Bearer "+tokenA)
	ledgerResp, _ := app.Test(moveReq, 1000)
	if ledgerResp.StatusCode != http.StatusOK {
		t.Fatalf("Movements query failed: %d", ledgerResp.StatusCode)
	}

	var movementsResp struct {
		Data struct {
			Items []dto.StockMovementResponse `json:"items"`
		} `json:"data"`
	}
	json.NewDecoder(ledgerResp.Body).Decode(&movementsResp)
	// Variant 1 should have 3 recorded movements: OPENING, PURCHASE, SALE
	if len(movementsResp.Data.Items) != 3 {
		t.Errorf("Expected 3 movements for Variant 1, got %d", len(movementsResp.Data.Items))
	}
}
