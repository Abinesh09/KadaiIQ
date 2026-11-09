package services

import (
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

	case "last_month":
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		start = firstOfThisMonth.AddDate(0, -1, 0)
		end = firstOfThisMonth.Add(-time.Nanosecond)
		prevStart = start.AddDate(0, -1, 0)
		prevEnd = start.Add(-time.Nanosecond)
		resolvedKey = "last_month"

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
