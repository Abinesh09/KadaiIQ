package services

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var gstinRegex = regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$`)

var stateCodeMap = map[string]string{
	"01": "Jammu & Kashmir",
	"02": "Himachal Pradesh",
	"03": "Punjab",
	"04": "Chandigarh",
	"06": "Haryana",
	"07": "Delhi",
	"08": "Rajasthan",
	"09": "Uttar Pradesh",
	"10": "Bihar",
	"18": "Assam",
	"19": "West Bengal",
	"21": "Odisha",
	"22": "Chhattisgarh",
	"23": "Madhya Pradesh",
	"24": "Gujarat",
	"27": "Maharashtra",
	"29": "Karnataka",
	"30": "Goa",
	"32": "Kerala",
	"33": "Tamil Nadu",
	"36": "Telangana",
	"37": "Andhra Pradesh",
}

type GSTService struct {
	repo *repositories.GSTRepository
}

func NewGSTService(repo *repositories.GSTRepository) *GSTService {
	return &GSTService{repo: repo}
}

func (s *GSTService) ValidateGSTIN(gstin string) (stateCode, stateName string, err error) {
	gstin = strings.TrimSpace(strings.ToUpper(gstin))
	if len(gstin) != 15 {
		return "", "", errors.New("GSTIN must be exactly 15 characters long")
	}
	if !gstinRegex.MatchString(gstin) {
		return "", "", errors.New("invalid Indian GSTIN format (e.g. 33AAAAA0000A1Z5)")
	}

	stateCode = gstin[:2]
	name, ok := stateCodeMap[stateCode]
	if !ok {
		name = "Other State/UT"
	}
	return stateCode, name, nil
}

func (s *GSTService) GetSettings(tenantID uuid.UUID) (*dto.GSTSettingsResponse, error) {
	setting, err := s.repo.GetSettings(tenantID)
	if err != nil {
		return nil, err
	}
	return s.formatSettingsResponse(setting), nil
}

func (s *GSTService) UpdateSettings(tenantID uuid.UUID, req *dto.UpdateGSTSettingsRequest) (*dto.GSTSettingsResponse, error) {
	setting, err := s.repo.GetSettings(tenantID)
	if err != nil {
		return nil, err
	}

	if req.GSTIN != "" {
		code, name, err := s.ValidateGSTIN(req.GSTIN)
		if err != nil {
			return nil, err
		}
		setting.GSTIN = strings.ToUpper(req.GSTIN)
		setting.StateCode = code
		setting.StateName = name
	}

	if req.LegalName != "" {
		setting.LegalName = req.LegalName
	}
	if req.TradeName != "" {
		setting.TradeName = req.TradeName
	}
	if req.RegistrationType != "" {
		setting.RegistrationType = models.GSTRegistrationType(req.RegistrationType)
	}
	if req.DefaultFootwearHSN != "" {
		setting.DefaultFootwearHSN = req.DefaultFootwearHSN
	}
	if req.DefaultGSTRate.IsPositive() {
		setting.DefaultGSTRate = req.DefaultGSTRate
	}
	setting.EnableComposition = req.EnableComposition

	if err := s.repo.SaveSettings(setting); err != nil {
		return nil, err
	}

	return s.formatSettingsResponse(setting), nil
}

func (s *GSTService) ResolveMonthPeriod(monthStr, yearStr string) (start, end time.Time, periodKey string) {
	now := time.Now().UTC()
	year, _ := strconv.Atoi(yearStr)
	if year < 2020 || year > 2035 {
		year = now.Year()
	}

	month, _ := strconv.Atoi(monthStr)
	if month < 1 || month > 12 {
		month = int(now.Month())
	}

	start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	periodKey = fmt.Sprintf("%02d/%d", month, year)
	return
}

func (s *GSTService) GetGSTR1Summary(tenantID uuid.UUID, monthStr, yearStr string) (*dto.GSTR1SummaryResponse, error) {
	start, end, periodKey := s.ResolveMonthPeriod(monthStr, yearStr)

	setting, err := s.repo.GetSettings(tenantID)
	if err != nil {
		return nil, err
	}

	slabs, err := s.repo.GetTaxSlabs(tenantID, start, end)
	if err != nil {
		return nil, err
	}

	two := decimal.NewFromInt(2)
	totalTaxable := decimal.Zero
	totalCGST := decimal.Zero
	totalSGST := decimal.Zero
	totalTax := decimal.Zero
	totalInvoices := int64(0)

	rateBreakdown := make([]dto.GSTR1RateSlab, 0, len(slabs))
	for _, sl := range slabs {
		// In Indian GST for local retail supplies, tax is evenly split between CGST and SGST
		cgst := sl.TaxAmount.Div(two).Round(2)
		sgst := sl.TaxAmount.Sub(cgst).Round(2) // Guarantees exact sum without rounding loss

		totalTaxable = totalTaxable.Add(sl.TaxableValue)
		totalCGST = totalCGST.Add(cgst)
		totalSGST = totalSGST.Add(sgst)
		totalTax = totalTax.Add(sl.TaxAmount)
		totalInvoices += sl.Invoices

		rateBreakdown = append(rateBreakdown, dto.GSTR1RateSlab{
			RatePct:      sl.TaxRate.Round(2),
			TaxableValue: sl.TaxableValue.Round(2),
			CGST:         cgst,
			SGST:         sgst,
			IGST:         decimal.Zero,
			TotalTax:     sl.TaxAmount.Round(2),
			Invoices:     sl.Invoices,
		})
	}

	// Fetch HSN summary
	hsnRows, err := s.repo.GetHSNAggregates(tenantID, start, end, setting.DefaultFootwearHSN)
	if err != nil {
		return nil, err
	}

	hsnSummary := make([]dto.GSTR1HSNRow, 0, len(hsnRows))
	for _, h := range hsnRows {
		if h.TotalQty > 0 || h.TaxableValue.IsPositive() {
			cgst := h.TaxAmount.Div(two).Round(2)
			sgst := h.TaxAmount.Sub(cgst).Round(2)

			hsnSummary = append(hsnSummary, dto.GSTR1HSNRow{
				HSNCode:      h.HSNCode,
				Description:  h.Description,
				UQC:          "PRS", // Pairs (Standard footwear unit)
				TotalQty:     h.TotalQty,
				TaxableValue: h.TaxableValue.Round(2),
				CGST:         cgst,
				SGST:         sgst,
				IGST:         decimal.Zero,
				TotalTax:     h.TaxAmount.Round(2),
			})
		}
	}

	return &dto.GSTR1SummaryResponse{
		ReturnPeriod:      periodKey,
		StartDate:         start,
		EndDate:           end,
		GSTIN:             setting.GSTIN,
		LegalName:         setting.LegalName,
		StateCode:         setting.StateCode,
		TotalTaxableValue: totalTaxable.Round(2),
		TotalCGST:         totalCGST.Round(2),
		TotalSGST:         totalSGST.Round(2),
		TotalIGST:         decimal.Zero,
		TotalTaxLiability: totalTax.Round(2),
		B2CSmallCount:     totalInvoices,
		RateBreakdown:     rateBreakdown,
		HSNSummary:        hsnSummary,
	}, nil
}

func (s *GSTService) GetGSTR3BTable31(tenantID uuid.UUID, monthStr, yearStr string) (*dto.GSTR3BTable31Response, error) {
	gstr1, err := s.GetGSTR1Summary(tenantID, monthStr, yearStr)
	if err != nil {
		return nil, err
	}

	return &dto.GSTR3BTable31Response{
		NatureOfSupplies: "(a) Outward taxable supplies (other than zero rated, nil rated and exempted)",
		TaxableValue:     gstr1.TotalTaxableValue,
		IntegratedTax:    gstr1.TotalIGST,
		CentralTax:       gstr1.TotalCGST,
		StateUTTax:       gstr1.TotalSGST,
		Cess:             decimal.Zero,
	}, nil
}

func (s *GSTService) formatSettingsResponse(st *models.TenantGSTSetting) *dto.GSTSettingsResponse {
	return &dto.GSTSettingsResponse{
		ID:                 st.ID,
		TenantID:           st.TenantID,
		GSTIN:              st.GSTIN,
		LegalName:          st.LegalName,
		TradeName:          st.TradeName,
		StateCode:          st.StateCode,
		StateName:          st.StateName,
		RegistrationType:   string(st.RegistrationType),
		DefaultFootwearHSN: st.DefaultFootwearHSN,
		DefaultGSTRate:     st.DefaultGSTRate,
		EnableComposition:  st.EnableComposition,
		UpdatedAt:          st.UpdatedAt,
	}
}
