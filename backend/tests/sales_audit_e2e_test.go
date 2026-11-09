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
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

func TestSalesAnalyticsControlledAuditE2E(t *testing.T) {
	app, _, dbManager := setupTestApp(t)
	db := dbManager.DB

	// 1. Onboard Tenant
	regPayload := dto.RegisterTenantRequest{
		OrganizationName: "Apex Retail Footwear",
		Slug:             "apex-retail-footwear",
		OwnerName:        "Suresh Varma",
		Email:            "suresh@apexfootwear.com",
		Password:         "SecurePassword@123",
	}
	body, _ := json.Marshal(regPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, 10000)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Tenant registration failed: %v", err)
	}

	var regResp struct {
		Data dto.AuthTokenResponse `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&regResp)
	token := regResp.Data.AccessToken
	tenantID := regResp.Data.ActiveTenant.TenantID

	// 2. Controlled Catalogue Dataset
	// Prod 1: Bata Pro Runner (2 variants)
	prod1 := models.Product{
		TenantID:    tenantID,
		Name:        "Bata Pro Runner",
		ProductCode: "BPR-001",
		Status:      "active",
		Variants: []models.ProductVariant{
			{
				TenantID:      tenantID,
				SKU:           "BPR-BLK-08",
				Barcode:       "890001",
				Size:          "8",
				Colour:        "Black",
				PurchasePrice: decimal.NewFromInt(800),
				SellingPrice:  decimal.NewFromInt(1500),
				CurrentStock:  60,
			},
			{
				TenantID:      tenantID,
				SKU:           "BPR-BLK-09",
				Barcode:       "890002",
				Size:          "9",
				Colour:        "Black",
				PurchasePrice: decimal.NewFromInt(800),
				SellingPrice:  decimal.NewFromInt(1500),
				CurrentStock:  50,
			},
		},
	}
	db.Create(&prod1)

	// Prod 2: Woodland Trek Boot
	prod2 := models.Product{
		TenantID:    tenantID,
		Name:        "Woodland Trek Boot",
		ProductCode: "WTB-001",
		Status:      "active",
		Variants: []models.ProductVariant{
			{
				TenantID:      tenantID,
				SKU:           "WTB-TAN-09",
				Barcode:       "890003",
				Size:          "9",
				Colour:        "Tan",
				PurchasePrice: decimal.NewFromInt(1200),
				SellingPrice:  decimal.NewFromInt(2500),
				CurrentStock:  30,
			},
		},
	}
	db.Create(&prod2)

	// Prod 3: Inactive Loafer (Zero Sales Alert)
	prod3 := models.Product{
		TenantID:    tenantID,
		Name:        "Relaxo Casual Slipper",
		ProductCode: "RCS-001",
		Status:      "active",
		Variants: []models.ProductVariant{
			{
				TenantID:      tenantID,
				SKU:           "RCS-BLU-07",
				Barcode:       "890004",
				Size:          "7",
				Colour:        "Blue",
				PurchasePrice: decimal.NewFromInt(150),
				SellingPrice:  decimal.NewFromInt(350),
				CurrentStock:  40,
			},
		},
	}
	db.Create(&prod3)

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	lastMonth := today.AddDate(0, -1, -5)

	// Controlled Transaction 1: Today, Completed
	// 2 units Bata Pro (Size 8), Gross: 3192, Discount: 150, Tax: 342, Net: 2542, COGS: 1600
	sale1 := models.Sale{
		TenantID:      tenantID,
		SourceSystem:  "ZenBiller",
		SourceBillID:  "ZB-AUDIT-001",
		BillNumber:    "AUDIT-001",
		SaleDate:      today,
		Subtotal:      decimal.NewFromInt(3000),
		Discount:      decimal.NewFromInt(150),
		TaxAmount:     decimal.NewFromInt(342),
		GrandTotal:    decimal.NewFromInt(3192),
		PaymentMethod: "UPI",
		Status:        models.SaleStatusCompleted,
		Items: []models.SaleItem{
			{
				TenantID:    tenantID,
				VariantID:   &prod1.Variants[0].ID,
				ProductName: "Bata Pro Runner - Size 8",
				Barcode:     "890001",
				SKU:         "BPR-BLK-08",
				Quantity:    2,
				UnitPrice:   decimal.NewFromInt(1500),
				Discount:    decimal.NewFromInt(150),
				LineTotal:   decimal.NewFromInt(2850),
			},
		},
	}
	db.Create(&sale1)

	// Controlled Transaction 2: Yesterday, Completed
	// 1 unit Woodland Trek (Size 9), Gross: 2950, Discount: 0, Tax: 450, Net: 2500, COGS: 1200
	sale2 := models.Sale{
		TenantID:      tenantID,
		SourceSystem:  "ZenBiller",
		SourceBillID:  "ZB-AUDIT-002",
		BillNumber:    "AUDIT-002",
		SaleDate:      yesterday,
		Subtotal:      decimal.NewFromInt(2500),
		Discount:      decimal.Zero,
		TaxAmount:     decimal.NewFromInt(450),
		GrandTotal:    decimal.NewFromInt(2950),
		PaymentMethod: "CARD",
		Status:        models.SaleStatusCompleted,
		Items: []models.SaleItem{
			{
				TenantID:    tenantID,
				VariantID:   &prod2.Variants[0].ID,
				ProductName: "Woodland Trek Boot - Size 9",
				Barcode:     "890003",
				SKU:         "WTB-TAN-09",
				Quantity:    1,
				UnitPrice:   decimal.NewFromInt(2500),
				Discount:    decimal.Zero,
				LineTotal:   decimal.NewFromInt(2500),
			},
		},
	}
	db.Create(&sale2)

	// Controlled Transaction 3: Last Month, Completed
	// 3 units Bata Pro (Size 9), Gross: 4704, Discount: 300, Tax: 504, Net: 4200, COGS: 2400
	sale3 := models.Sale{
		TenantID:      tenantID,
		SourceSystem:  "ZenBiller",
		SourceBillID:  "ZB-AUDIT-003",
		BillNumber:    "AUDIT-003",
		SaleDate:      lastMonth,
		Subtotal:      decimal.NewFromInt(4500),
		Discount:      decimal.NewFromInt(300),
		TaxAmount:     decimal.NewFromInt(504),
		GrandTotal:    decimal.NewFromInt(4704),
		PaymentMethod: "CASH",
		Status:        models.SaleStatusCompleted,
		Items: []models.SaleItem{
			{
				TenantID:    tenantID,
				VariantID:   &prod1.Variants[1].ID,
				ProductName: "Bata Pro Runner - Size 9",
				Barcode:     "890002",
				SKU:         "BPR-BLK-09",
				Quantity:    3,
				UnitPrice:   decimal.NewFromInt(1500),
				Discount:    decimal.NewFromInt(300),
				LineTotal:   decimal.NewFromInt(4200),
			},
		},
	}
	db.Create(&sale3)

	// Controlled Transaction 4: Today, Return worth 500
	sale4 := models.Sale{
		TenantID:      tenantID,
		SourceSystem:  "ZenBiller",
		SourceBillID:  "ZB-AUDIT-RET",
		BillNumber:    "AUDIT-RET-001",
		SaleDate:      today,
		Subtotal:      decimal.NewFromInt(500),
		Discount:      decimal.Zero,
		TaxAmount:     decimal.Zero,
		GrandTotal:    decimal.NewFromInt(500),
		PaymentMethod: "CASH",
		Status:        models.SaleStatusReturned,
	}
	db.Create(&sale4)

	// Controlled Transaction 5: Today, Cancelled Bill (MUST BE EXCLUDED)
	sale5 := models.Sale{
		TenantID:      tenantID,
		SourceSystem:  "ZenBiller",
		SourceBillID:  "ZB-AUDIT-CANCEL",
		BillNumber:    "AUDIT-CAN-001",
		SaleDate:      today,
		Subtotal:      decimal.NewFromInt(1500),
		Discount:      decimal.Zero,
		TaxAmount:     decimal.NewFromInt(180),
		GrandTotal:    decimal.NewFromInt(1680),
		PaymentMethod: "CASH",
		Status:        models.SaleStatusCancelled,
	}
	db.Create(&sale5)

	t.Run("Audit_Today_Financial_Accuracy", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/performance?range=today", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req, 10000)

		var apiResp struct {
			Data dto.SalesPerformanceSummary `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		data := apiResp.Data

		// Expected for Today:
		// Completed Bills: 1 (Sale 1)
		if data.TotalBills != 1 {
			t.Errorf("Expected 1 completed bill today, got %d", data.TotalBills)
		}
		// Units Sold: 2
		if data.TotalUnitsSold != 2 {
			t.Errorf("Expected 2 units sold today, got %d", data.TotalUnitsSold)
		}
		// Gross Sales: 3192.00
		if data.GrossSales.StringFixed(2) != "3192.00" {
			t.Errorf("Expected Gross Sales 3192.00, got %s", data.GrossSales.StringFixed(2))
		}
		// Total Discounts: 150.00
		if data.TotalDiscounts.StringFixed(2) != "150.00" {
			t.Errorf("Expected Discounts 150.00, got %s", data.TotalDiscounts.StringFixed(2))
		}
		// Returns: 500.00
		if data.TotalReturns.StringFixed(2) != "500.00" {
			t.Errorf("Expected Returns 500.00, got %s", data.TotalReturns.StringFixed(2))
		}
		// Net Sales: 3192 - 150 - 500 = 2542.00
		if data.NetSales.StringFixed(2) != "2542.00" {
			t.Errorf("Expected Net Sales 2542.00, got %s", data.NetSales.StringFixed(2))
		}
		// COGS: 2 * 800 = 1600.00
		if data.CostOfGoodsSold.StringFixed(2) != "1600.00" {
			t.Errorf("Expected COGS 1600.00, got %s", data.CostOfGoodsSold.StringFixed(2))
		}
		// Gross Profit: 2542 - 1600 = 942.00
		if data.GrossProfit.StringFixed(2) != "942.00" {
			t.Errorf("Expected Gross Profit 942.00, got %s", data.GrossProfit.StringFixed(2))
		}
		// Profit Margin %: 942 / 2542 * 100 = 37.06%
		if data.ProfitMarginPct.StringFixed(2) != "37.06" {
			t.Errorf("Expected Margin 37.06%%, got %s", data.ProfitMarginPct.StringFixed(2))
		}
	})

	t.Run("Audit_Reconciliation_Across_Daily_Monthly_Yearly_Views", func(t *testing.T) {
		// Query All Time Performance (Yearly / Custom window encompassing all 3 sales)
		fromStr := lastMonth.AddDate(0, 0, -1).Format("2006-01-02")
		toStr := today.AddDate(0, 0, 1).Format("2006-01-02")

		// 1. Overall Performance Summary
		reqPerf := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/performance?range=custom&from="+fromStr+"&to="+toStr, nil)
		reqPerf.Header.Set("Authorization", "Bearer "+token)
		respPerf, _ := app.Test(reqPerf, 10000)
		var pResp struct {
			Data dto.SalesPerformanceSummary `json:"data"`
		}
		json.NewDecoder(respPerf.Body).Decode(&pResp)
		totalNetExpected := "9896.00" // (3192 + 2950 + 4704) - 450 - 500 = 9896.00
		if pResp.Data.NetSales.StringFixed(2) != totalNetExpected {
			t.Fatalf("Expected overall Net Sales %s, got %s", totalNetExpected, pResp.Data.NetSales.StringFixed(2))
		}

		// 2. Daily Breakdown Sum Reconciliation
		reqDaily := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/daily?range=custom&from="+fromStr+"&to="+toStr, nil)
		reqDaily.Header.Set("Authorization", "Bearer "+token)
		respDaily, _ := app.Test(reqDaily, 10000)
		var dResp struct {
			Data []dto.DailyBreakdownItem `json:"data"`
		}
		json.NewDecoder(respDaily.Body).Decode(&dResp)
		dailySumNet := decimal.Zero
		var dailySumUnits, dailySumBills int64
		for _, d := range dResp.Data {
			dailySumNet = dailySumNet.Add(d.NetSales)
			dailySumUnits += d.UnitsSold
			dailySumBills += d.Bills
		}
		if dailySumNet.StringFixed(2) != totalNetExpected {
			t.Errorf("RECONCILIATION FAILURE: Daily Net Sales Sum %s does not equal expected %s", dailySumNet.StringFixed(2), totalNetExpected)
		}
		if dailySumBills != 3 {
			t.Errorf("RECONCILIATION FAILURE: Daily Bills Sum %d != 3", dailySumBills)
		}
		if dailySumUnits != 6 {
			t.Errorf("RECONCILIATION FAILURE: Daily Units Sum %d != 6", dailySumUnits)
		}

		// 3. Monthly Breakdown Sum Reconciliation
		reqMo := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/monthly?range=custom&from="+fromStr+"&to="+toStr, nil)
		reqMo.Header.Set("Authorization", "Bearer "+token)
		respMo, _ := app.Test(reqMo, 10000)
		var mResp struct {
			Data []dto.MonthlyBreakdownItem `json:"data"`
		}
		json.NewDecoder(respMo.Body).Decode(&mResp)
		monthlySumNet := decimal.Zero
		var monthlySumUnits, monthlySumBills int64
		for _, m := range mResp.Data {
			monthlySumNet = monthlySumNet.Add(m.NetSales)
			monthlySumUnits += m.UnitsSold
			monthlySumBills += m.Bills
		}
		if monthlySumNet.StringFixed(2) != totalNetExpected {
			t.Errorf("RECONCILIATION FAILURE: Monthly Net Sales Sum %s does not equal expected %s", monthlySumNet.StringFixed(2), totalNetExpected)
		}
		if monthlySumBills != 3 {
			t.Errorf("RECONCILIATION FAILURE: Monthly Bills Sum %d != 3", monthlySumBills)
		}
		if monthlySumUnits != 6 {
			t.Errorf("RECONCILIATION FAILURE: Monthly Units Sum %d != 6", monthlySumUnits)
		}

		// 4. Yearly Breakdown Sum Reconciliation
		reqYr := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/yearly?range=custom&from="+fromStr+"&to="+toStr, nil)
		reqYr.Header.Set("Authorization", "Bearer "+token)
		respYr, _ := app.Test(reqYr, 10000)
		var yResp struct {
			Data []dto.YearlyBreakdownItem `json:"data"`
		}
		json.NewDecoder(respYr.Body).Decode(&yResp)
		yearlySumNet := decimal.Zero
		for _, y := range yResp.Data {
			yearlySumNet = yearlySumNet.Add(y.NetSales)
		}
		if yearlySumNet.StringFixed(2) != totalNetExpected {
			t.Errorf("RECONCILIATION FAILURE: Yearly Net Sales Sum %s does not equal expected %s", yearlySumNet.StringFixed(2), totalNetExpected)
		}

		t.Logf("[Audit PASS] Perfect Reconciliation Verified: Daily (%s) == Monthly (%s) == Yearly (%s) == Summary (%s)",
			dailySumNet.StringFixed(2), monthlySumNet.StringFixed(2), yearlySumNet.StringFixed(2), totalNetExpected)
	})

	t.Run("Audit_Zero_Sales_Alert_Verification", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/products/detailed?range=this_year", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ := app.Test(req, 10000)

		var apiResp struct {
			Data dto.ProductPerformanceResponse `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&apiResp)
		data := apiResp.Data

		// Relaxo Casual Slipper has 40 units in stock and 0 sales
		foundZero := false
		for _, p := range data.ZeroSales {
			if p.ProductCode == "RCS-001" {
				foundZero = true
				if p.UnitsSold != 0 {
					t.Errorf("Expected 0 units sold, got %d", p.UnitsSold)
				}
				if p.StockOnHand != 40 {
					t.Errorf("Expected 40 stock on hand, got %d", p.StockOnHand)
				}
			}
		}
		if !foundZero {
			t.Errorf("Zero Sales product RCS-001 not detected in ZeroSales alerts")
		}
	})

	t.Run("Audit_Excel_Multi_Worksheet_Integrity", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/export/excel?range=this_year", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req, 10000)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Excel export failed: status=%d, err=%v", resp.StatusCode, err)
		}

		excelFile, err := excelize.OpenReader(resp.Body)
		if err != nil {
			t.Fatalf("Failed to open Excel workbook: %v", err)
		}
		defer excelFile.Close()

		sheetList := excelFile.GetSheetList()
		requiredSheets := []string{
			"Summary",
			"Daily Sales",
			"Monthly Sales",
			"Yearly Summary",
			"Product Sales",
			"Category & Brand Sales",
			"Zero Sales Alert",
			"Footwear Matrix",
			"Payment Summary",
			"Transaction Details",
		}

		for _, required := range requiredSheets {
			found := false
			for _, actual := range sheetList {
				if actual == required {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Required Excel worksheet '%s' missing from workbook. Sheets present: %v", required, sheetList)
			}
		}

		// Verify Title in Summary
		val, _ := excelFile.GetCellValue("Summary", "A1")
		if val != "KadaiIQ - Sales Performance & Revenue Analytics" {
			t.Errorf("Unexpected Summary Title in A1: %s", val)
		}

		// Verify Zero Sales Alert sheet contains Product Code column
		zCol, _ := excelFile.GetCellValue("Zero Sales Alert", "A1")
		if zCol != "Product Code" {
			t.Errorf("Expected 'Product Code' in Zero Sales Alert A1, got %s", zCol)
		}
	})
}

