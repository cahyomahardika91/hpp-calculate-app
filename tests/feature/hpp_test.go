package feature

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/app/models"
	"goravel/app/services"
	"goravel/tests"
)

type HppTestSuite struct {
	suite.Suite
	tests.TestCase
	service *services.HppService
}

func TestHppTestSuite(t *testing.T) {
	suite.Run(t, new(HppTestSuite))
}

func (s *HppTestSuite) SetupTest() {
	s.service = services.NewHppService()
}

// Test HppService Manufacture Calculation with Markup on Cost
func (s *HppTestSuite) TestServiceCalculateManufacture_MarkupOnCost() {
	req := models.ManufactureHppRequest{
		ProductName:     "Kopi Susu Aren",
		ProductionUnits: 100,
		MarkupPercent:   50, // 50%
		PricingMethod:   "markup_on_cost",
		RawMaterials: []models.CostItem{
			{Name: "Kopi", Quantity: 1, PricePerUnit: 100000},
			{Name: "Susu", Quantity: 5, PricePerUnit: 20000},
		},
		LaborCosts: []models.LaborItem{
			{Role: "Barista", WorkerCount: 1, WageRate: 50000, HoursOrDays: 1},
		},
		OverheadCosts: []models.OverheadItem{
			{Name: "Listrik", Category: "Utilitas", Cost: 50000},
		},
	}

	result, err := s.service.CalculateManufacture(req)
	s.NoError(err)
	s.NotNil(result)

	// Total Raw Material = 100,000 + 100,000 = 200,000
	s.Equal(200000.0, result.TotalRawMaterialCost)
	// Total Labor = 50,000
	s.Equal(50000.0, result.TotalLaborCost)
	// Total Overhead = 50,000
	s.Equal(50000.0, result.TotalOverheadCost)
	// Total Production Cost = 300,000
	s.Equal(300000.0, result.TotalProductionCost)
	// HPP per Unit = 300,000 / 100 = 3,000
	s.Equal(3000.0, result.HppPerUnit)
	// Markup on Cost 50%: Selling price = 3,000 * 1.50 = 4,500
	s.Equal(4500.0, result.SellingPricePerUnit)
	// Profit per unit = 1,500
	s.Equal(1500.0, result.ProfitPerUnit)
	// Total Revenue = 450,000
	s.Equal(450000.0, result.TotalRevenue)
	// Total Profit = 150,000
	s.Equal(150000.0, result.TotalProfit)
	// Cost percentages: RM = 66.67%, Labor = 16.67%, Overhead = 16.67%
	s.InDelta(66.67, result.RawMaterialPercentage, 0.1)
	s.InDelta(16.67, result.LaborPercentage, 0.1)
	s.InDelta(16.67, result.OverheadPercentage, 0.1)
}

// Test HppService Manufacture Calculation with Margin on Sales
func (s *HppTestSuite) TestServiceCalculateManufacture_MarginOnSales() {
	req := models.ManufactureHppRequest{
		ProductName:     "Roti Manis",
		ProductionUnits: 10,
		MarkupPercent:   20, // 20% margin on sales
		PricingMethod:   "margin_on_sales",
		RawMaterials: []models.CostItem{
			{Name: "Tepung", Quantity: 2, PricePerUnit: 10000}, // 20,000
			{Name: "Gula", Quantity: 1, PricePerUnit: 12000},   // 12,000
		},
		LaborCosts: []models.LaborItem{
			{Role: "Baker", WorkerCount: 1, WageRate: 8000, HoursOrDays: 1}, // 8,000
		},
	}

	result, err := s.service.CalculateManufacture(req)
	s.NoError(err)
	s.NotNil(result)

	// Total Production Cost = 20,000 + 12,000 + 8,000 = 40,000
	s.Equal(40000.0, result.TotalProductionCost)
	// HPP per Unit = 4,000
	s.Equal(4000.0, result.HppPerUnit)
	// Selling Price = 4,000 / (1 - 0.20) = 5,000
	s.Equal(5000.0, result.SellingPricePerUnit)
	s.Equal(1000.0, result.ProfitPerUnit)
	// Profit Margin = 1,000 / 5,000 = 20%
	s.Equal(20.0, result.ProfitMarginPercent)
}

// Test HppService Manufacture validation on zero units
func (s *HppTestSuite) TestServiceCalculateManufacture_ZeroUnits() {
	req := models.ManufactureHppRequest{
		ProductionUnits: 0,
	}

	result, err := s.service.CalculateManufacture(req)
	s.Error(err)
	s.Nil(result)
}

