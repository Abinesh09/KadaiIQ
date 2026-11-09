package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/shopspring/decimal"
)

func TestGSTComplianceAndReports(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// Step 1: Register Tenant A (Chennai Footwear Mart)
	regA := dto.RegisterTenantRequest{
		OrganizationName: "Chennai Footwear Mart",
		Slug:             "chennai-footwear-mart",
		OwnerName:        "Murugan S",
		Email:            "murugan@chennaimart.com",
		Password:         "Password@123",
	}
	bodyA, _ := json.Marshal(regA)
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

	// Step 2: Register Tenant B (Bengaluru Shoe Plaza)
	regB := dto.RegisterTenantRequest{
		OrganizationName: "Bengaluru Shoe Plaza",
		Slug:             "bengaluru-shoe-plaza",
		OwnerName:        "Venkatesh R",
		Email:            "venkatesh@bengaluruplaza.com",
		Password:         "Password@123",
	}
	bodyB, _ := json.Marshal(regB)
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

	// Step 3: Seed Catalogue and Sales in Tenant A
	prodPayload := dto.CreateProductRequest{
		Name:        "VKC Pride Formal",
		ProductCode: "VKC-001",
		Gender:      "Men",
		Variants: []dto.CreateVariantRequest{
			{
				Size:              "UK 8",
				Colour:            "Black",
				Barcode:           "890200000008",
				SKU:               "VKC-BLK-8",
				PurchasePrice:     decimal.NewFromInt(600),
				SellingPrice:      decimal.NewFromInt(1200),
				GSTRate:           decimal.NewFromInt(12),
				OpeningStock:      40,
				MinStockThreshold: 5,
			},
		},
	}
	pBody, _ := json.Marshal(prodPayload)
	pReq := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(pBody))
	pReq.Header.Set("Content-Type", "application/json")
	pReq.Header.Set("Authorization", "Bearer "+tokenA)
	pResp, err := app.Test(pReq, 1000)
	if err != nil || pResp.StatusCode != http.StatusCreated {
		t.Fatalf("Product creation failed: %v", err)
	}

	// Import POS sale: 2 units at rate 1200, GST 12%
	// Taxable: 2400, GST: 288 (CGST: 144, SGST: 144), Total: 2688
	csvContent := "Bill Number,Date,Barcode,Product Name,Quantity,Unit Price,Discount,Tax Rate,Tax Amount,Total,Payment Mode\n" +
		"BILL-GST-101,2026-11-09,890200000008,VKC Pride Formal,2,1200,0,12,288,2688,UPI\n"

	impPayload := dto.ExecuteImportRequest{
		SourceSystem:    "ZenBiller",
		FileName:        "gst_sales.csv",
		CSVData:         csvContent,
		DeductInventory: true,
	}
	impBody, _ := json.Marshal(impPayload)
	impReq := httptest.NewRequest("POST", "/api/v1/sales/import", bytes.NewReader(impBody))
	impReq.Header.Set("Content-Type", "application/json")
	impReq.Header.Set("Authorization", "Bearer "+tokenA)
	impResp, err := app.Test(impReq, 1000)
	if err != nil || impResp.StatusCode != http.StatusOK {
		t.Fatalf("Sales import failed: %v", err)
	}

	// TEST 1: Retrieve and Update Tenant GST Settings
	t.Run("GST Settings and GSTIN Validation", func(t *testing.T) {
		// Get default settings
		getReq := httptest.NewRequest("GET", "/api/v1/gst/settings", nil)
		getReq.Header.Set("Authorization", "Bearer "+tokenA)
		getResp, err := app.Test(getReq, 1000)
		if err != nil || getResp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to get GST settings: %v", err)
		}

		var apiResp struct {
			Data dto.GSTSettingsResponse `json:"data"`
		}
		json.NewDecoder(getResp.Body).Decode(&apiResp)
		if apiResp.Data.StateCode != "33" {
			t.Errorf("Expected default state code 33, got %s", apiResp.Data.StateCode)
		}

		// Update with invalid GSTIN format (should fail 400)
		invalidPayload := dto.UpdateGSTSettingsRequest{
			GSTIN: "INVALID_GSTIN_123",
		}
		invBody, _ := json.Marshal(invalidPayload)
		invReq := httptest.NewRequest("PUT", "/api/v1/gst/settings", bytes.NewReader(invBody))
		invReq.Header.Set("Content-Type", "application/json")
		invReq.Header.Set("Authorization", "Bearer "+tokenA)
		invResp, _ := app.Test(invReq, 1000)
		if invResp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 for invalid GSTIN, got %d", invResp.StatusCode)
		}

		// Update with valid GSTIN (Karnataka code: 29)
		validPayload := dto.UpdateGSTSettingsRequest{
			GSTIN:              "29ABCDE1234F1Z5",
			LegalName:          "Chennai Footwear Enterprises Pvt Ltd",
			TradeName:          "Chennai Footwear Mart",
			DefaultFootwearHSN: "6403",
			DefaultGSTRate:     decimal.NewFromInt(12),
		}
		valBody, _ := json.Marshal(validPayload)
		valReq := httptest.NewRequest("PUT", "/api/v1/gst/settings", bytes.NewReader(valBody))
		valReq.Header.Set("Content-Type", "application/json")
		valReq.Header.Set("Authorization", "Bearer "+tokenA)
		valResp, err := app.Test(valReq, 1000)
		if err != nil || valResp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to update valid GSTIN: %v", err)
		}

		var updateResp struct {
			Data dto.GSTSettingsResponse `json:"data"`
		}
		json.NewDecoder(valResp.Body).Decode(&updateResp)
		if updateResp.Data.StateCode != "29" || updateResp.Data.StateName != "Karnataka" {
			t.Errorf("Expected state code 29 (Karnataka), got %s (%s)", updateResp.Data.StateCode, updateResp.Data.StateName)
		}
	})

	// TEST 2: GSTR-1 Return Filing Summary
	t.Run("GSTR-1 Monthly Return Summary", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/gst/gstr1?month=11&year=2026", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("GSTR-1 request failed: %v", err)
		}

		var apiResp struct {
			Data dto.GSTR1SummaryResponse `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		gstr1 := apiResp.Data

		// Taxable Value: 2400
		expectedTaxable := decimal.NewFromInt(2400)
		if !gstr1.TotalTaxableValue.Equal(expectedTaxable) {
			t.Errorf("Expected taxable value %v, got %v", expectedTaxable, gstr1.TotalTaxableValue)
		}

		// Total Tax: 288 (Split into CGST 144, SGST 144)
		expectedTax := decimal.NewFromInt(288)
		if !gstr1.TotalTaxLiability.Equal(expectedTax) {
			t.Errorf("Expected tax liability %v, got %v", expectedTax, gstr1.TotalTaxLiability)
		}
		expectedCGST := decimal.NewFromInt(144)
		if !gstr1.TotalCGST.Equal(expectedCGST) {
			t.Errorf("Expected CGST %v, got %v", expectedCGST, gstr1.TotalCGST)
		}
		expectedSGST := decimal.NewFromInt(144)
		if !gstr1.TotalSGST.Equal(expectedSGST) {
			t.Errorf("Expected SGST %v, got %v", expectedSGST, gstr1.TotalSGST)
		}

		// Rate Slab: 12%
		if len(gstr1.RateBreakdown) == 0 {
			t.Fatalf("Expected at least 1 rate slab")
		}
		if !gstr1.RateBreakdown[0].RatePct.Equal(decimal.NewFromInt(12)) {
			t.Errorf("Expected rate 12%%, got %v", gstr1.RateBreakdown[0].RatePct)
		}

		// HSN Summary: Footwear (6403)
		if len(gstr1.HSNSummary) == 0 {
			t.Fatalf("Expected HSN summary entry")
		}
		if gstr1.HSNSummary[0].TotalQty != 2 {
			t.Errorf("Expected HSN total quantity 2, got %d", gstr1.HSNSummary[0].TotalQty)
		}
	})

	// TEST 3: GSTR-3B Table 3.1
	t.Run("GSTR-3B Table 3.1 Outward Taxable Supplies", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/gst/gstr3b?month=11&year=2026", nil)
		req.Header.Set("Authorization", "Bearer "+tokenA)
		resp, err := app.Test(req, 1000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("GSTR-3B request failed: %v", err)
		}

		var apiResp struct {
			Data dto.GSTR3BTable31Response `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		t31 := apiResp.Data

		if !t31.TaxableValue.Equal(decimal.NewFromInt(2400)) {
			t.Errorf("Expected taxable value 2400, got %v", t31.TaxableValue)
		}
		if !t31.CentralTax.Equal(decimal.NewFromInt(144)) {
			t.Errorf("Expected central tax 144, got %v", t31.CentralTax)
		}
		if !t31.StateUTTax.Equal(decimal.NewFromInt(144)) {
			t.Errorf("Expected state tax 144, got %v", t31.StateUTTax)
		}
	})

	// TEST 4: Reports (JSON and CSV Streaming)
	t.Run("Reports JSON and CSV Exports", func(t *testing.T) {
		// 1. Sales Report JSON
		sReq := httptest.NewRequest("GET", "/api/v1/reports/sales", nil)
		sReq.Header.Set("Authorization", "Bearer "+tokenA)
		sResp, err := app.Test(sReq, 1000)
		if err != nil || sResp.StatusCode != http.StatusOK {
			t.Fatalf("Sales report JSON failed: %v", err)
		}

		var sReport struct {
			Data dto.ReportResponse[dto.SalesSummaryRow] `json:"data"`
		}
		json.NewDecoder(sResp.Body).Decode(&sReport)
		if sReport.Data.TotalRecords != 1 {
			t.Errorf("Expected 1 sales record, got %d", sReport.Data.TotalRecords)
		}

		// 2. Stock Valuation Report JSON
		vReq := httptest.NewRequest("GET", "/api/v1/reports/inventory", nil)
		vReq.Header.Set("Authorization", "Bearer "+tokenA)
		vResp, err := app.Test(vReq, 1000)
		if err != nil || vResp.StatusCode != http.StatusOK {
			t.Fatalf("Inventory report JSON failed: %v", err)
		}

		var vReport struct {
			Data dto.ReportResponse[dto.StockValuationRow] `json:"data"`
		}
		json.NewDecoder(vResp.Body).Decode(&vReport)
		if vReport.Data.TotalRecords == 0 {
			t.Errorf("Expected stock valuation records, got 0")
		}

		// 3. Stock Movement Report JSON
		mReq := httptest.NewRequest("GET", "/api/v1/reports/movements", nil)
		mReq.Header.Set("Authorization", "Bearer "+tokenA)
		mResp, err := app.Test(mReq, 1000)
		if err != nil || mResp.StatusCode != http.StatusOK {
			t.Fatalf("Stock movements report JSON failed: %v", err)
		}

		// 4. CSV Download Streaming
		csvReq := httptest.NewRequest("GET", "/api/v1/reports/download?type=sales", nil)
		csvReq.Header.Set("Authorization", "Bearer "+tokenA)
		csvResp, err := app.Test(csvReq, 1000)
		if err != nil || csvResp.StatusCode != http.StatusOK {
			t.Fatalf("CSV download failed: %v", err)
		}

		contentType := csvResp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/csv") {
			t.Errorf("Expected text/csv content type, got %s", contentType)
		}
		contentDisp := csvResp.Header.Get("Content-Disposition")
		if !strings.Contains(contentDisp, "attachment") {
			t.Errorf("Expected attachment header, got %s", contentDisp)
		}
	})

	// TEST 5: Strict Multi-Tenant Isolation
	t.Run("Strict Multi-Tenant Isolation for GST and Reports", func(t *testing.T) {
		// Tenant B has 0 sales. Querying GSTR-1 must yield 0 taxable value, completely isolated from Tenant A.
		bReq := httptest.NewRequest("GET", "/api/v1/gst/gstr1?month=11&year=2026", nil)
		bReq.Header.Set("Authorization", "Bearer "+tokenB)
		bResp, err := app.Test(bReq, 1000)
		if err != nil || bResp.StatusCode != http.StatusOK {
			t.Fatalf("Tenant B GSTR-1 request failed: %v", err)
		}

		var bApiResp struct {
			Data dto.GSTR1SummaryResponse `json:"data"`
		}
		json.NewDecoder(bResp.Body).Decode(&bApiResp)
		if !bApiResp.Data.TotalTaxableValue.IsZero() {
			t.Errorf("Tenant B expected 0 taxable value, got %v", bApiResp.Data.TotalTaxableValue)
		}
		if len(bApiResp.Data.RateBreakdown) != 0 {
			t.Errorf("Tenant B expected 0 rate slabs, got %d", len(bApiResp.Data.RateBreakdown))
		}
	})
}
