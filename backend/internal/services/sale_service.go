package services

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type SaleService struct {
	saleRepo    *repositories.SaleRepository
	productRepo *repositories.ProductRepository
}

func NewSaleService(
	saleRepo *repositories.SaleRepository,
	productRepo *repositories.ProductRepository,
) *SaleService {
	return &SaleService{
		saleRepo:    saleRepo,
		productRepo: productRepo,
	}
}

func (s *SaleService) ListSales(
	tenantID uuid.UUID,
	search string,
	status string,
	dateFrom *time.Time,
	dateTo *time.Time,
	page int,
	pageSize int,
) (*dto.SaleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sales, total, err := s.saleRepo.FindSales(tenantID, search, status, dateFrom, dateTo, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]dto.SaleResponse, 0, len(sales))
	for _, sale := range sales {
		items = append(items, s.formatSaleResponse(&sale))
	}

	return &dto.SaleListResponse{
		Items:      items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *SaleService) GetSale(tenantID, saleID uuid.UUID) (*dto.SaleResponse, error) {
	sale, err := s.saleRepo.GetSaleByID(tenantID, saleID)
	if err != nil {
		return nil, errors.New("sale record not found")
	}
	res := s.formatSaleResponse(sale)
	return &res, nil
}

// PreviewImport parses the top rows of a CSV or XLSX for validation
func (s *SaleService) PreviewImport(tenantID uuid.UUID, content []byte, filename string) (*dto.ImportPreviewResponse, error) {
	rows, err := s.extractRows(content, filename)
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, errors.New("file is empty or missing headers")
	}

	headers := rows[0]
	headerMap := s.buildHeaderMap(headers)

	sampleRows := make([]dto.ImportPreviewItem, 0, 5)
	limit := 6
	if len(rows) < limit {
		limit = len(rows)
	}

	for i := 1; i < limit; i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		billNo := s.getCol(row, headerMap, "bill_number", "invoice", "bill no", "billno")
		dateStr := s.getCol(row, headerMap, "date", "bill date", "sale date")
		barcode := s.getCol(row, headerMap, "barcode", "item code", "code", "sku")
		prodName := s.getCol(row, headerMap, "product", "item name", "name", "description")
		qtyStr := s.getCol(row, headerMap, "quantity", "qty")
		priceStr := s.getCol(row, headerMap, "rate", "price", "unit price", "mrp")
		totalStr := s.getCol(row, headerMap, "total", "net amount", "amount", "line total")

		qty, _ := strconv.Atoi(qtyStr)
		if qty <= 0 {
			qty = 1
		}

		matched := false
		if barcode != "" {
			_, err := s.productRepo.FindVariantByBarcode(tenantID, barcode)
			matched = (err == nil)
		}

		sampleRows = append(sampleRows, dto.ImportPreviewItem{
			BillNumber:  billNo,
			Date:        dateStr,
			Barcode:     barcode,
			ProductName: prodName,
			Quantity:    qty,
			UnitPrice:   priceStr,
			Total:       totalStr,
			Matched:     matched,
		})
	}

	return &dto.ImportPreviewResponse{
		TotalRows:       len(rows) - 1,
		DetectedColumns: headers,
		SampleRows:      sampleRows,
	}, nil
}

