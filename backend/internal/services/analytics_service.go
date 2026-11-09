package services

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

type AnalyticsService struct {
	repo *repositories.AnalyticsRepository
}

func NewAnalyticsService(repo *repositories.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// ResolveDateWindow calculates current and matching comparison date boundaries.
func (s *AnalyticsService) ResolveDateWindow(rangeKey, fromStr, toStr string) (start, end, prevStart, prevEnd time.Time, resolvedKey string) {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)

	switch rangeKey {
	case "today":
		start = startOfDay
		end = endOfDay
		prevStart = start.AddDate(0, 0, -1)
		prevEnd = end.AddDate(0, 0, -1)
		resolvedKey = "today"

	case "yesterday":
		start = startOfDay.AddDate(0, 0, -1)
		end = endOfDay.AddDate(0, 0, -1)
		prevStart = start.AddDate(0, 0, -1)
		prevEnd = end.AddDate(0, 0, -1)
		resolvedKey = "yesterday"

	case "daily":
		// Past 30 days daily window
		start = startOfDay.AddDate(0, 0, -29)
		end = endOfDay
		prevStart = start.AddDate(0, 0, -30)
		prevEnd = end.AddDate(0, 0, -30)
		resolvedKey = "daily"

	case "7d":
		start = startOfDay.AddDate(0, 0, -6)
		end = endOfDay
		prevStart = start.AddDate(0, 0, -7)
		prevEnd = end.AddDate(0, 0, -7)
		resolvedKey = "7d"

	case "this_month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = endOfDay
		prevStart = start.AddDate(0, -1, 0)
		prevEnd = end.AddDate(0, -1, 0)
		resolvedKey = "this_month"

	case "monthly":
		// Current year to date monthly window
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		end = endOfDay
		prevStart = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)
		prevEnd = time.Date(now.Year()-1, 12, 31, 23, 59, 59, 999999999, time.UTC)
		resolvedKey = "monthly"

	case "last_month":
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		start = firstOfThisMonth.AddDate(0, -1, 0)
		end = firstOfThisMonth.Add(-time.Nanosecond)
		prevStart = start.AddDate(0, -1, 0)
		prevEnd = start.Add(-time.Nanosecond)
		resolvedKey = "last_month"

	case "this_year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		end = endOfDay
		prevStart = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)
		prevEnd = end.AddDate(-1, 0, 0)
		resolvedKey = "this_year"

	case "last_year":
		start = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(now.Year()-1, 12, 31, 23, 59, 59, 999999999, time.UTC)
		prevStart = time.Date(now.Year()-2, 1, 1, 0, 0, 0, 0, time.UTC)
		prevEnd = time.Date(now.Year()-2, 12, 31, 23, 59, 59, 999999999, time.UTC)
		resolvedKey = "last_year"

	case "yearly":
		// Past 5 years
		start = time.Date(now.Year()-4, 1, 1, 0, 0, 0, 0, time.UTC)
		end = endOfDay
		prevStart = time.Date(now.Year()-9, 1, 1, 0, 0, 0, 0, time.UTC)
		prevEnd = time.Date(now.Year()-5, 12, 31, 23, 59, 59, 999999999, time.UTC)
		resolvedKey = "yearly"

	case "custom":
		tStart, err1 := time.Parse("2006-01-02", fromStr)
		tEnd, err2 := time.Parse("2006-01-02", toStr)
		if err1 == nil && err2 == nil {
			start = time.Date(tStart.Year(), tStart.Month(), tStart.Day(), 0, 0, 0, 0, time.UTC)
			end = time.Date(tEnd.Year(), tEnd.Month(), tEnd.Day(), 23, 59, 59, 999999999, time.UTC)
			duration := end.Sub(start)
			prevEnd = start.Add(-time.Nanosecond)
			prevStart = prevEnd.Add(-duration)
			resolvedKey = "custom"
			return
		}
		// Single individual day e.g. "2026-11-09"
		if err1 == nil && toStr == "" {
			start = time.Date(tStart.Year(), tStart.Month(), tStart.Day(), 0, 0, 0, 0, time.UTC)
			end = time.Date(tStart.Year(), tStart.Month(), tStart.Day(), 23, 59, 59, 999999999, time.UTC)
			prevStart = start.AddDate(0, 0, -1)
			prevEnd = end.AddDate(0, 0, -1)
			resolvedKey = "day"
			return
		}
		// Selected month e.g. "2026-10"
		if tMonth, errM := time.Parse("2006-01", fromStr); errM == nil {
			start = time.Date(tMonth.Year(), tMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
			prevStart = start.AddDate(0, -1, 0)
			prevEnd = start.Add(-time.Nanosecond)
			resolvedKey = "month"
			return
		}
		// Selected year e.g. "2026"
		if tYear, errY := time.Parse("2006", fromStr); errY == nil {
			start = time.Date(tYear.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			end = time.Date(tYear.Year(), 12, 31, 23, 59, 59, 999999999, time.UTC)
			prevStart = time.Date(tYear.Year()-1, 1, 1, 0, 0, 0, 0, time.UTC)
			prevEnd = time.Date(tYear.Year()-1, 12, 31, 23, 59, 59, 999999999, time.UTC)
			resolvedKey = "year"
			return
		}
		fallthrough

	default: // Default: 30d
		start = startOfDay.AddDate(0, 0, -29)
		end = endOfDay
		prevStart = start.AddDate(0, 0, -30)
		prevEnd = end.AddDate(0, 0, -30)
		resolvedKey = "30d"
	}

	return
}

func calcGrowth(curr, prev decimal.Decimal) decimal.Decimal {
	if prev.IsZero() {
		if curr.IsPositive() {
			return decimal.NewFromInt(100)
		}
		return decimal.Zero
	}
	return curr.Sub(prev).Div(prev).Mul(decimal.NewFromInt(100)).Round(2)
}

func calcGrowthInt(curr, prev int64) decimal.Decimal {
	if prev == 0 {
		if curr > 0 {
			return decimal.NewFromInt(100)
		}
		return decimal.Zero
	}
	diff := decimal.NewFromInt(curr - prev)
	prevDec := decimal.NewFromInt(prev)
	return diff.Div(prevDec).Mul(decimal.NewFromInt(100)).Round(2)
}

