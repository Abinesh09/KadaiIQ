package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/shopspring/decimal"
)

func TestSalesImportAndDeduplication(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// Step 1: Register Tenant (Zen Retailer)
	regPayload := dto.RegisterTenantRequest{
		OrganizationName: "Zen Retailer",
		Slug:             "zen-retailer",
		OwnerName:        "Suresh Prabhu",
		Email:            "suresh@zenretail.com",
		Password:         "Password@123",
	}
	body, _ := json.Marshal(regPayload)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 1000)

	var authResp struct {
		Data dto.AuthTokenResponse `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&authResp)
	token := authResp.Data.AccessToken

	// Step 2: Create a Product Variant with barcode "890111222333" and 50 opening stock
	prodPayload := dto.CreateProductRequest{
		Name:        "Classic Leather Loafer",
		ProductCode: "CLL-01",
		Gender:      "Men",
		Variants: []dto.CreateVariantRequest{
			{
				Size:              "UK 9",
				Colour:            "Brown",
				SKU:               "CLL-BRN-09",
				Barcode:           "890111222333",
				PurchasePrice:     decimal.NewFromFloat(1000.00),
				SellingPrice:      decimal.NewFromFloat(1999.00),
				GSTRate:           decimal.NewFromFloat(12.00),
				OpeningStock:      50,
				MinStockThreshold: 5,
			},
		},
	}
	prodBody, _ := json.Marshal(prodPayload)
	createReq := httptest.NewRequest("POST", "/api/v1/products", bytes.NewReader(prodBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	prodResp, _ := app.Test(createReq, 1000)
	if prodResp.StatusCode != http.StatusCreated {
		t.Fatalf("Product setup failed: %d", prodResp.StatusCode)
	}

	// Step 3: Prepare POS CSV Data (2 bills: BILL-001 completed with 2 units, BILL-002 cancelled)
	sampleCSV := `Bill Number,Date,Barcode,Item Name,Quantity,Rate,Discount,Tax Rate,Tax Amount,Total,Payment Mode,Status
BILL-001,2026-11-01,890111222333,Classic Leather Loafer,2,1999,0,12,479.76,4477.76,UPI,Completed
BILL-002,2026-11-01,890111222333,Classic Leather Loafer,1,1999,0,12,239.88,2238.88,CASH,Cancelled
`

	importPayload := dto.ExecuteImportRequest{
		SourceSystem:    "ZenBiller",
		FileName:        "zenbiller_sales_nov01.csv",
		CSVData:         sampleCSV,
		DeductInventory: true,
	}
	importBody, _ := json.Marshal(importPayload)
	importReq := httptest.NewRequest("POST", "/api/v1/sales/import", bytes.NewReader(importBody))
	importReq.Header.Set("Content-Type", "application/json")
	importReq.Header.Set("Authorization", "Bearer "+token)
	importResp, err := app.Test(importReq, 1000)
	if err != nil || importResp.StatusCode != http.StatusOK {
		t.Fatalf("Sales import failed: %v, code: %d", err, importResp.StatusCode)
	}

	var jobResp struct {
		Data dto.ImportJobResponse `json:"data"`
	}
	json.NewDecoder(importResp.Body).Decode(&jobResp)
	job := jobResp.Data

	if job.ImportedRecords != 2 {
		t.Errorf("Expected 2 imported bills, got %d", job.ImportedRecords)
	}
	if job.DuplicateRecords != 0 {
		t.Errorf("Expected 0 duplicates on first run, got %d", job.DuplicateRecords)
	}

	// Step 4: Verify Stock Deduction
	// Initial stock was 50. BILL-001 sold 2 units. BILL-002 was Cancelled (should NOT deduct stock).
	// Expected stock: 50 - 2 = 48 units.
	barcodeReq := httptest.NewRequest("GET", "/api/v1/products/barcode/890111222333", nil)
	barcodeReq.Header.Set("Authorization", "Bearer "+token)
	barcodeResp, _ := app.Test(barcodeReq, 1000)

	var variantData struct {
		Data dto.ProductVariantResponse `json:"data"`
	}
	json.NewDecoder(barcodeResp.Body).Decode(&variantData)
	if variantData.Data.CurrentStock != 48 {
		t.Errorf("Expected stock after sales import to be 48 (50 - 2 sold), got %d", variantData.Data.CurrentStock)
	}

	// Step 5: Test Idempotent Re-Import (Re-importing identical file must detect duplicates and NEVER deduct stock again)
	reimportReq := httptest.NewRequest("POST", "/api/v1/sales/import", bytes.NewReader(importBody))
	reimportReq.Header.Set("Content-Type", "application/json")
	reimportReq.Header.Set("Authorization", "Bearer "+token)
	reimportResp, _ := app.Test(reimportReq, 1000)
	if reimportResp.StatusCode != http.StatusOK {
		t.Fatalf("Re-import request failed: %d", reimportResp.StatusCode)
	}

	var reimportJobResp struct {
		Data dto.ImportJobResponse `json:"data"`
	}
	json.NewDecoder(reimportResp.Body).Decode(&reimportJobResp)
	reimportJob := reimportJobResp.Data

	if reimportJob.DuplicateRecords != 2 {
		t.Errorf("Expected 2 duplicate records detected on re-import, got %d", reimportJob.DuplicateRecords)
	}
	if reimportJob.ImportedRecords != 0 {
		t.Errorf("Expected 0 new records on re-import, got %d", reimportJob.ImportedRecords)
	}

	// Verify stock is still 48 (not deducted a second time)
	barcodeReq2 := httptest.NewRequest("GET", "/api/v1/products/barcode/890111222333", nil)
	barcodeReq2.Header.Set("Authorization", "Bearer "+token)
	barcodeResp2, _ := app.Test(barcodeReq2, 1000)
	json.NewDecoder(barcodeResp2.Body).Decode(&variantData)
	if variantData.Data.CurrentStock != 48 {
		t.Errorf("Stock was incorrectly deducted on duplicate import! Expected 48, got %d", variantData.Data.CurrentStock)
	}

	// Step 6: Test Sales Listing and Query
	listReq := httptest.NewRequest("GET", "/api/v1/sales?status=COMPLETED", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp, _ := app.Test(listReq, 1000)
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("Sales listing failed: %d", listResp.StatusCode)
	}

	var listData struct {
		Data dto.SaleListResponse `json:"data"`
	}
	json.NewDecoder(listResp.Body).Decode(&listData)
	if listData.Data.TotalCount != 1 {
		t.Errorf("Expected 1 completed sale, got %d", listData.Data.TotalCount)
	}
	if listData.Data.Items[0].BillNumber != "BILL-001" {
		t.Errorf("Expected bill BILL-001, got %s", listData.Data.Items[0].BillNumber)
	}
}