// ExecuteImport processes all rows with transactional deduplication & stock sync
func (s *SaleService) ExecuteImport(
	tenantID uuid.UUID,
	userID uuid.UUID,
	content []byte,
	fileName string,
	sourceSystem string,
	deductInventory bool,
) (*dto.ImportJobResponse, error) {
	if sourceSystem == "" {
		sourceSystem = "ZenBiller"
	}

	// 1. Create PENDING ImportJob
	job := models.ImportJob{
		ID:           uuid.New(),
		TenantID:     tenantID,
		SourceSystem: sourceSystem,
		FileName:     fileName,
		Status:       models.ImportStatusProcessing,
		CreatedBy:    userID,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := s.saleRepo.CreateImportJob(&job); err != nil {
		return nil, fmt.Errorf("failed to initialize import job: %w", err)
	}

	rows, err := s.extractRows(content, fileName)
	if err != nil {
		job.Status = models.ImportStatusFailed
		job.ErrorDetails = err.Error()
		_ = s.saleRepo.UpdateImportJob(&job)
		return nil, err
	}

	if len(rows) < 2 {
		job.Status = models.ImportStatusFailed
		job.ErrorDetails = "File contains no data rows"
		_ = s.saleRepo.UpdateImportJob(&job)
		return nil, errors.New("file contains no data rows")
	}

	headers := rows[0]
	headerMap := s.buildHeaderMap(headers)

	// Group rows by BillNumber/SourceBillID
	type ParsedLine struct {
		SourceLineID string
		Barcode      string
		ProductName  string
		Quantity     int
		UnitPrice    decimal.Decimal
		Discount     decimal.Decimal
		TaxRate      decimal.Decimal
		TaxAmount    decimal.Decimal
		LineTotal    decimal.Decimal
	}

	type ParsedBill struct {
		BillNumber    string
		SourceBillID  string
		SaleDate      time.Time
		PaymentMethod string
		Status        models.SaleStatus
		Lines         []ParsedLine
	}

	billsMap := make(map[string]*ParsedBill)
	var errorsList []string

	for rowIndex, row := range rows[1:] {
		if len(row) == 0 {
			continue
		}

		billNo := s.getCol(row, headerMap, "bill_number", "invoice", "bill no", "billno")
		if billNo == "" {
			billNo = fmt.Sprintf("BILL-ROW-%d", rowIndex+2)
		}

		sourceBillID := s.getCol(row, headerMap, "source_bill_id", "source id", "bill_id", "bill id")
		if sourceBillID == "" {
			sourceBillID = billNo
		}

		dateStr := s.getCol(row, headerMap, "date", "bill date", "sale date")
		saleDate, parseErr := s.parseDate(dateStr)
		if parseErr != nil {
			saleDate = time.Now().UTC()
		}

		payMode := s.getCol(row, headerMap, "payment mode", "payment method", "mode")
		if payMode == "" {
			payMode = "CASH"
		}

		statusStr := strings.ToUpper(s.getCol(row, headerMap, "status", "bill status"))
		status := models.SaleStatusCompleted
		if strings.Contains(statusStr, "CANCEL") {
			status = models.SaleStatusCancelled
		} else if strings.Contains(statusStr, "RETURN") {
			status = models.SaleStatusReturned
		} else if strings.Contains(statusStr, "VOID") {
			status = models.SaleStatusVoided
		}

		barcode := s.getCol(row, headerMap, "barcode", "item code", "code", "sku")
		prodName := s.getCol(row, headerMap, "product", "item name", "name", "description")
		if prodName == "" {
			prodName = "Retail Item"
		}

		qty, _ := strconv.Atoi(s.getCol(row, headerMap, "quantity", "qty"))
		if qty <= 0 {
			qty = 1
		}

		unitPrice, _ := decimal.NewFromString(s.cleanAmount(s.getCol(row, headerMap, "rate", "price", "unit price", "mrp")))
		discount, _ := decimal.NewFromString(s.cleanAmount(s.getCol(row, headerMap, "discount", "disc")))
		taxRate, _ := decimal.NewFromString(s.cleanAmount(s.getCol(row, headerMap, "tax rate", "gst %", "tax %", "tax_rate")))
		taxAmount, _ := decimal.NewFromString(s.cleanAmount(s.getCol(row, headerMap, "tax amount", "tax", "gst amount")))
		lineTotal, _ := decimal.NewFromString(s.cleanAmount(s.getCol(row, headerMap, "total", "net amount", "amount", "line total")))

		if lineTotal.IsZero() {
			lineTotal = unitPrice.Mul(decimal.NewFromInt(int64(qty))).Sub(discount).Add(taxAmount)
		}

		bill, exists := billsMap[sourceBillID]
		if !exists {
			bill = &ParsedBill{
				BillNumber:    billNo,
				SourceBillID:  sourceBillID,
				SaleDate:      saleDate,
				PaymentMethod: payMode,
				Status:        status,
				Lines:         make([]ParsedLine, 0, 4),
			}
			billsMap[sourceBillID] = bill
		}

		bill.Lines = append(bill.Lines, ParsedLine{
			SourceLineID: fmt.Sprintf("%s-L%d", sourceBillID, len(bill.Lines)+1),
			Barcode:      barcode,
			ProductName:  prodName,
			Quantity:     qty,
			UnitPrice:    unitPrice,
			Discount:     discount,
			TaxRate:      taxRate,
			TaxAmount:    taxAmount,
			LineTotal:    lineTotal,
		})
	}

	job.TotalRecords = len(billsMap)
	importedCount := 0
	duplicateCount := 0
	failedCount := 0

	// 2. Transactionally ingest each bill with deduplication and stock sync
	for _, b := range billsMap {
		txErr := s.saleRepo.DB().Transaction(func(tx *gorm.DB) error {
			// Deduplication check
			var existing models.Sale
			err := tx.Where("tenant_id = ? AND source_system = ? AND source_bill_id = ?",
				tenantID, sourceSystem, b.SourceBillID).First(&existing).Error

			if err == nil {
				// Duplicate detected!
				duplicateCount++
				return nil
			} else if err != gorm.ErrRecordNotFound {
				return err
			}

			// Calculate Bill Totals
			var subtotal, totalDiscount, totalTax, grandTotal decimal.Decimal
			for _, line := range b.Lines {
				subtotal = subtotal.Add(line.UnitPrice.Mul(decimal.NewFromInt(int64(line.Quantity))))
				totalDiscount = totalDiscount.Add(line.Discount)
				totalTax = totalTax.Add(line.TaxAmount)
				grandTotal = grandTotal.Add(line.LineTotal)
			}

			now := time.Now().UTC()
			sale := models.Sale{
				ID:                uuid.New(),
				TenantID:          tenantID,
				SourceSystem:      sourceSystem,
				SourceBillID:      b.SourceBillID,
				BillNumber:        b.BillNumber,
				SaleDate:          b.SaleDate,
				Subtotal:          subtotal,
				Discount:          totalDiscount,
				TaxAmount:         totalTax,
				GrandTotal:        grandTotal,
				PaymentMethod:     b.PaymentMethod,
				Status:            b.Status,
				InventoryDeducted: false,
				ImportJobID:       &job.ID,
				CreatedAt:         now,
				UpdatedAt:         now,
			}

			if err := tx.Create(&sale).Error; err != nil {
				return fmt.Errorf("failed to insert sale: %w", err)
			}

			// Insert Line Items and Match Product Variants
			for _, line := range b.Lines {
				var variantID *uuid.UUID
				if line.Barcode != "" {
					var v models.ProductVariant
					if err := tx.Where("tenant_id = ? AND barcode = ?", tenantID, line.Barcode).First(&v).Error; err == nil {
						variantID = &v.ID
					}
				}

				item := models.SaleItem{
					ID:           uuid.New(),
					TenantID:     tenantID,
					SaleID:       sale.ID,
					SourceLineID: line.SourceLineID,
					VariantID:    variantID,
					ProductName:  line.ProductName,
					Barcode:      line.Barcode,
					Quantity:     line.Quantity,
					UnitPrice:    line.UnitPrice,
					Discount:     line.Discount,
					TaxRate:      line.TaxRate,
					TaxAmount:    line.TaxAmount,
					LineTotal:    line.LineTotal,
					CreatedAt:    now,
				}
				if err := tx.Create(&item).Error; err != nil {
					return fmt.Errorf("failed to insert sale item: %w", err)
				}

				// Deduct Inventory if configured, sale is COMPLETED, and variant is matched
				if deductInventory && sale.Status == models.SaleStatusCompleted && variantID != nil {
					// Directly update stock in tx to preserve transaction boundary
					var variant models.ProductVariant
					if err := tx.Where("id = ? AND tenant_id = ?", *variantID, tenantID).First(&variant).Error; err == nil {
						prev := variant.CurrentStock
						newStock := prev - line.Quantity
						_ = tx.Model(&variant).Update("current_stock", newStock)

						// Record movement in stock ledger
						mov := models.StockMovement{
							ID:              uuid.New(),
							TenantID:        tenantID,
							VariantID:       *variantID,
							MovementType:    models.MovementSale,
							Quantity:        line.Quantity,
							PreviousStock:   prev,
							NewStock:        newStock,
							UnitCost:        variant.PurchasePrice,
							ReferenceNumber: sale.BillNumber,
							Notes:           fmt.Sprintf("Auto-deducted from %s import", sourceSystem),
							CreatedBy:       userID,
							CreatedAt:       now,
						}
						_ = tx.Create(&mov)
						sale.InventoryDeducted = true
					}
				}
			}

			if sale.InventoryDeducted {
				_ = tx.Model(&sale).Update("inventory_deducted", true)
			}

			importedCount++
			return nil
		})

		if txErr != nil {
			failedCount++
			errorsList = append(errorsList, fmt.Sprintf("Bill %s: %v", b.BillNumber, txErr))
		}
	}

	job.ImportedRecords = importedCount
	job.DuplicateRecords = duplicateCount
	job.FailedRecords = failedCount
	job.Status = models.ImportStatusCompleted
	if failedCount > 0 && importedCount == 0 {
		job.Status = models.ImportStatusFailed
	}

	if len(errorsList) > 0 {
		errBytes, _ := json.Marshal(errorsList)
		job.ErrorDetails = string(errBytes)
	}

	job.UpdatedAt = time.Now().UTC()
	_ = s.saleRepo.UpdateImportJob(&job)

	return &dto.ImportJobResponse{
		ID:               job.ID,
		SourceSystem:     job.SourceSystem,
		FileName:         job.FileName,
		TotalRecords:     job.TotalRecords,
		ImportedRecords:  job.ImportedRecords,
		FailedRecords:    job.FailedRecords,
		DuplicateRecords: job.DuplicateRecords,
		Status:           string(job.Status),
		ErrorDetails:     job.ErrorDetails,
		CreatedAt:        job.CreatedAt,
	}, nil
}

func (s *SaleService) GetImportHistory(tenantID uuid.UUID) ([]dto.ImportJobResponse, error) {
	jobs, err := s.saleRepo.GetImportJobs(tenantID, 50)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ImportJobResponse, 0, len(jobs))
	for _, j := range jobs {
		res = append(res, dto.ImportJobResponse{
			ID:               j.ID,
			SourceSystem:     j.SourceSystem,
			FileName:         j.FileName,
			TotalRecords:     j.TotalRecords,
			ImportedRecords:  j.ImportedRecords,
			FailedRecords:    j.FailedRecords,
			DuplicateRecords: j.DuplicateRecords,
			Status:           string(j.Status),
			ErrorDetails:     j.ErrorDetails,
			CreatedAt:        j.CreatedAt,
		})
	}
	return res, nil
}