func (s *AnalyticsService) GetOverview(tenantID uuid.UUID, rangeKey, fromStr, toStr string) (*dto.AnalyticsOverviewResponse, error) {
	start, end, prevStart, prevEnd, resolvedKey := s.ResolveDateWindow(rangeKey, fromStr, toStr)

	currSales, err := s.repo.GetSalesAggregates(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	currItems, err := s.repo.GetItemsAggregates(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	prevSales, err := s.repo.GetSalesAggregates(tenantID, prevStart, prevEnd)
	if err != nil {
		return nil, err
	}

	grossProfit := currSales.NetSales.Sub(currItems.EstimatedCost)
	profitMargin := decimal.Zero
	if currSales.NetSales.IsPositive() {
		profitMargin = grossProfit.Div(currSales.NetSales).Mul(decimal.NewFromInt(100)).Round(2)
	}

	aov := decimal.Zero
	if currSales.TotalBills > 0 {
		aov = currSales.GrossSales.Div(decimal.NewFromInt(currSales.TotalBills)).Round(2)
	}

	salesGrowth := calcGrowth(currSales.GrossSales, prevSales.GrossSales)
	billsGrowth := calcGrowthInt(currSales.TotalBills, prevSales.TotalBills)

	return &dto.AnalyticsOverviewResponse{
		DateRange:       resolvedKey,
		StartDate:       start,
		EndDate:         end,
		GrossSales:      currSales.GrossSales.Round(2),
		TotalDiscounts:  currSales.TotalDiscounts.Round(2),
		TotalTax:        currSales.TotalTax.Round(2),
		NetSales:        currSales.NetSales.Round(2),
		EstimatedCost:   currItems.EstimatedCost.Round(2),
		GrossProfit:     grossProfit.Round(2),
		ProfitMarginPct: profitMargin,
		TotalBills:      currSales.TotalBills,
		TotalItemsSold:  currItems.TotalUnitsSold,
		AverageOrderVal: aov,
		PrevGrossSales:  prevSales.GrossSales.Round(2),
		PrevNetSales:    prevSales.NetSales.Round(2),
		PrevTotalBills:  prevSales.TotalBills,
		SalesGrowthPct:  salesGrowth,
		BillsGrowthPct:  billsGrowth,
	}, nil
}

// GetSalesPerformance produces the unified executive sales & revenue report
func (s *AnalyticsService) GetSalesPerformance(tenantID uuid.UUID, rangeKey, fromStr, toStr string) (*dto.SalesPerformanceSummary, error) {
	start, end, prevStart, prevEnd, resolvedKey := s.ResolveDateWindow(rangeKey, fromStr, toStr)

	currSales, err := s.repo.GetSalesAggregates(tenantID, start, end)
	if err != nil {
		return nil, err
	}
	currItems, err := s.repo.GetItemsAggregates(tenantID, start, end)
	if err != nil {
		return nil, err
	}
	prevSales, err := s.repo.GetSalesAggregates(tenantID, prevStart, prevEnd)
	if err != nil {
		return nil, err
	}
	prevItems, err := s.repo.GetItemsAggregates(tenantID, prevStart, prevEnd)
	if err != nil {
		return nil, err
	}

	returnsMetric, err := s.repo.GetReturnsMetrics(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	// Net Sales = Gross Sales - Total Discounts - Total Returns
	netSales := currSales.GrossSales.Sub(currSales.TotalDiscounts).Sub(returnsMetric.TotalReturns)
	if netSales.IsNegative() {
		netSales = decimal.Zero
	}

	// Amount collected is total money collected from completed bills
	amountCollected := currSales.GrossSales

	// COGS & Gross Profit
	cogs := currItems.EstimatedCost
	grossProfit := netSales.Sub(cogs)
	profitMargin := decimal.Zero
	if netSales.IsPositive() {
		profitMargin = grossProfit.Div(netSales).Mul(decimal.NewFromInt(100)).Round(2)
	}

	// Average Bill Value (AOV)
	aov := decimal.Zero
	if currSales.TotalBills > 0 {
		aov = netSales.Div(decimal.NewFromInt(currSales.TotalBills)).Round(2)
	}

	// Average Daily Sales
	durationDays := int64(end.Sub(start).Hours()/24) + 1
	if durationDays <= 0 {
		durationDays = 1
	}
	avgDaily := netSales.Div(decimal.NewFromInt(durationDays)).Round(2)

	// Best Selling Day
	dailyRows, _ := s.repo.GetDetailedDailyBreakdown(tenantID, start, end)
	var bestDay *dto.BestSellingMetric
	if len(dailyRows) > 0 {
		var maxSales decimal.Decimal
		var maxUnits int64
		var maxLabel string
		for _, row := range dailyRows {
			if row.GrossSales.GreaterThan(maxSales) {
				maxSales = row.GrossSales
				maxUnits = row.UnitsSold
				maxLabel = row.Date
			}
		}
		if !maxSales.IsZero() {
			bestDay = &dto.BestSellingMetric{
				Label: maxLabel,
				Sales: maxSales.Round(2),
				Units: maxUnits,
			}
		}
	}

	// Best Selling Month
	monthlyMap := make(map[string]*dto.BestSellingMetric)
	for _, row := range dailyRows {
		if len(row.Date) >= 7 {
			moKey := row.Date[:7]
			if metric, exists := monthlyMap[moKey]; exists {
				metric.Sales = metric.Sales.Add(row.GrossSales)
				metric.Units += row.UnitsSold
			} else {
				monthlyMap[moKey] = &dto.BestSellingMetric{
					Label: moKey,
					Sales: row.GrossSales,
					Units: row.UnitsSold,
				}
			}
		}
	}
	var bestMonth *dto.BestSellingMetric
	var maxMoSales decimal.Decimal
	for _, m := range monthlyMap {
		if m.Sales.GreaterThan(maxMoSales) {
			maxMoSales = m.Sales
			bestMonth = m
		}
	}

	// Previous period calculations
	prevReturns, _ := s.repo.GetReturnsMetrics(tenantID, prevStart, prevEnd)
	prevReturnVal := decimal.Zero
	if prevReturns != nil {
		prevReturnVal = prevReturns.TotalReturns
	}
	prevNetSales := prevSales.GrossSales.Sub(prevSales.TotalDiscounts).Sub(prevReturnVal)
	if prevNetSales.IsNegative() {
		prevNetSales = decimal.Zero
	}

	salesGrowth := calcGrowth(netSales, prevNetSales)
	billsGrowth := calcGrowthInt(currSales.TotalBills, prevSales.TotalBills)

	return &dto.SalesPerformanceSummary{
		PeriodType:      resolvedKey,
		StartDate:       start,
		EndDate:         end,
		GrossSales:      currSales.GrossSales.Round(2),
		TotalDiscounts:  currSales.TotalDiscounts.Round(2),
		TotalReturns:    returnsMetric.TotalReturns.Round(2),
		TotalTax:        currSales.TotalTax.Round(2),
		NetSales:        netSales.Round(2),
		AmountCollected: amountCollected.Round(2),
		CostOfGoodsSold: cogs.Round(2),
		GrossProfit:     grossProfit.Round(2),
		ProfitMarginPct: profitMargin,
		TotalBills:      currSales.TotalBills,
		TotalUnitsSold:  currItems.TotalUnitsSold,
		AverageBillVal:  aov,
		AverageDailyVal: avgDaily,
		BestSellingDay:  bestDay,
		BestSellingMo:   bestMonth,
		PrevGrossSales:  prevSales.GrossSales.Round(2),
		PrevNetSales:    prevNetSales.Round(2),
		PrevUnitsSold:   prevItems.TotalUnitsSold,
		PrevTotalBills:  prevSales.TotalBills,
		SalesGrowthPct:  salesGrowth,
		BillsGrowthPct:  billsGrowth,
	}, nil
}

// GetDetailedDailyBreakdown returns daily figures
func (s *AnalyticsService) GetDetailedDailyBreakdown(tenantID uuid.UUID, rangeKey, fromStr, toStr string) ([]dto.DailyBreakdownItem, error) {
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)
	rows, err := s.repo.GetDetailedDailyBreakdown(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	items := make([]dto.DailyBreakdownItem, len(rows))
	for i, r := range rows {
		items[i] = dto.DailyBreakdownItem{
			Date:       r.Date,
			Bills:      r.Bills,
			UnitsSold:  r.UnitsSold,
			GrossSales: r.GrossSales.Round(2),
			Discounts:  r.Discounts.Round(2),
			Returns:    r.Returns.Round(2),
			TaxAmount:  r.TaxAmount.Round(2),
			NetSales:   r.NetSales.Round(2),
		}
	}
	return items, nil
}

// GetMonthlyBreakdown groups sales by year and month
func (s *AnalyticsService) GetMonthlyBreakdown(tenantID uuid.UUID, rangeKey, fromStr, toStr string) ([]dto.MonthlyBreakdownItem, error) {
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)
	dailyRows, err := s.repo.GetDetailedDailyBreakdown(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	type mKey struct {
		Year  int
		Month int
	}
	grouped := make(map[mKey]*dto.MonthlyBreakdownItem)

	for _, d := range dailyRows {
		t, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			continue
		}
		k := mKey{Year: t.Year(), Month: int(t.Month())}
		item, ok := grouped[k]
		if !ok {
			item = &dto.MonthlyBreakdownItem{
				Year:      t.Year(),
				Month:     int(t.Month()),
				MonthName: t.Month().String(),
			}
			grouped[k] = item
		}
		item.Bills += d.Bills
		item.UnitsSold += d.UnitsSold
		item.GrossSales = item.GrossSales.Add(d.GrossSales)
		item.Discounts = item.Discounts.Add(d.Discounts)
		item.Returns = item.Returns.Add(d.Returns)
		item.NetSales = item.NetSales.Add(d.NetSales)
	}

	// Sort chronologically
	var keys []mKey
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Year != keys[j].Year {
			return keys[i].Year < keys[j].Year
		}
		return keys[i].Month < keys[j].Month
	})

	results := make([]dto.MonthlyBreakdownItem, len(keys))
	for i, k := range keys {
		item := grouped[k]
		item.GrossSales = item.GrossSales.Round(2)
		item.Discounts = item.Discounts.Round(2)
		item.Returns = item.Returns.Round(2)
		item.NetSales = item.NetSales.Round(2)
		results[i] = *item
	}
	return results, nil
}

