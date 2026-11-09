package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func setupSalesPerformanceTestDB(t *testing.T) (*gorm.DB, *fiber.App, string, uuid.UUID) {
	app, _, dbManager := setupTestApp(t)
	db := dbManager.DB

	// Register user & tenant
	regPayload := dto.RegisterTenantRequest{
		OrganizationName: "Metro Shoes Chennai",
		Slug:             "metro-shoes-chennai",
		OwnerName:        "Metro Admin",
		Email:            "metro.admin@example.com",
		Password:         "Password@123",
	}
	body, _ := json.Marshal(regPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 10000)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to register tenant: status=%d, err=%v", resp.StatusCode, err)
	}

	var regResp struct {
		Data dto.AuthTokenResponse `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&regResp)
	token := regResp.Data.AccessToken
	tenantID := regResp.Data.ActiveTenant.TenantID

	// Create Product 1: Sports Running Shoes (with 2 variants)
	prod1 := models.Product{
		TenantID:    tenantID,
		Name:        "SpeedRunner Pro",
		ProductCode: "PROD-SR01",
		Status:      "active",
		Variants: []models.ProductVariant{
			{
				TenantID:      tenantID,
				SKU:           "SR-BLK-09",
				Barcode:       "8901001001",
				Size:          "9",
				Colour:        "Black",
				PurchasePrice: decimal.NewFromInt(1000),
				SellingPrice:  decimal.NewFromInt(2000),
				CurrentStock:  50,
			},
			{
				TenantID:      tenantID,
				SKU:           "SR-BLU-08",
				Barcode:       "8901001002",
				Size:          "8",
				Colour:        "Blue",
				PurchasePrice: decimal.NewFromInt(1000),
				SellingPrice:  decimal.NewFromInt(2000),
				CurrentStock:  40,
			},
		},
	}
	if err := db.Create(&prod1).Error; err != nil {
		t.Fatalf("Failed to create prod1: %v", err)
	}

	// Create Product 2: Zero Sales Product (in catalogue, has stock, but 0 sales)
	prod2 := models.Product{
		TenantID:    tenantID,
		Name:        "Old Classic Leather Loafer",
		ProductCode: "PROD-OLL01",
		Status:      "active",
		Variants: []models.ProductVariant{
			{
				TenantID:      tenantID,
				SKU:           "OLL-BRN-10",
				Barcode:       "8902002001",
				Size:          "10",
				Colour:        "Brown",
				PurchasePrice: decimal.NewFromInt(1500),
				SellingPrice:  decimal.NewFromInt(3000),
				CurrentStock:  25,
			},
		},
	}
	if err := db.Create(&prod2).Error; err != nil {
		t.Fatalf("Failed to create prod2: %v", err)
	}

	now := time.Now().UTC()

	// Create Completed Sale 1: Today, 2 units of SR-BLK-09 (Gross 4000, Discount 200, Tax 380, GrandTotal 4180)
	sale1 := models.Sale{
		TenantID:      tenantID,
		SourceSystem:  "ZenBiller",
		SourceBillID:  "ZB-PERF-001",
		BillNumber:    "INV-001",
		SaleDate:      now,
		Subtotal:      decimal.NewFromInt(3800),
		Discount:      decimal.NewFromInt(200),
		TaxAmount:     decimal.NewFromInt(380),
		GrandTotal:    decimal.NewFromInt(4180),
		PaymentMethod: "UPI",
		Status:        models.SaleStatusCompleted,
		Items: []models.SaleItem{
			{
				TenantID:    tenantID,
				VariantID:   &prod1.Variants[0].ID,
				ProductName: "SpeedRunner Pro - Black 9",
				Barcode:     "8901001001",
				SKU:         "SR-BLK-09",
				Quantity:    2,
				UnitPrice:   decimal.NewFromInt(2000),
				Discount:    decimal.NewFromInt(200),
				LineTotal:   decimal.NewFromInt(3800),
			},
		},
	}
	db.Create(&sale1)

	// Create Completed Sale 2: Yesterday, 1 unit of SR-BLU-08 (Gross 2000, Discount 0, Tax 200, GrandTotal 2200)
	sale2 := models.Sale{
		TenantID:      tenantID,
		SourceSystem:  "ZenBiller",
		SourceBillID:  "ZB-PERF-002",
		BillNumber:    "INV-002",
		SaleDate:      now.AddDate(0, 0, -1),
		Subtotal:      decimal.NewFromInt(2000),
		Discount:      decimal.Zero,
		TaxAmount:     decimal.NewFromInt(200),
		GrandTotal:    decimal.NewFromInt(2200),
		PaymentMethod: "CASH",
		Status:        models.SaleStatusCompleted,
		Items: []models.SaleItem{
			{
				TenantID:    tenantID,
				VariantID:   &prod1.Variants[1].ID,
				ProductName: "SpeedRunner Pro - Blue 8",
				Barcode:     "8901001002",
				SKU:         "SR-BLU-08",
				Quantity:    1,
				UnitPrice:   decimal.NewFromInt(2000),
				Discount:    decimal.Zero,
				LineTotal:   decimal.NewFromInt(2000),
			},
		},
	}
	db.Create(&sale2)

	// Create Returned Sale 3: Today, Return worth 500
	sale3 := models.Sale{
		TenantID:      tenantID,
		SourceSystem:  "ZenBiller",
		SourceBillID:  "ZB-PERF-RET",
		BillNumber:    "RET-001",
		SaleDate:      now,
		Subtotal:      decimal.NewFromInt(500),
		Discount:      decimal.Zero,
		TaxAmount:     decimal.Zero,
		GrandTotal:    decimal.NewFromInt(500),
		PaymentMethod: "CASH",
		Status:        models.SaleStatusReturned,
	}
	db.Create(&sale3)

	return db, app, token, tenantID
}

func TestSalesPerformanceAndRevenueAnalytics(t *testing.T) {
	_, app, token, _ := setupSalesPerformanceTestDB(t)

	t.Run("Executive_Sales_Performance_Today", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/performance?range=today", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req, 10000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Performance endpoint failed: status=%d, err=%v", resp.StatusCode, err)
		}

		var apiResp struct {
			Data dto.SalesPerformanceSummary `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		p := apiResp.Data
		if p.TotalBills != 1 {
			t.Errorf("Expected 1 completed bill today, got %d", p.TotalBills)
		}
		if p.TotalUnitsSold != 2 {
			t.Errorf("Expected 2 units sold today, got %d", p.TotalUnitsSold)
		}
		if p.TotalDiscounts.StringFixed(2) != "200.00" {
			t.Errorf("Expected discount 200, got %s", p.TotalDiscounts)
		}
		if p.TotalReturns.StringFixed(2) != "500.00" {
			t.Errorf("Expected returns 500, got %s", p.TotalReturns)
		}
		// Gross Sales: 4180
		// Net Sales: 4180 - 200 - 500 = 3480
		expectedNet := decimal.NewFromInt(4180).Sub(decimal.NewFromInt(200)).Sub(decimal.NewFromInt(500)).StringFixed(2)
		if p.NetSales.StringFixed(2) != expectedNet {
			t.Errorf("Expected Net Sales %s, got %s", expectedNet, p.NetSales.StringFixed(2))
		}
		// COGS: 2 * 1000 = 2000
		if p.CostOfGoodsSold.StringFixed(2) != "2000.00" {
			t.Errorf("Expected COGS 2000.00, got %s", p.CostOfGoodsSold)
		}
		// Gross Profit: NetSales - COGS = 3480 - 2000 = 1480
		expectedProfit := decimal.RequireFromString(expectedNet).Sub(decimal.NewFromInt(2000)).StringFixed(2)
		if p.GrossProfit.StringFixed(2) != expectedProfit {
			t.Errorf("Expected Gross Profit %s, got %s", expectedProfit, p.GrossProfit.StringFixed(2))
		}
	})

	t.Run("Detailed_Product_Performance_And_Zero_Sales_Alert", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/products/detailed?range=this_month", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req, 10000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Detailed products failed: status=%d, err=%v", resp.StatusCode, err)
		}

		var apiResp struct {
			Data dto.ProductPerformanceResponse `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)

		data := apiResp.Data
		if len(data.AllProducts) < 3 {
			t.Fatalf("Expected at least 3 product variants, got %d", len(data.AllProducts))
		}

		// Verify zero-sales alert identifies the inactive loafer
		foundZero := false
		for _, z := range data.ZeroSales {
			if z.ProductCode == "PROD-OLL01" {
				foundZero = true
				if z.UnitsSold != 0 {
					t.Errorf("Expected 0 units sold for zero-sales alert, got %d", z.UnitsSold)
				}
				if z.StockOnHand != 25 {
					t.Errorf("Expected 25 idle stock units, got %d", z.StockOnHand)
				}
			}
		}
		if !foundZero {
			t.Errorf("Expected to find PROD-OLL01 in zero-sales alerts")
		}
	})

	t.Run("Multi_Sheet_Excel_Report_Generation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/export/excel?range=this_month", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req, 10000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Excel export failed: status=%d, err=%v", resp.StatusCode, err)
		}

		contentType := resp.Header.Get("Content-Type")
		expectedCT := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		if contentType != expectedCT {
			t.Errorf("Expected Content-Type %s, got %s", expectedCT, contentType)
		}

		// Read response body into Excelize
		excelFile, err := excelize.OpenReader(resp.Body)
		if err != nil {
			t.Fatalf("Failed to open generated Excel spreadsheet: %v", err)
		}
		defer excelFile.Close()

		sheets := excelFile.GetSheetList()
		expectedSheets := []string{
			"Summary",
			"Daily Sales",
			"Monthly Breakdown",
			"Product Sales",
			"Zero Sales Alert",
			"Footwear Matrix",
			"Payment Summary",
			"Transaction Details",
		}

		for _, sheetName := range expectedSheets {
			found := false
			for _, actual := range sheets {
				if actual == sheetName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected sheet '%s' not found in Excel workbook. Actual sheets: %v", sheetName, sheets)
			}
		}

		// Check cell value in Summary sheet
		val, err := excelFile.GetCellValue("Summary", "A1")
		if err != nil || val != "KadaiIQ - Sales Performance & Revenue Analytics" {
			t.Errorf("Unexpected Summary title in A1: %s (err: %v)", val, err)
		}

		// Check zero sales alert sheet
		zVal, _ := excelFile.GetCellValue("Zero Sales Alert", "A1")
		if zVal != "Product Code" {
			t.Errorf("Expected 'Product Code' in Zero Sales Alert A1, got %s", zVal)
		}
	})
}