// Helpers
func (s *SaleService) extractRows(content []byte, filename string) ([][]string, error) {
	lower := strings.ToLower(filename)
	if strings.HasSuffix(lower, ".xlsx") {
		f, err := excelize.OpenReader(bytes.NewReader(content))
		if err != nil {
			return nil, fmt.Errorf("failed to read excel file: %w", err)
		}
		defer f.Close()

		sheetList := f.GetSheetList()
		if len(sheetList) == 0 {
			return nil, errors.New("excel workbook contains no sheets")
		}
		return f.GetRows(sheetList[0])
	}

	// CSV Parser
	reader := csv.NewReader(bytes.NewReader(content))
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	var rows [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Skip corrupted rows
			continue
		}
		rows = append(rows, record)
	}
	return rows, nil
}

func (s *SaleService) buildHeaderMap(headers []string) map[string]int {
	m := make(map[string]int)
	for i, h := range headers {
		clean := strings.ToLower(strings.TrimSpace(h))
		clean = strings.ReplaceAll(clean, "_", " ")
		m[clean] = i
	}
	return m
}

func (s *SaleService) getCol(row []string, headerMap map[string]int, aliases ...string) string {
	for _, alias := range aliases {
		cleanAlias := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(alias, "_", " ")))
		if idx, ok := headerMap[cleanAlias]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
	}
	return ""
}