// GetYearlyBreakdown groups sales by year
func (s *AnalyticsService) GetYearlyBreakdown(tenantID uuid.UUID, rangeKey, fromStr, toStr string) ([]dto.YearlyBreakdownItem, error) {
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)
	dailyRows, err := s.repo.GetDetailedDailyBreakdown(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	grouped := make(map[int]*dto.YearlyBreakdownItem)
	for _, d := range dailyRows {
		t, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			continue
		}
		year := t.Year()
		item, ok := grouped[year]
		if !ok {
			item = &dto.YearlyBreakdownItem{Year: year}
			grouped[year] = item
		}
		item.Bills += d.Bills
		item.UnitsSold += d.UnitsSold
		item.GrossSales = item.GrossSales.Add(d.GrossSales)
		item.Discounts = item.Discounts.Add(d.Discounts)
		item.Returns = item.Returns.Add(d.Returns)
		item.NetSales = item.NetSales.Add(d.NetSales)
	}

	var years []int
	for y := range grouped {
		years = append(years, y)
	}
	sort.Ints(years)

	results := make([]dto.YearlyBreakdownItem, len(years))
	for i, y := range years {
		item := grouped[y]
		item.GrossSales = item.GrossSales.Round(2)
		item.Discounts = item.Discounts.Round(2)
		item.Returns = item.Returns.Round(2)
		item.NetSales = item.NetSales.Round(2)
		results[i] = *item
	}
	return results, nil
}

// GetDetailedProductPerformance returns product intelligence categorized by sales performance
func (s *AnalyticsService) GetDetailedProductPerformance(tenantID uuid.UUID, rangeKey, fromStr, toStr string) (*dto.ProductPerformanceResponse, error) {
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)
	deepRows, err := s.repo.GetDeepProductSales(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	var allItems []dto.DetailedProductPerformanceItem
	var topSelling []dto.DetailedProductPerformanceItem
	var slowMoving []dto.DetailedProductPerformanceItem
	var zeroSales []dto.DetailedProductPerformanceItem

	for _, r := range deepRows {
		pID := r.ProductID
		profit := r.NetSales.Sub(r.CostOfGoods).Round(2)
		margin := decimal.Zero
		if r.NetSales.IsPositive() {
			margin = profit.Div(r.NetSales).Mul(decimal.NewFromInt(100)).Round(2)
		}

		status := "ACTIVE"
		if r.UnitsSold == 0 {
			status = "ZERO_SALES"
		} else if r.UnitsSold >= 10 {
			status = "BEST_SELLER"
		} else if r.UnitsSold < 3 {
			status = "SLOW_MOVING"
		}

		item := dto.DetailedProductPerformanceItem{
			ProductID:    &pID,
			ProductName:  r.ProductName,
			ProductCode:  r.ProductCode,
			SKU:          r.SKU,
			Barcode:      r.Barcode,
			BrandName:    r.BrandName,
			CategoryName: r.CategoryName,
			Size:         r.Size,
			Colour:       r.Colour,
			UnitsSold:    r.UnitsSold,
			GrossSales:   r.GrossSales.Round(2),
			Discounts:    r.Discounts.Round(2),
			Returns:      decimal.Zero,
			NetSales:     r.NetSales.Round(2),
			CostOfGoods:  r.CostOfGoods.Round(2),
			GrossProfit:  profit,
			ProfitMargin: margin,
			StockOnHand:  r.CurrentStock,
			Status:       status,
		}

		allItems = append(allItems, item)

		switch status {
		case "ZERO_SALES":
			zeroSales = append(zeroSales, item)
		case "SLOW_MOVING":
			slowMoving = append(slowMoving, item)
		case "BEST_SELLER":
			topSelling = append(topSelling, item)
		}
	}

	// Sort top selling by UnitsSold descending
	sort.Slice(topSelling, func(i, j int) bool {
		return topSelling[i].UnitsSold > topSelling[j].UnitsSold
	})
	// Sort zero sales by CurrentStock descending (highest idle capital first)
	sort.Slice(zeroSales, func(i, j int) bool {
		return zeroSales[i].StockOnHand > zeroSales[j].StockOnHand
	})

	return &dto.ProductPerformanceResponse{
		TopSelling:    topSelling,
		SlowMoving:    slowMoving,
		ZeroSales:     zeroSales,
		AllProducts:   allItems,
		TotalProducts: int64(len(allItems)),
	}, nil
}

