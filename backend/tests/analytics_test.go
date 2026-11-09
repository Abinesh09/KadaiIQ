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

func TestAdvancedSalesAnalyticsEngine(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// 1. Register Tenant A (Chennai Footwear)
	tenantAPayload := dto.RegisterTenantRequest{
		OrganizationName: "Chennai Footwear",
		Slug:             "chennai-footwear",
		OwnerName:        "Ramesh Kumar",
		Email:            "ramesh@chennaifootwear.com",
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

	// 2. Register Tenant B (Coimbatore Shoes) for strict isolation testing
	tenantBPayload := dto.RegisterTenantRequest{
		OrganizationName: "Coimbatore Shoes",
		Slug:             "coimbatore-shoes",
		OwnerName:        "Prakash Raj",
		Email:            "prakash@coimbatoreshoes.com",
		Password:         "Password@123",
	}
	bodyB, _ := json.Marshal(tenantBPayload)
	reqB := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(bodyB))
	reqB.Header.Set("Content-Type", "application/json")
	respB, err := app.Test(reqB, 1000)
	if err != nil || respB.StatusCode != http.StatusCreated {
		t.Fatalf("Tenant B registration failed: %v", err)
	}

	var authB struct {
		Data dto.AuthTokenResponse `json:"data"`
	}
	json.NewDecoder(respB.Body).Decode(&authB)
	tokenB := authB.Data.AccessToken

	// 3. Create Product in Tenant A with 2 variants
	prodPayload := dto.CreateProductRequest{
		Name:        "Bata Power Walker",
		ProductCode: "BPW-001",
		Gender:      "Men",
		Description: "Comfort footwear",
		Variants: []dto.CreateVariantRequest{
			{
				Size:              "UK 8",
				Colour:            "Black",
				Barcode:           "890100000008",
				SKU:               "BPW-BLK-8",
				PurchasePrice:     decimal.NewFromInt(1000),
				SellingPrice:      decimal.NewFromInt(2000),
				GSTRate:           decimal.NewFromInt(12),
				OpeningStock:      50,
				MinStockThreshold: 5,
			},
			{
				Size:              "UK 9",
				Colour:            "Blue",
				Barcode:           "890100000009",
				SKU:               "BPW-BLU-9",
				PurchasePrice:     decimal.NewFromInt(1200),
				SellingPrice:      decimal.NewFromInt(2500),
				GSTRate:           decimal.NewFromInt(12),
				OpeningStock:      50,
				MinStockThreshold: 5,
			},
		},
	}
	prodBody, _ := json.Marshal(prodPayload)
	reqProd := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(prodBody))
	reqProd.Header.Set("Content-Type", "application/json")
	reqProd.Header.Set("Authorization", "Bearer "+tokenA)
	respProd, err := app.Test(reqProd, 1000)
	if err != nil || respProd.StatusCode != http.StatusCreated {
		t.Fatalf("Product creation failed: %v", err)
	}

	// 4. Import POS Sales into Tenant A
	// Bill 1: 1 unit of UK 8 (Gross: 2000, Disc: 100, Net: 1900, Tax: 228, Grand Total: 2128, Payment: UPI)
	// Bill 2: 2 units of UK 9 (Gross: 5000, Disc: 200, Net: 4800, Tax: 576, Grand Total: 5376, Payment: Card)
	csvContentA := "Bill Number,Date,Barcode,Product Name,Quantity,Unit Price,Discount,Tax Amount,Total,Payment Mode\n" +
		"BILL-A-101,2026-11-09,890100000008,Bata Power Walker,1,2000,100,228,2128,UPI\n" +
		"BILL-A-102,2026-11-09,890100000009,Bata Power Walker,2,2500,200,576,5376,Card\n"

	importReqA := dto.ExecuteImportRequest{
		SourceSystem:    "ZenBiller",
		FileName:        "daily_sales_a.csv",
		CSVData:         csvContentA,
		DeductInventory: true,
	}
	importBodyA, _ := json.Marshal(importReqA)
	reqImpA := httptest.NewRequest("POST", "/api/v1/sales/import", bytes.NewReader(importBodyA))
	reqImpA.Header.Set("Content-Type", "application/json")
	reqImpA.Header.Set("Authorization", "Bearer "+tokenA)
	respImpA, err := app.Test(reqImpA, 1000)
	if err != nil || respImpA.StatusCode != http.StatusOK {
		t.Fatalf("Sales import for Tenant A failed: %v", err)
	}

	// 5. Import different sale into Tenant B
	csvContentB := "Bill Number,Date,Barcode,Product Name,Quantity,Unit Price,Discount,Tax Amount,Total,Payment Mode\n" +
		"BILL-B-999,2026-11-09,890999999999,Other Brand Shoe,5,1000,0,0,5000,Cash\n"

	importReqB := dto.ExecuteImportRequest{
		SourceSystem:    "ZenBiller",
		FileName:        "daily_sales_b.csv",
		CSVData:         csvContentB,
		DeductInventory: false,
	}
	importBodyB, _ := json.Marshal(importReqB)
	reqImpB := httptest.NewRequest("POST", "/api/v1/sales/import", bytes.NewReader(importBodyB))
	reqImpB.Header.Set("Content-Type", "application/json")
	reqImpB.Header.Set("Authorization", "Bearer "+tokenB)
	respImpB, err := app.Test(reqImpB, 1000)
	if err != nil || respImpB.StatusCode != http.StatusOK {
		t.Fatalf("Sales import for Tenant B failed: %v", err)
	}

	// TEST 1: Overview KPIs for Tenant A
	t.Run("Tenant A Overview KPIs", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/analytics/overview?range=30d", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Overview request failed: %v (status %d)", err, resp.StatusCode)
		}

		var apiResp struct {
			Data dto.AnalyticsOverviewResponse `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		overview := apiResp.Data

		if overview.TotalBills != 2 {
			t.Errorf("Expected 2 bills, got %d", overview.TotalBills)
		}
		if overview.TotalItemsSold != 3 {
			t.Errorf("Expected 3 items sold, got %d", overview.TotalItemsSold)
		}
		expectedGross := decimal.NewFromInt(7504) // 2128 + 5376
		if !overview.GrossSales.Equal(expectedGross) {
			t.Errorf("Expected gross sales %v, got %v", expectedGross, overview.GrossSales)
		}
		expectedNet := decimal.NewFromInt(7000) // 2000 + 5000
		if !overview.NetSales.Equal(expectedNet) {
			t.Errorf("Expected net sales %v, got %v", expectedNet, overview.NetSales)
		}
		expectedCost := decimal.NewFromInt(3400) // (1 * 1000) + (2 * 1200)
		if !overview.EstimatedCost.Equal(expectedCost) {
			t.Errorf("Expected cost %v, got %v", expectedCost, overview.EstimatedCost)
		}
		expectedProfit := decimal.NewFromInt(3600) // 7000 - 3400
		if !overview.GrossProfit.Equal(expectedProfit) {
			t.Errorf("Expected gross profit %v, got %v", expectedProfit, overview.GrossProfit)
		}
		if overview.AverageOrderVal.IsZero() {
			t.Errorf("Expected non-zero Average Order Value")
		}
	})

	// TEST 2: Daily Trends
	t.Run("Daily Sales Trends", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/analytics/trends?range=30d", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			var errResp map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errResp)
			t.Fatalf("Trends request failed (status %d): %v", resp.StatusCode, errResp)
		}

		var apiResp struct {
			Data []dto.TrendPoint `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		if len(apiResp.Data) == 0 {
			t.Errorf("Expected trend data points, got 0")
		}
	})

	// TEST 3: Payment Method Distribution
	t.Run("Payment Methods Breakdown", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/analytics/payment-methods?range=30d", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Payment methods request failed: %v", err)
		}

		var apiResp struct {
			Data []dto.PaymentMethodAnalytics `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		if len(apiResp.Data) != 2 {
			t.Errorf("Expected 2 payment methods, got %d", len(apiResp.Data))
		}
	})

	// TEST 4: Top Products Performance
	t.Run("Top Products Analytics", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/analytics/products?range=30d&limit=5", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Top products request failed: %v", err)
		}

		var apiResp struct {
			Data []dto.ProductSalesAnalytics `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		if len(apiResp.Data) == 0 {
			t.Fatalf("Expected top products, got 0")
		}
		if apiResp.Data[0].ProductName != "Bata Power Walker" {
			t.Errorf("Expected product 'Bata Power Walker', got '%s'", apiResp.Data[0].ProductName)
		}
		if apiResp.Data[0].UnitsSold != 3 {
			t.Errorf("Expected 3 units sold, got %d", apiResp.Data[0].UnitsSold)
		}
	})

	// TEST 5: Footwear Intelligence (Sizes, Colours, Brands)
	t.Run("Footwear Variant Intelligence", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/analytics/footwear?range=30d", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Footwear intelligence request failed: %v", err)
		}

		var apiResp struct {
			Data dto.FootwearIntelligenceResponse `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		intel := apiResp.Data

		if len(intel.TopSizes) != 2 {
			t.Errorf("Expected 2 sizes, got %d", len(intel.TopSizes))
		}
		// UK 9 had 2 units sold, UK 8 had 1 unit sold
		if intel.TopSizes[0].AttributeVal != "UK 9" || intel.TopSizes[0].UnitsSold != 2 {
			t.Errorf("Expected UK 9 to be top size with 2 units, got %s (%d units)", intel.TopSizes[0].AttributeVal, intel.TopSizes[0].UnitsSold)
		}
		if len(intel.TopColours) != 2 {
			t.Errorf("Expected 2 colours, got %d", len(intel.TopColours))
		}
		if len(intel.TopBrands) == 0 {
			t.Errorf("Expected brand intelligence data")
		}
	})

	// TEST 6: Strict Multi-Tenant Isolation
	t.Run("Strict Multi-Tenant Isolation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/analytics/overview?range=30d", nil)
		req.Header.Set("Authorization", "Bearer "+tokenB)
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Tenant B overview request failed: %v", err)
		}

		var apiResp struct {
			Data dto.AnalyticsOverviewResponse `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		overviewB := apiResp.Data

		// Tenant B must see exactly 1 bill and 5000 gross sales, strictly isolated from Tenant A's 7504
		if overviewB.TotalBills != 1 {
			t.Errorf("Tenant B expected 1 bill, got %d", overviewB.TotalBills)
		}
		expectedGrossB := decimal.NewFromInt(5000)
		if !overviewB.GrossSales.Equal(expectedGrossB) {
			t.Errorf("Tenant B expected gross %v, got %v", expectedGrossB, overviewB.GrossSales)
		}
		fmt.Printf("[Test] Strict tenant isolation verified: Tenant A: 7504, Tenant B: 5000\n")
	})
}
