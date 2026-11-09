package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/Abinesh09/KadaiIQ/backend/internal/dto"
	"github.com/Abinesh09/KadaiIQ/backend/internal/models"
	"github.com/Abinesh09/KadaiIQ/backend/internal/repositories"
	"github.com/google/uuid"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) ListProducts(
	tenantID uuid.UUID,
	search string,
	brandID *uuid.UUID,
	categoryID *uuid.UUID,
	gender string,
	status string,
	page int,
	pageSize int,
) (*dto.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	prods, total, err := s.productRepo.FindProducts(
		tenantID,
		search,
		brandID,
		categoryID,
		gender,
		status,
		page,
		pageSize,
	)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ProductResponse, 0, len(prods))
	for _, p := range prods {
		items = append(items, s.formatProductResponse(&p))
	}

	return &dto.ProductListResponse{
		Items:      items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *ProductService) GetProduct(tenantID, productID uuid.UUID) (*dto.ProductResponse, error) {
	prod, err := s.productRepo.GetProductByID(tenantID, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	res := s.formatProductResponse(prod)
	return &res, nil
}

func (s *ProductService) CreateProduct(
	tenantID uuid.UUID,
	userID uuid.UUID,
	req dto.CreateProductRequest,
) (*dto.ProductResponse, error) {
	if req.Name == "" || req.ProductCode == "" {
		return nil, errors.New("product name and product code are required")
	}
	if len(req.Variants) == 0 {
		return nil, errors.New("at least one product variant must be provided")
	}

	now := time.Now().UTC()
	product := models.Product{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        req.Name,
		ProductCode: req.ProductCode,
		SKU:         req.SKU,
		Barcode:     req.Barcode,
		BrandID:     req.BrandID,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Gender:      req.Gender,
		Status:      models.ProductStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if product.Gender == "" {
		product.Gender = "Unisex"
	}

	variants := make([]models.ProductVariant, len(req.Variants))
	openingQuantities := make([]int, len(req.Variants))

	for i, v := range req.Variants {
		if v.Size == "" || v.Colour == "" || v.Barcode == "" || v.SKU == "" {
			return nil, fmt.Errorf("variant %d is missing required fields (size, colour, sku, barcode)", i+1)
		}
		if v.SellingPrice.IsNegative() || v.PurchasePrice.IsNegative() {
			return nil, fmt.Errorf("variant %d prices cannot be negative", i+1)
		}

		variants[i] = models.ProductVariant{
			ID:                uuid.New(),
			TenantID:          tenantID,
			Size:              v.Size,
			Colour:            v.Colour,
			SKU:               v.SKU,
			Barcode:           v.Barcode,
			PurchasePrice:     v.PurchasePrice,
			SellingPrice:      v.SellingPrice,
			GSTRate:           v.GSTRate,
			MinStockThreshold: v.MinStockThreshold,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		openingQuantities[i] = v.OpeningStock
	}

	if err := s.productRepo.CreateProductWithVariantsTx(&product, variants, openingQuantities, userID); err != nil {
		return nil, err
	}

	// Fetch refreshed record
	created, err := s.productRepo.GetProductByID(tenantID, product.ID)
	if err != nil {
		return nil, err
	}

	res := s.formatProductResponse(created)
	return &res, nil
}

func (s *ProductService) LookupBarcode(tenantID uuid.UUID, barcode string) (*dto.ProductVariantResponse, error) {
	variant, err := s.productRepo.FindVariantByBarcode(tenantID, barcode)
	if err != nil {
		return nil, fmt.Errorf("no product variant found with barcode '%s'", barcode)
	}

	status := "IN_STOCK"
	if variant.CurrentStock <= 0 {
		status = "OUT_OF_STOCK"
	} else if variant.CurrentStock <= variant.MinStockThreshold {
		status = "LOW_STOCK"
	}

	return &dto.ProductVariantResponse{
		ID:                variant.ID,
		ProductID:         variant.ProductID,
		Size:              variant.Size,
		Colour:            variant.Colour,
		SKU:               variant.SKU,
		Barcode:           variant.Barcode,
		PurchasePrice:     variant.PurchasePrice,
		SellingPrice:      variant.SellingPrice,
		GSTRate:           variant.GSTRate,
		CurrentStock:      variant.CurrentStock,
		MinStockThreshold: variant.MinStockThreshold,
		StockStatus:       status,
	}, nil
}

func (s *ProductService) ArchiveProduct(tenantID, productID uuid.UUID) error {
	return s.productRepo.ArchiveProduct(tenantID, productID)
}

func (s *ProductService) GetBrands(tenantID uuid.UUID) ([]models.Brand, error) {
	return s.productRepo.GetBrands(tenantID)
}

func (s *ProductService) GetCategories(tenantID uuid.UUID) ([]models.Category, error) {
	return s.productRepo.GetCategories(tenantID)
}

func (s *ProductService) RecordStockAdjustment(
	tenantID uuid.UUID,
	userID uuid.UUID,
	req dto.StockAdjustmentRequest,
) (*dto.StockMovementResponse, error) {
	movType := models.MovementType(req.MovementType)
	movement, err := s.productRepo.RecordStockMovementTx(
		tenantID,
		req.VariantID,
		movType,
		req.Quantity,
		req.UnitCost,
		req.ReferenceNumber,
		req.Notes,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return &dto.StockMovementResponse{
		ID:              movement.ID,
		VariantID:       movement.VariantID,
		MovementType:    string(movement.MovementType),
		Quantity:        movement.Quantity,
		PreviousStock:   movement.PreviousStock,
		NewStock:        movement.NewStock,
		UnitCost:        movement.UnitCost,
		ReferenceNumber: movement.ReferenceNumber,
		Notes:           movement.Notes,
		CreatedAt:       movement.CreatedAt,
	}, nil
}

func (s *ProductService) GetStockMovements(
	tenantID uuid.UUID,
	variantID *uuid.UUID,
	page int,
	pageSize int,
) ([]dto.StockMovementResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	movements, total, err := s.productRepo.GetStockMovements(tenantID, variantID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	items := make([]dto.StockMovementResponse, 0, len(movements))
	for _, m := range movements {
		prodName := "Product"
		details := ""
		barcode := ""
		if m.Variant != nil {
			barcode = m.Variant.Barcode
			details = fmt.Sprintf("Size: %s | Colour: %s", m.Variant.Size, m.Variant.Colour)
			if m.Variant.Product != nil {
				prodName = m.Variant.Product.Name
			}
		}
		userName := "System"
		if m.User != nil {
			userName = m.User.Name
		}

		items = append(items, dto.StockMovementResponse{
			ID:              m.ID,
			VariantID:       m.VariantID,
			ProductName:     prodName,
			VariantDetails:  details,
			Barcode:         barcode,
			MovementType:    string(m.MovementType),
			Quantity:        m.Quantity,
			PreviousStock:   m.PreviousStock,
			NewStock:        m.NewStock,
			UnitCost:        m.UnitCost,
			ReferenceNumber: m.ReferenceNumber,
			Notes:           m.Notes,
			CreatedByName:   userName,
			CreatedAt:       m.CreatedAt,
		})
	}

	return items, total, nil
}

func (s *ProductService) GetInventorySummary(tenantID uuid.UUID) (*dto.InventorySummaryResponse, error) {
	summary, err := s.productRepo.GetInventorySummary(tenantID)
	if err != nil {
		return nil, err
	}
	return &dto.InventorySummaryResponse{
		TotalProducts:   summary.TotalProducts,
		TotalVariants:   summary.TotalVariants,
		TotalStockUnits: summary.TotalStockUnits,
		LowStockCount:   summary.LowStockCount,
		OutOfStockCount: summary.OutOfStockCount,
		StockValuation:  summary.StockValuation,
	}, nil
}

func (s *ProductService) formatProductResponse(p *models.Product) dto.ProductResponse {
	brandName := ""
	if p.Brand != nil {
		brandName = p.Brand.Name
	}
	categoryName := ""
	if p.Category != nil {
		categoryName = p.Category.Name
	}

	totalStock := 0
	variants := make([]dto.ProductVariantResponse, 0, len(p.Variants))
	for _, v := range p.Variants {
		totalStock += v.CurrentStock
		status := "IN_STOCK"
		if v.CurrentStock <= 0 {
			status = "OUT_OF_STOCK"
		} else if v.CurrentStock <= v.MinStockThreshold {
			status = "LOW_STOCK"
		}

		variants = append(variants, dto.ProductVariantResponse{
			ID:                v.ID,
			ProductID:         v.ProductID,
			Size:              v.Size,
			Colour:            v.Colour,
			SKU:               v.SKU,
			Barcode:           v.Barcode,
			PurchasePrice:     v.PurchasePrice,
			SellingPrice:      v.SellingPrice,
			GSTRate:           v.GSTRate,
			CurrentStock:      v.CurrentStock,
			MinStockThreshold: v.MinStockThreshold,
			StockStatus:       status,
		})
	}

	return dto.ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		ProductCode:  p.ProductCode,
		SKU:          p.SKU,
		Barcode:      p.Barcode,
		BrandName:    brandName,
		CategoryName: categoryName,
		BrandID:      p.BrandID,
		CategoryID:   p.CategoryID,
		Description:  p.Description,
		Gender:       p.Gender,
		Status:       string(p.Status),
		TotalStock:   totalStock,
		Variants:     variants,
		CreatedAt:    p.CreatedAt,
	}
}