// Test HppService Retail Calculation
func (s *HppTestSuite) TestServiceCalculateRetail() {
	req := models.RetailHppRequest{
		PeriodName:         "Oktober 2026",
		BeginningInventory: 10000000,
		Purchases:          20000000,
		FreightIn:          500000,
		PurchaseReturns:    300000,
		PurchaseDiscounts:  200000,
		EndingInventory:    8000000,
		TotalSales:         30000000,
		SalesReturns:       500000,
		OperatingExpenses:  3000000,
	}

	result, err := s.service.CalculateRetail(req)
	s.NoError(err)
	s.NotNil(result)

	// Purchase Deductions = 300,000 + 200,000 = 500,000
	s.Equal(500000.0, result.PurchaseDeductions)
	// Net Purchases = 20,000,000 + 500,000 - 500,000 = 20,000,000
	s.Equal(20000000.0, result.NetPurchases)
	// Goods Available for Sale = 10,000,000 + 20,000,000 = 30,000,000
	s.Equal(30000000.0, result.GoodsAvailableForSale)
	// Total HPP = 30,000,000 - 8,000,000 = 22,000,000
	s.Equal(22000000.0, result.TotalHpp)
	// Net Sales = 30,000,000 - 500,000 = 29,500,000
	s.Equal(29500000.0, result.NetSales)
	// Gross Profit = 29,500,000 - 22,000,000 = 7,500,000
	s.Equal(7500000.0, result.GrossProfit)
	// Net Profit = 7,500,000 - 3,000,000 = 4,500,000
	s.Equal(4500000.0, result.NetProfit)
}

// Test HppService Presets
func (s *HppTestSuite) TestServiceGetPresets() {
	presets := s.service.GetPresets()
	s.NotEmpty(presets)
	s.GreaterOrEqual(len(presets), 3)

	foundManufacture := false
	foundRetail := false
	for _, p := range presets {
		if p.Type == "manufacture" {
			foundManufacture = true
		}
		if p.Type == "retail" {
			foundRetail = true
		}
	}
	s.True(foundManufacture)
	s.True(foundRetail)
}

// Test HTTP API: GET /api/hpp/presets
func (s *HppTestSuite) TestApiPresets() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/hpp/presets", nil)
	s.NoError(err)

	facades.Route().ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("success", resp["status"])
	s.NotNil(resp["data"])
}

// Test HTTP API: POST /api/hpp/calculate-manufacture
func (s *HppTestSuite) TestApiCalculateManufacture() {
	payload := models.ManufactureHppRequest{
		ProductName:     "Baju Kaos",
		ProductionUnits: 50,
		MarkupPercent:   40,
		PricingMethod:   "markup_on_cost",
		RawMaterials: []models.CostItem{
			{Name: "Kain", Quantity: 10, PricePerUnit: 100000, Subtotal: 1000000},
		},
		LaborCosts: []models.LaborItem{
			{Role: "Jahit", WageRate: 200000, HoursOrDays: 1, WorkerCount: 1, Subtotal: 200000},
		},
		OverheadCosts: []models.OverheadItem{
			{Name: "Listrik", Cost: 50000},
		},
	}

	body, err := json.Marshal(payload)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/hpp/calculate-manufacture", bytes.NewReader(body))
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	facades.Route().ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var resp struct {
		Status  string                      `json:"status"`
		Message string                      `json:"message"`
		Data    models.ManufactureHppResult `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("success", resp.Status)
	s.Equal("Baju Kaos", resp.Data.ProductName)
	s.Equal(1250000.0, resp.Data.TotalProductionCost)
	s.Equal(25000.0, resp.Data.HppPerUnit)
}

// Test HTTP API: POST /api/hpp/calculate-retail
func (s *HppTestSuite) TestApiCalculateRetail() {
	payload := models.RetailHppRequest{
		PeriodName:         "Toko Baru",
		BeginningInventory: 5000000,
		Purchases:          10000000,
		EndingInventory:    4000000,
		TotalSales:         15000000,
		OperatingExpenses:  2000000,
	}

	body, err := json.Marshal(payload)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/hpp/calculate-retail", bytes.NewReader(body))
	s.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	facades.Route().ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var resp struct {
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Data    models.RetailHppResult `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("success", resp.Status)
	// Total HPP = 5,000,000 + 10,000,000 - 4,000,000 = 11,000,000
	s.Equal(11000000.0, resp.Data.TotalHpp)
	// Gross Profit = 15,000,000 - 11,000,000 = 4,000,000
	s.Equal(4000000.0, resp.Data.GrossProfit)
	// Net Profit = 4,000,000 - 2,000,000 = 2,000,000
	s.Equal(2000000.0, resp.Data.NetProfit)
}

// Test Web Page View rendering
func (s *HppTestSuite) TestWebHppPage() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	s.NoError(err)

	facades.Route().ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Kalkulator HPP")
}