func (s *AnalyticsService) GetTrends(tenantID uuid.UUID, rangeKey, fromStr, toStr string) ([]dto.TrendPoint, error) {
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)
	rows, err := s.repo.GetDailyTrends(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	points := make([]dto.TrendPoint, len(rows))
	for i, r := range rows {
		t, _ := time.Parse("2006-01-02", r.Period)
		points[i] = dto.TrendPoint{
			Period:    r.Period,
			Timestamp: t,
			Sales:     r.Sales.Round(2),
			NetSales:  r.NetSales.Round(2),
			Bills:     r.Bills,
			UnitsSold: r.UnitsSold,
		}
	}
	return points, nil
}

func (s *AnalyticsService) GetPaymentMethods(tenantID uuid.UUID, rangeKey, fromStr, toStr string) ([]dto.PaymentMethodAnalytics, error) {
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)
	rows, err := s.repo.GetPaymentMethodDistribution(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	totalSales := decimal.Zero
	for _, r := range rows {
		totalSales = totalSales.Add(r.TotalSales)
	}

	res := make([]dto.PaymentMethodAnalytics, len(rows))
	for i, r := range rows {
		pct := decimal.Zero
		if totalSales.IsPositive() {
			pct = r.TotalSales.Div(totalSales).Mul(decimal.NewFromInt(100)).Round(2)
		}
		res[i] = dto.PaymentMethodAnalytics{
			Method:     r.Method,
			TotalSales: r.TotalSales.Round(2),
			BillCount:  r.BillCount,
			Percentage: pct,
		}
	}
	return res, nil
}

func (s *AnalyticsService) GetTopProducts(tenantID uuid.UUID, rangeKey, fromStr, toStr string, limit int) ([]dto.ProductSalesAnalytics, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)
	rows, err := s.repo.GetTopProducts(tenantID, start, end, limit)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ProductSalesAnalytics, len(rows))
	for i, r := range rows {
		profit := r.Revenue.Sub(r.EstCost).Round(2)
		res[i] = dto.ProductSalesAnalytics{
			ProductName:  r.ProductName,
			BrandName:    r.BrandName,
			CategoryName: r.CategoryName,
			UnitsSold:    r.UnitsSold,
			Revenue:      r.Revenue.Round(2),
			GrossProfit:  profit,
		}
	}
	return res, nil
}

func (s *AnalyticsService) GetCategorySales(tenantID uuid.UUID, rangeKey, fromStr, toStr string) ([]dto.CategorySalesAnalytics, error) {
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)
	rows, err := s.repo.GetCategorySales(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	totalRev := decimal.Zero
	for _, r := range rows {
		totalRev = totalRev.Add(r.Revenue)
	}

	res := make([]dto.CategorySalesAnalytics, len(rows))
	for i, r := range rows {
		pct := decimal.Zero
		if totalRev.IsPositive() {
			pct = r.Revenue.Div(totalRev).Mul(decimal.NewFromInt(100)).Round(2)
		}
		res[i] = dto.CategorySalesAnalytics{
			CategoryName: r.CategoryName,
			Revenue:      r.Revenue.Round(2),
			UnitsSold:    r.UnitsSold,
			Percentage:   pct,
		}
	}
	return res, nil
}

func (s *AnalyticsService) GetFootwearIntelligence(tenantID uuid.UUID, rangeKey, fromStr, toStr string) (*dto.FootwearIntelligenceResponse, error) {
	start, end, _, _, _ := s.ResolveDateWindow(rangeKey, fromStr, toStr)

	sizes, err := s.repo.GetTopSizes(tenantID, start, end, 10)
	if err != nil {
		return nil, err
	}
	colours, err := s.repo.GetTopColours(tenantID, start, end, 10)
	if err != nil {
		return nil, err
	}
	brands, err := s.repo.GetTopBrands(tenantID, start, end, 10)
	if err != nil {
		return nil, err
	}

	mapAttributes := func(attrType string, items []repositories.RawAttributeSales) []dto.VariantAttributeAnalytics {
		totalUnits := int64(0)
		for _, it := range items {
			totalUnits += it.UnitsSold
		}
		out := make([]dto.VariantAttributeAnalytics, len(items))
		for i, it := range items {
			pct := decimal.Zero
			if totalUnits > 0 {
				pct = decimal.NewFromInt(it.UnitsSold).Div(decimal.NewFromInt(totalUnits)).Mul(decimal.NewFromInt(100)).Round(2)
			}
			out[i] = dto.VariantAttributeAnalytics{
				AttributeType: attrType,
				AttributeVal:  it.AttributeVal,
				UnitsSold:     it.UnitsSold,
				Revenue:       it.Revenue.Round(2),
				Percentage:    pct,
			}
		}
		return out
	}

	return &dto.FootwearIntelligenceResponse{
		TopSizes:   mapAttributes("size", sizes),
		TopColours: mapAttributes("colour", colours),
		TopBrands:  mapAttributes("brand", brands),
	}, nil
}

