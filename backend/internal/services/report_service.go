package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) parseDates(fromStr, toStr string) (start, end time.Time) {
	now := time.Now().UTC()
	start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)

	if t, err := time.Parse("2006-01-02", fromStr); err == nil {
		start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}
	if t, err := time.Parse("2006-01-02", toStr); err == nil {
		end = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
	}
	return
}

func (s *ReportService) GetSalesSummary(tenantID uuid.UUID, fromStr, toStr string) (*dto.ReportResponse[dto.SalesSummaryRow], error) {
	start, end := s.parseDates(fromStr, toStr)
	rows, err := s.repo.GetSalesSummaryReport(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	totalVal := decimal.Zero
	for _, r := range rows {
		totalVal = totalVal.Add(r.GrandTotal)
	}

	return &dto.ReportResponse[dto.SalesSummaryRow]{
		ReportName:   "Sales Ledger Summary Report",
		GeneratedAt:  time.Now().UTC(),
		DateFrom:     &start,
		DateTo:       &end,
		TotalRecords: len(rows),
		SummaryValue: totalVal.Round(2),
		Data:         rows,
	}, nil
}

func (s *ReportService) GetStockValuation(tenantID uuid.UUID) (*dto.ReportResponse[dto.StockValuationRow], error) {
	rows, err := s.repo.GetStockValuationReport(tenantID)
	if err != nil {
		return nil, err
	}

	totalVal := decimal.Zero
	for _, r := range rows {
		totalVal = totalVal.Add(r.CostValuation)
	}

	return &dto.ReportResponse[dto.StockValuationRow]{
		ReportName:   "Footwear Stock Valuation Report",
		GeneratedAt:  time.Now().UTC(),
		TotalRecords: len(rows),
		SummaryValue: totalVal.Round(2),
		Data:         rows,
	}, nil
}

func (s *ReportService) GetStockMovements(tenantID uuid.UUID, fromStr, toStr string) (*dto.ReportResponse[dto.StockMovementReportRow], error) {
	start, end := s.parseDates(fromStr, toStr)
	rows, err := s.repo.GetStockMovementReport(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	return &dto.ReportResponse[dto.StockMovementReportRow]{
		ReportName:   "Stock Movement Audit Trail",
		GeneratedAt:  time.Now().UTC(),
		DateFrom:     &start,
		DateTo:       &end,
		TotalRecords: len(rows),
		Data:         rows,
	}, nil
}

func (s *ReportService) ExportSalesCSV(tenantID uuid.UUID, fromStr, toStr string) ([]byte, string, error) {
	report, err := s.GetSalesSummary(tenantID, fromStr, toStr)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Write CSV Header
	w.Write([]string{
		"Bill Number",
		"Sale Date",
		"Payment Method",
		"Item Count",
		"Subtotal (₹)",
		"Discount (₹)",
		"Tax / GST (₹)",
		"Grand Total (₹)",
		"Status",
	})

	for _, r := range report.Data {
		w.Write([]string{
			r.BillNumber,
			r.SaleDate.Format("2006-01-02 15:04:05"),
			r.PaymentMethod,
			fmt.Sprintf("%d", r.ItemCount),
			r.Subtotal.StringFixed(2),
			r.Discount.StringFixed(2),
			r.TaxAmount.StringFixed(2),
			r.GrandTotal.StringFixed(2),
			r.Status,
		})
	}
	w.Flush()

	fileName := fmt.Sprintf("sales_report_%s.csv", time.Now().Format("20060102_150405"))
	return buf.Bytes(), fileName, nil
}

func (s *ReportService) ExportStockValuationCSV(tenantID uuid.UUID) ([]byte, string, error) {
	report, err := s.GetStockValuation(tenantID)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.Write([]string{
		"Product Name",
		"Product Code",
		"Size",
		"Colour",
		"SKU",
		"Barcode",
		"Current Stock",
		"Purchase Cost (₹)",
		"Selling Price (₹)",
		"Valuation at Cost (₹)",
		"Valuation at Retail (₹)",
		"Stock Status",
	})

	for _, r := range report.Data {
		w.Write([]string{
			r.ProductName,
			r.ProductCode,
			r.Size,
			r.Colour,
			r.SKU,
			r.Barcode,
			fmt.Sprintf("%d", r.CurrentStock),
			r.PurchasePrice.StringFixed(2),
			r.SellingPrice.StringFixed(2),
			r.CostValuation.StringFixed(2),
			r.RetailValuation.StringFixed(2),
			r.StockStatus,
		})
	}
	w.Flush()

	fileName := fmt.Sprintf("stock_valuation_%s.csv", time.Now().Format("20060102_150405"))
	return buf.Bytes(), fileName, nil
}

func (s *ReportService) ExportStockMovementCSV(tenantID uuid.UUID, fromStr, toStr string) ([]byte, string, error) {
	report, err := s.GetStockMovements(tenantID, fromStr, toStr)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.Write([]string{
		"Timestamp",
		"Product Name",
		"Barcode",
		"Variant Details",
		"Movement Type",
		"Quantity",
		"Previous Stock",
		"New Stock",
		"Unit Cost (₹)",
		"Reference / Bill #",
		"Notes",
		"Created By",
	})

	for _, r := range report.Data {
		w.Write([]string{
			r.Date.Format("2006-01-02 15:04:05"),
			r.ProductName,
			r.Barcode,
			r.VariantDetails,
			r.MovementType,
			fmt.Sprintf("%d", r.Quantity),
			fmt.Sprintf("%d", r.PreviousStock),
			fmt.Sprintf("%d", r.NewStock),
			r.UnitCost.StringFixed(2),
			r.ReferenceNumber,
			r.Notes,
			r.CreatedBy,
		})
	}
	w.Flush()

	fileName := fmt.Sprintf("stock_movements_%s.csv", time.Now().Format("20060102_150405"))
	return buf.Bytes(), fileName, nil
}