func (s *SaleService) parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"02-01-2006",
		"02/01/2006",
		"2006/01/02",
		"2006-01-02 15:04:05",
		"02-01-2006 15:04:05",
		"02/01/2006 15:04:05",
		time.RFC3339,
	}
	for _, f := range formats {
		if t, err := time.Parse(f, dateStr); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Now().UTC(), errors.New("unsupported date format")
}

func (s *SaleService) cleanAmount(val string) string {
	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, "₹", "")
	val = strings.ReplaceAll(val, "$", "")
	val = strings.ReplaceAll(val, ",", "")
	if val == "" {
		return "0"
	}
	return val
}

func (s *SaleService) formatSaleResponse(sale *models.Sale) dto.SaleResponse {
	items := make([]dto.SaleItemResponse, 0, len(sale.Items))
	for _, it := range sale.Items {
		items = append(items, dto.SaleItemResponse{
			ID:          it.ID,
			ProductName: it.ProductName,
			SKU:         it.SKU,
			Barcode:     it.Barcode,
			VariantID:   it.VariantID,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
			Discount:    it.Discount,
			TaxRate:     it.TaxRate,
			TaxAmount:   it.TaxAmount,
			LineTotal:   it.LineTotal,
		})
	}

	return dto.SaleResponse{
		ID:                sale.ID,
		SourceSystem:      sale.SourceSystem,
		SourceBillID:      sale.SourceBillID,
		BillNumber:        sale.BillNumber,
		SaleDate:          sale.SaleDate,
		Subtotal:          sale.Subtotal,
		Discount:          sale.Discount,
		TaxAmount:         sale.TaxAmount,
		GrandTotal:        sale.GrandTotal,
		PaymentMethod:     sale.PaymentMethod,
		Status:            string(sale.Status),
		InventoryDeducted: sale.InventoryDeducted,
		ItemCount:         len(sale.Items),
		Items:             items,
		CreatedAt:         sale.CreatedAt,
	}
}