// ExportComprehensiveExcel generates an 8-worksheet Excel workbook
func (s *AnalyticsService) ExportComprehensiveExcel(tenantID uuid.UUID, rangeKey, fromStr, toStr string) ([]byte, string, error) {
	summary, err := s.GetSalesPerformance(tenantID, rangeKey, fromStr, toStr)
	if err != nil {
		return nil, "", err
	}

	dailyBreakdown, _ := s.GetDetailedDailyBreakdown(tenantID, rangeKey, fromStr, toStr)
	monthlyBreakdown, _ := s.GetMonthlyBreakdown(tenantID, rangeKey, fromStr, toStr)
	yearlyBreakdown, _ := s.GetYearlyBreakdown(tenantID, rangeKey, fromStr, toStr)
	productPerf, _ := s.GetDetailedProductPerformance(tenantID, rangeKey, fromStr, toStr)
	categorySales, _ := s.GetCategorySales(tenantID, rangeKey, fromStr, toStr)
	footwearIntel, _ := s.GetFootwearIntelligence(tenantID, rangeKey, fromStr, toStr)
	paymentMethods, _ := s.GetPaymentMethods(tenantID, rangeKey, fromStr, toStr)
	rawTransactions, _ := s.repo.GetRawTransactions(tenantID, summary.StartDate, summary.EndDate)

	f := excelize.NewFile()
	defer f.Close()

	// Setup professional styles
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF", Size: 11, Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#0F172A"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "#475569", Style: 1},
			{Type: "top", Color: "#475569", Style: 1},
			{Type: "bottom", Color: "#475569", Style: 1},
			{Type: "right", Color: "#475569", Style: 1},
		},
	})

	cardHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF", Size: 10, Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#2563EB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})

	totalRowStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#0F172A", Size: 11, Family: "Calibri"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E2E8F0"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "top", Color: "#0F172A", Style: 1},
			{Type: "bottom", Color: "#0F172A", Style: 6}, // double line
		},
	})

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#0F172A", Size: 16, Family: "Calibri"},
	})

	subTitleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Italic: true, Color: "#64748B", Size: 10, Family: "Calibri"},
	})

	// 1. SUMMARY SHEET
	f.SetSheetName("Sheet1", "Summary")
	f.SetCellValue("Summary", "A1", "KadaiIQ - Sales Performance & Revenue Analytics")
	f.SetCellStyle("Summary", "A1", "A1", titleStyle)
	f.SetCellValue("Summary", "A2", fmt.Sprintf("Period: %s (%s to %s) | Generated: %s UTC",
		strings.ToUpper(summary.PeriodType),
		summary.StartDate.Format("02 Jan 2006"),
		summary.EndDate.Format("02 Jan 2006"),
		time.Now().UTC().Format("02 Jan 2006 15:04:05")))
	f.SetCellStyle("Summary", "A2", "A2", subTitleStyle)

	// Summary KPI Table
	f.SetCellValue("Summary", "A4", "Key Financial Metric")
	f.SetCellValue("Summary", "B4", "Amount / Value")
	f.SetCellStyle("Summary", "A4", "B4", cardHeaderStyle)

	metrics := [][]interface{}{
		{"Gross Sales", "₹ " + summary.GrossSales.StringFixed(2)},
		{"Total Discounts Given", "₹ " + summary.TotalDiscounts.StringFixed(2)},
		{"Sales Returns", "₹ " + summary.TotalReturns.StringFixed(2)},
		{"GST / Taxes Collected", "₹ " + summary.TotalTax.StringFixed(2)},
		{"Net Sales (Gross - Discounts - Returns)", "₹ " + summary.NetSales.StringFixed(2)},
		{"Actual Amount Collected (Tenders)", "₹ " + summary.AmountCollected.StringFixed(2)},
		{"Cost of Goods Sold (COGS)", "₹ " + summary.CostOfGoodsSold.StringFixed(2)},
		{"Gross Profit (Net Sales - COGS)", "₹ " + summary.GrossProfit.StringFixed(2)},
		{"Gross Profit Margin (%)", summary.ProfitMarginPct.StringFixed(2) + " %"},
		{"Completed Invoices / Bills", summary.TotalBills},
		{"Total Footwear Products / Units Sold", summary.TotalUnitsSold},
		{"Average Order Value (AOV)", "₹ " + summary.AverageBillVal.StringFixed(2)},
		{"Average Daily Revenue", "₹ " + summary.AverageDailyVal.StringFixed(2)},
	}

	for idx, row := range metrics {
		rNum := 5 + idx
		f.SetCellValue("Summary", fmt.Sprintf("A%d", rNum), row[0])
		f.SetCellValue("Summary", fmt.Sprintf("B%d", rNum), row[1])
	}

	f.SetCellValue("Summary", "A19", "Growth & Period Benchmarks")
	f.SetCellValue("Summary", "B19", "Metric")
	f.SetCellStyle("Summary", "A19", "B19", cardHeaderStyle)

	growthData := [][]interface{}{
		{"Prior Period Net Sales", "₹ " + summary.PrevNetSales.StringFixed(2)},
		{"Sales Growth Rate vs Prior Period", summary.SalesGrowthPct.StringFixed(2) + " %"},
		{"Prior Period Invoices Count", summary.PrevTotalBills},
		{"Invoices Growth Rate", summary.BillsGrowthPct.StringFixed(2) + " %"},
	}
	if summary.BestSellingDay != nil {
		growthData = append(growthData, []interface{}{"Best-Selling Day", fmt.Sprintf("%s (₹ %s, %d units)", summary.BestSellingDay.Label, summary.BestSellingDay.Sales.StringFixed(2), summary.BestSellingDay.Units)})
	}
	if summary.BestSellingMo != nil {
		growthData = append(growthData, []interface{}{"Best-Selling Month", fmt.Sprintf("%s (₹ %s, %d units)", summary.BestSellingMo.Label, summary.BestSellingMo.Sales.StringFixed(2), summary.BestSellingMo.Units)})
	}

	for idx, row := range growthData {
		rNum := 20 + idx
		f.SetCellValue("Summary", fmt.Sprintf("A%d", rNum), row[0])
		f.SetCellValue("Summary", fmt.Sprintf("B%d", rNum), row[1])
	}
	f.SetColWidth("Summary", "A", "A", 42)
	f.SetColWidth("Summary", "B", "B", 35)

	// 2. DAILY SALES SHEET
	f.NewSheet("Daily Sales")
	dHeaders := []string{"Date", "Bills Count", "Units Sold", "Gross Sales (₹)", "Discounts (₹)", "Returns (₹)", "Tax (₹)", "Net Sales (₹)"}
	for c, h := range dHeaders {
		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		f.SetCellValue("Daily Sales", cell, h)
	}
	f.SetCellStyle("Daily Sales", "A1", "H1", headerStyle)

	dRow := 2
	var dTotBills, dTotUnits int64
	dTotGross, dTotDisc, dTotRet, dTotTax, dTotNet := decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero

	for _, d := range dailyBreakdown {
		f.SetCellValue("Daily Sales", fmt.Sprintf("A%d", dRow), d.Date)
		f.SetCellValue("Daily Sales", fmt.Sprintf("B%d", dRow), d.Bills)
		f.SetCellValue("Daily Sales", fmt.Sprintf("C%d", dRow), d.UnitsSold)
		f.SetCellValue("Daily Sales", fmt.Sprintf("D%d", dRow), d.GrossSales.InexactFloat64())
		f.SetCellValue("Daily Sales", fmt.Sprintf("E%d", dRow), d.Discounts.InexactFloat64())
		f.SetCellValue("Daily Sales", fmt.Sprintf("F%d", dRow), d.Returns.InexactFloat64())
		f.SetCellValue("Daily Sales", fmt.Sprintf("G%d", dRow), d.TaxAmount.InexactFloat64())
		f.SetCellValue("Daily Sales", fmt.Sprintf("H%d", dRow), d.NetSales.InexactFloat64())

		dTotBills += d.Bills
		dTotUnits += d.UnitsSold
		dTotGross = dTotGross.Add(d.GrossSales)
		dTotDisc = dTotDisc.Add(d.Discounts)
		dTotRet = dTotRet.Add(d.Returns)
		dTotTax = dTotTax.Add(d.TaxAmount)
		dTotNet = dTotNet.Add(d.NetSales)
		dRow++
	}

	// Total Row
	f.SetCellValue("Daily Sales", fmt.Sprintf("A%d", dRow), "TOTAL")
	f.SetCellValue("Daily Sales", fmt.Sprintf("B%d", dRow), dTotBills)
	f.SetCellValue("Daily Sales", fmt.Sprintf("C%d", dRow), dTotUnits)
	f.SetCellValue("Daily Sales", fmt.Sprintf("D%d", dRow), dTotGross.InexactFloat64())
	f.SetCellValue("Daily Sales", fmt.Sprintf("E%d", dRow), dTotDisc.InexactFloat64())
	f.SetCellValue("Daily Sales", fmt.Sprintf("F%d", dRow), dTotRet.InexactFloat64())
	f.SetCellValue("Daily Sales", fmt.Sprintf("G%d", dRow), dTotTax.InexactFloat64())
	f.SetCellValue("Daily Sales", fmt.Sprintf("H%d", dRow), dTotNet.InexactFloat64())
	f.SetCellStyle("Daily Sales", fmt.Sprintf("A%d", dRow), fmt.Sprintf("H%d", dRow), totalRowStyle)
	f.SetColWidth("Daily Sales", "A", "H", 18)

	// 3. MONTHLY SALES SHEET (Both names for strict backwards compatibility)
	f.NewSheet("Monthly Sales")
	f.NewSheet("Monthly Breakdown")
	mHeaders := []string{"Year", "Month", "Bills", "Units Sold", "Gross Sales (₹)", "Discounts (₹)", "Returns (₹)", "Net Sales (₹)"}
	for c, h := range mHeaders {
		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		f.SetCellValue("Monthly Sales", cell, h)
		f.SetCellValue("Monthly Breakdown", cell, h)
	}
	f.SetCellStyle("Monthly Sales", "A1", "H1", headerStyle)
	f.SetCellStyle("Monthly Breakdown", "A1", "H1", headerStyle)

	mRow := 2
	for _, m := range monthlyBreakdown {
		for _, sName := range []string{"Monthly Sales", "Monthly Breakdown"} {
			f.SetCellValue(sName, fmt.Sprintf("A%d", mRow), m.Year)
			f.SetCellValue(sName, fmt.Sprintf("B%d", mRow), m.MonthName)
			f.SetCellValue(sName, fmt.Sprintf("C%d", mRow), m.Bills)
			f.SetCellValue(sName, fmt.Sprintf("D%d", mRow), m.UnitsSold)
			f.SetCellValue(sName, fmt.Sprintf("E%d", mRow), m.GrossSales.InexactFloat64())
			f.SetCellValue(sName, fmt.Sprintf("F%d", mRow), m.Discounts.InexactFloat64())
			f.SetCellValue(sName, fmt.Sprintf("G%d", mRow), m.Returns.InexactFloat64())
			f.SetCellValue(sName, fmt.Sprintf("H%d", mRow), m.NetSales.InexactFloat64())
		}
		mRow++
	}
	f.SetColWidth("Monthly Sales", "A", "H", 16)
	f.SetColWidth("Monthly Breakdown", "A", "H", 16)

	// 4. YEARLY SUMMARY SHEET
	f.NewSheet("Yearly Summary")
	yHeaders := []string{"Year", "Bills / Invoices", "Units Sold", "Gross Sales (₹)", "Discounts (₹)", "Returns (₹)", "Net Sales (₹)"}
	for c, h := range yHeaders {
		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		f.SetCellValue("Yearly Summary", cell, h)
	}
	f.SetCellStyle("Yearly Summary", "A1", "G1", headerStyle)

	yRow := 2
	var yTotBills, yTotUnits int64
	yTotGross, yTotDisc, yTotRet, yTotNet := decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero
	for _, y := range yearlyBreakdown {
		f.SetCellValue("Yearly Summary", fmt.Sprintf("A%d", yRow), y.Year)
		f.SetCellValue("Yearly Summary", fmt.Sprintf("B%d", yRow), y.Bills)
		f.SetCellValue("Yearly Summary", fmt.Sprintf("C%d", yRow), y.UnitsSold)
		f.SetCellValue("Yearly Summary", fmt.Sprintf("D%d", yRow), y.GrossSales.InexactFloat64())
		f.SetCellValue("Yearly Summary", fmt.Sprintf("E%d", yRow), y.Discounts.InexactFloat64())
		f.SetCellValue("Yearly Summary", fmt.Sprintf("F%d", yRow), y.Returns.InexactFloat64())
		f.SetCellValue("Yearly Summary", fmt.Sprintf("G%d", yRow), y.NetSales.InexactFloat64())

		yTotBills += y.Bills
		yTotUnits += y.UnitsSold
		yTotGross = yTotGross.Add(y.GrossSales)
		yTotDisc = yTotDisc.Add(y.Discounts)
		yTotRet = yTotRet.Add(y.Returns)
		yTotNet = yTotNet.Add(y.NetSales)
		yRow++
	}
	// Total row for yearly summary
	f.SetCellValue("Yearly Summary", fmt.Sprintf("A%d", yRow), "TOTAL")
	f.SetCellValue("Yearly Summary", fmt.Sprintf("B%d", yRow), yTotBills)
	f.SetCellValue("Yearly Summary", fmt.Sprintf("C%d", yRow), yTotUnits)
	f.SetCellValue("Yearly Summary", fmt.Sprintf("D%d", yRow), yTotGross.InexactFloat64())
	f.SetCellValue("Yearly Summary", fmt.Sprintf("E%d", yRow), yTotDisc.InexactFloat64())
	f.SetCellValue("Yearly Summary", fmt.Sprintf("F%d", yRow), yTotRet.InexactFloat64())
	f.SetCellValue("Yearly Summary", fmt.Sprintf("G%d", yRow), yTotNet.InexactFloat64())
	f.SetCellStyle("Yearly Summary", fmt.Sprintf("A%d", yRow), fmt.Sprintf("G%d", yRow), totalRowStyle)
	f.SetColWidth("Yearly Summary", "A", "G", 18)

	// 4. PRODUCT SALES SHEET
	f.NewSheet("Product Sales")
	pHeaders := []string{"Product Code", "Product Name", "Brand", "Category", "Size", "Colour", "Barcode", "Units Sold", "Gross Sales (₹)", "Discounts (₹)", "Net Sales (₹)", "COGS (₹)", "Gross Profit (₹)", "Margin %", "Stock on Hand", "Status"}
	for c, h := range pHeaders {
		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		f.SetCellValue("Product Sales", cell, h)
	}
	f.SetCellStyle("Product Sales", "A1", "P1", headerStyle)

	pRow := 2
	if productPerf != nil {
		for _, p := range productPerf.AllProducts {
			f.SetCellValue("Product Sales", fmt.Sprintf("A%d", pRow), p.ProductCode)
			f.SetCellValue("Product Sales", fmt.Sprintf("B%d", pRow), p.ProductName)
			f.SetCellValue("Product Sales", fmt.Sprintf("C%d", pRow), p.BrandName)
			f.SetCellValue("Product Sales", fmt.Sprintf("D%d", pRow), p.CategoryName)
			f.SetCellValue("Product Sales", fmt.Sprintf("E%d", pRow), p.Size)
			f.SetCellValue("Product Sales", fmt.Sprintf("F%d", pRow), p.Colour)
			f.SetCellValue("Product Sales", fmt.Sprintf("G%d", pRow), p.Barcode)
			f.SetCellValue("Product Sales", fmt.Sprintf("H%d", pRow), p.UnitsSold)
			f.SetCellValue("Product Sales", fmt.Sprintf("I%d", pRow), p.GrossSales.InexactFloat64())
			f.SetCellValue("Product Sales", fmt.Sprintf("J%d", pRow), p.Discounts.InexactFloat64())
			f.SetCellValue("Product Sales", fmt.Sprintf("K%d", pRow), p.NetSales.InexactFloat64())
			f.SetCellValue("Product Sales", fmt.Sprintf("L%d", pRow), p.CostOfGoods.InexactFloat64())
			f.SetCellValue("Product Sales", fmt.Sprintf("M%d", pRow), p.GrossProfit.InexactFloat64())
			f.SetCellValue("Product Sales", fmt.Sprintf("N%d", pRow), p.ProfitMargin.InexactFloat64())
			f.SetCellValue("Product Sales", fmt.Sprintf("O%d", pRow), p.StockOnHand)
			f.SetCellValue("Product Sales", fmt.Sprintf("P%d", pRow), p.Status)
			pRow++
		}
	}
	f.SetColWidth("Product Sales", "A", "P", 16)
	f.SetColWidth("Product Sales", "B", "B", 26)

	// 5. CATEGORY & BRAND SALES SHEET
	f.NewSheet("Category & Brand Sales")
	f.SetCellValue("Category & Brand Sales", "A1", "Product Category Breakdown")
	f.SetCellValue("Category & Brand Sales", "B1", "Units Sold")
	f.SetCellValue("Category & Brand Sales", "C1", "Revenue (₹)")
	f.SetCellValue("Category & Brand Sales", "D1", "Revenue Share %")
	f.SetCellStyle("Category & Brand Sales", "A1", "D1", headerStyle)

	cbRow := 2
	for _, cat := range categorySales {
		f.SetCellValue("Category & Brand Sales", fmt.Sprintf("A%d", cbRow), cat.CategoryName)
		f.SetCellValue("Category & Brand Sales", fmt.Sprintf("B%d", cbRow), cat.UnitsSold)
		f.SetCellValue("Category & Brand Sales", fmt.Sprintf("C%d", cbRow), cat.Revenue.InexactFloat64())
		f.SetCellValue("Category & Brand Sales", fmt.Sprintf("D%d", cbRow), cat.Percentage.InexactFloat64())
		cbRow++
	}

	cbRow += 2
	f.SetCellValue("Category & Brand Sales", fmt.Sprintf("A%d", cbRow), "Brand Performance Breakdown")
	f.SetCellValue("Category & Brand Sales", fmt.Sprintf("B%d", cbRow), "Units Sold")
	f.SetCellValue("Category & Brand Sales", fmt.Sprintf("C%d", cbRow), "Revenue (₹)")
	f.SetCellValue("Category & Brand Sales", fmt.Sprintf("D%d", cbRow), "Revenue Share %")
	f.SetCellStyle("Category & Brand Sales", fmt.Sprintf("A%d", cbRow), fmt.Sprintf("D%d", cbRow), headerStyle)
	cbRow++

	if footwearIntel != nil {
		for _, br := range footwearIntel.TopBrands {
			f.SetCellValue("Category & Brand Sales", fmt.Sprintf("A%d", cbRow), br.AttributeVal)
			f.SetCellValue("Category & Brand Sales", fmt.Sprintf("B%d", cbRow), br.UnitsSold)
			f.SetCellValue("Category & Brand Sales", fmt.Sprintf("C%d", cbRow), br.Revenue.InexactFloat64())
			f.SetCellValue("Category & Brand Sales", fmt.Sprintf("D%d", cbRow), br.Percentage.InexactFloat64())
			cbRow++
		}
	}
	f.SetColWidth("Category & Brand Sales", "A", "D", 22)

	// 6. ZERO SALES ALERT SHEET
	f.NewSheet("Zero Sales Alert")
	zHeaders := []string{"Product Code", "Product Name", "Brand", "Category", "Size", "Colour", "Barcode", "Current Stock (Units)", "Capital Locked"}
	for c, h := range zHeaders {
		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		f.SetCellValue("Zero Sales Alert", cell, h)
	}
	f.SetCellStyle("Zero Sales Alert", "A1", "I1", headerStyle)

	zRow := 2
	if productPerf != nil {
		for _, z := range productPerf.ZeroSales {
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("A%d", zRow), z.ProductCode)
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("B%d", zRow), z.ProductName)
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("C%d", zRow), z.BrandName)
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("D%d", zRow), z.CategoryName)
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("E%d", zRow), z.Size)
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("F%d", zRow), z.Colour)
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("G%d", zRow), z.Barcode)
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("H%d", zRow), z.StockOnHand)
			f.SetCellValue("Zero Sales Alert", fmt.Sprintf("I%d", zRow), "Capital Idle")
			zRow++
		}
	}
	f.SetColWidth("Zero Sales Alert", "A", "I", 18)
	f.SetColWidth("Zero Sales Alert", "B", "B", 28)

	// 6. FOOTWEAR MATRIX SHEET
	f.NewSheet("Footwear Matrix")
	f.SetCellValue("Footwear Matrix", "A1", "Footwear Size Breakdown")
	f.SetCellValue("Footwear Matrix", "B1", "Units Sold")
	f.SetCellValue("Footwear Matrix", "C1", "Revenue (₹)")
	f.SetCellValue("Footwear Matrix", "D1", "Sales Share %")
	f.SetCellStyle("Footwear Matrix", "A1", "D1", headerStyle)

	fmRow := 2
	if footwearIntel != nil {
		for _, s := range footwearIntel.TopSizes {
			f.SetCellValue("Footwear Matrix", fmt.Sprintf("A%d", fmRow), "Size "+s.AttributeVal)
			f.SetCellValue("Footwear Matrix", fmt.Sprintf("B%d", fmRow), s.UnitsSold)
			f.SetCellValue("Footwear Matrix", fmt.Sprintf("C%d", fmRow), s.Revenue.InexactFloat64())
			f.SetCellValue("Footwear Matrix", fmt.Sprintf("D%d", fmRow), s.Percentage.InexactFloat64())
			fmRow++
		}

		fmRow += 2
		f.SetCellValue("Footwear Matrix", fmt.Sprintf("A%d", fmRow), "Footwear Colour Breakdown")
		f.SetCellValue("Footwear Matrix", fmt.Sprintf("B%d", fmRow), "Units Sold")
		f.SetCellValue("Footwear Matrix", fmt.Sprintf("C%d", fmRow), "Revenue (₹)")
		f.SetCellValue("Footwear Matrix", fmt.Sprintf("D%d", fmRow), "Sales Share %")
		f.SetCellStyle("Footwear Matrix", fmt.Sprintf("A%d", fmRow), fmt.Sprintf("D%d", fmRow), headerStyle)
		fmRow++

		for _, c := range footwearIntel.TopColours {
			f.SetCellValue("Footwear Matrix", fmt.Sprintf("A%d", fmRow), c.AttributeVal)
			f.SetCellValue("Footwear Matrix", fmt.Sprintf("B%d", fmRow), c.UnitsSold)
			f.SetCellValue("Footwear Matrix", fmt.Sprintf("C%d", fmRow), c.Revenue.InexactFloat64())
			f.SetCellValue("Footwear Matrix", fmt.Sprintf("D%d", fmRow), c.Percentage.InexactFloat64())
			fmRow++
		}
	}
	f.SetColWidth("Footwear Matrix", "A", "D", 20)

	// 7. PAYMENT SUMMARY SHEET
	f.NewSheet("Payment Summary")
	f.SetCellValue("Payment Summary", "A1", "Payment Mode")
	f.SetCellValue("Payment Summary", "B1", "Bill Count")
	f.SetCellValue("Payment Summary", "C1", "Total Collected (₹)")
	f.SetCellValue("Payment Summary", "D1", "Share %")
	f.SetCellStyle("Payment Summary", "A1", "D1", headerStyle)

	payRow := 2
	for _, pm := range paymentMethods {
		f.SetCellValue("Payment Summary", fmt.Sprintf("A%d", payRow), pm.Method)
		f.SetCellValue("Payment Summary", fmt.Sprintf("B%d", payRow), pm.BillCount)
		f.SetCellValue("Payment Summary", fmt.Sprintf("C%d", payRow), pm.TotalSales.InexactFloat64())
		f.SetCellValue("Payment Summary", fmt.Sprintf("D%d", payRow), pm.Percentage.InexactFloat64())
		payRow++
	}
	f.SetColWidth("Payment Summary", "A", "D", 20)

	// 8. TRANSACTION DETAILS SHEET
	f.NewSheet("Transaction Details")
	tHeaders := []string{"Bill Number", "Source Bill ID", "Sale Date", "Payment Mode", "Status", "Items", "Subtotal (₹)", "Discount (₹)", "Tax (₹)", "Grand Total (₹)"}
	for c, h := range tHeaders {
		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		f.SetCellValue("Transaction Details", cell, h)
	}
	f.SetCellStyle("Transaction Details", "A1", "J1", headerStyle)

	txRow := 2
	for _, tx := range rawTransactions {
		f.SetCellValue("Transaction Details", fmt.Sprintf("A%d", txRow), tx.BillNumber)
		f.SetCellValue("Transaction Details", fmt.Sprintf("B%d", txRow), tx.SourceBillID)
		f.SetCellValue("Transaction Details", fmt.Sprintf("C%d", txRow), tx.SaleDate.Format("2006-01-02 15:04"))
		f.SetCellValue("Transaction Details", fmt.Sprintf("D%d", txRow), tx.PaymentMethod)
		f.SetCellValue("Transaction Details", fmt.Sprintf("E%d", txRow), string(tx.Status))
		f.SetCellValue("Transaction Details", fmt.Sprintf("F%d", txRow), len(tx.Items))
		f.SetCellValue("Transaction Details", fmt.Sprintf("G%d", txRow), tx.Subtotal.InexactFloat64())
		f.SetCellValue("Transaction Details", fmt.Sprintf("H%d", txRow), tx.Discount.InexactFloat64())
		f.SetCellValue("Transaction Details", fmt.Sprintf("I%d", txRow), tx.TaxAmount.InexactFloat64())
		f.SetCellValue("Transaction Details", fmt.Sprintf("J%d", txRow), tx.GrandTotal.InexactFloat64())
		txRow++
	}
	f.SetColWidth("Transaction Details", "A", "J", 18)

	// Write buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("KadaiIQ_Sales_Analytics_%s_%s.xlsx", summary.PeriodType, time.Now().Format("20060102_150405"))
	return buf.Bytes(), filename, nil
}
