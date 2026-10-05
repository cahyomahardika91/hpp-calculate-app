package services

import (
	"errors"
	"math"

	"goravel/app/models"
)

type HppService struct{}

func NewHppService() *HppService {
	return &HppService{}
}

// roundTo rounds a float to the specified decimal places
func roundTo(val float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(val*pow) / pow
}

// CalculateManufacture performs HPP calculation for production/manufacturing businesses
func (s *HppService) CalculateManufacture(req models.ManufactureHppRequest) (*models.ManufactureHppResult, error) {
	if req.ProductionUnits <= 0 {
		return nil, errors.New("jumlah unit produksi harus lebih besar dari 0")
	}

	if req.PricingMethod == "" {
		req.PricingMethod = "markup_on_cost"
	}

	// 1. Calculate Raw Materials
	totalRawMaterial := 0.0
	for i := range req.RawMaterials {
		item := &req.RawMaterials[i]
		if item.Subtotal <= 0 && item.Quantity > 0 && item.PricePerUnit > 0 {
			item.Subtotal = item.Quantity * item.PricePerUnit
		}
		item.Subtotal = roundTo(item.Subtotal, 2)
		totalRawMaterial += item.Subtotal
	}
	totalRawMaterial = roundTo(totalRawMaterial, 2)

	// 2. Calculate Labor Costs
	totalLabor := 0.0
	for i := range req.LaborCosts {
		item := &req.LaborCosts[i]
		if item.Subtotal <= 0 {
			if item.WorkerCount > 0 && item.WageRate > 0 {
				multiplier := item.HoursOrDays
				if multiplier <= 0 {
					multiplier = 1
				}
				item.Subtotal = item.WorkerCount * item.WageRate * multiplier
			} else if item.WageRate > 0 {
				item.Subtotal = item.WageRate
			}
		}
		item.Subtotal = roundTo(item.Subtotal, 2)
		totalLabor += item.Subtotal
	}
	totalLabor = roundTo(totalLabor, 2)

	// 3. Calculate Overhead Costs
	totalOverhead := 0.0
	for i := range req.OverheadCosts {
		item := &req.OverheadCosts[i]
		item.Cost = roundTo(item.Cost, 2)
		totalOverhead += item.Cost
	}
	totalOverhead = roundTo(totalOverhead, 2)

	// Total Production Cost & HPP per Unit
	totalProductionCost := roundTo(totalRawMaterial+totalLabor+totalOverhead, 2)
	hppPerUnit := roundTo(totalProductionCost/req.ProductionUnits, 2)

	// Pricing calculation
	var sellingPricePerUnit float64
	markup := req.MarkupPercent
	if req.PricingMethod == "margin_on_sales" {
		if markup >= 100 {
			markup = 99.9 // prevent division by zero or negative
		}
		marginRatio := markup / 100.0
		sellingPricePerUnit = roundTo(hppPerUnit/(1.0-marginRatio), 2)
	} else {
		// Default: markup_on_cost
		markupRatio := markup / 100.0
		sellingPricePerUnit = roundTo(hppPerUnit*(1.0+markupRatio), 2)
	}

	profitPerUnit := roundTo(sellingPricePerUnit-hppPerUnit, 2)
	totalRevenue := roundTo(sellingPricePerUnit*req.ProductionUnits, 2)
	totalProfit := roundTo(profitPerUnit*req.ProductionUnits, 2)

	var profitMarginPercent float64
	if totalRevenue > 0 {
		profitMarginPercent = roundTo((totalProfit/totalRevenue)*100.0, 2)
	}

	// Cost breakdown percentages
	var rawMaterialPct, laborPct, overheadPct float64
	if totalProductionCost > 0 {
		rawMaterialPct = roundTo((totalRawMaterial/totalProductionCost)*100.0, 2)
		laborPct = roundTo((totalLabor/totalProductionCost)*100.0, 2)
		overheadPct = roundTo((totalOverhead/totalProductionCost)*100.0, 2)
	}

	// Break-Even Point (BEP) Analysis
	// Fixed Cost ~ Overhead; Variable Cost per Unit ~ (Raw Material + Labor) / Units
	variableCostPerUnit := roundTo((totalRawMaterial+totalLabor)/req.ProductionUnits, 2)
	contributionMargin := roundTo(sellingPricePerUnit-variableCostPerUnit, 2)

	var breakEvenUnits float64
	var breakEvenRevenue float64
	if contributionMargin > 0 && totalOverhead > 0 {
		breakEvenUnits = roundTo(totalOverhead/contributionMargin, 2)
		breakEvenRevenue = roundTo(breakEvenUnits*sellingPricePerUnit, 2)
	}

	productName := req.ProductName
	if productName == "" {
		productName = "Produk Tanpa Nama"
	}

	return &models.ManufactureHppResult{
		ProductName:           productName,
		ProductionUnits:       req.ProductionUnits,
		RawMaterials:          req.RawMaterials,
		LaborCosts:            req.LaborCosts,
		OverheadCosts:         req.OverheadCosts,
		TotalRawMaterialCost:  totalRawMaterial,
		TotalLaborCost:        totalLabor,
		TotalOverheadCost:     totalOverhead,
		TotalProductionCost:   totalProductionCost,
		HppPerUnit:            hppPerUnit,
		PricingMethod:         req.PricingMethod,
		MarkupPercent:         req.MarkupPercent,
		SellingPricePerUnit:   sellingPricePerUnit,
		ProfitPerUnit:         profitPerUnit,
		TotalRevenue:          totalRevenue,
		TotalProfit:           totalProfit,
		ProfitMarginPercent:   profitMarginPercent,
		RawMaterialPercentage: rawMaterialPct,
		LaborPercentage:       laborPct,
		OverheadPercentage:    overheadPct,
		BreakEvenUnits:        breakEvenUnits,
		BreakEvenRevenue:      breakEvenRevenue,
		CalculatedAt:          models.CurrentTimestamp(),
	}, nil
}

// CalculateRetail performs standard accounting HPP calculation for retail/trading
func (s *HppService) CalculateRetail(req models.RetailHppRequest) (*models.RetailHppResult, error) {
	purchaseDeductions := roundTo(req.PurchaseReturns+req.PurchaseDiscounts, 2)
	netPurchases := roundTo(req.Purchases+req.FreightIn-purchaseDeductions, 2)
	goodsAvailableForSale := roundTo(req.BeginningInventory+netPurchases, 2)
	totalHpp := roundTo(goodsAvailableForSale-req.EndingInventory, 2)

	netSales := roundTo(req.TotalSales-req.SalesReturns, 2)
	grossProfit := roundTo(netSales-totalHpp, 2)

	var grossProfitMarginPercent float64
	if netSales > 0 {
		grossProfitMarginPercent = roundTo((grossProfit/netSales)*100.0, 2)
	}

	netProfit := roundTo(grossProfit-req.OperatingExpenses, 2)
	var netProfitMarginPercent float64
	if netSales > 0 {
		netProfitMarginPercent = roundTo((netProfit/netSales)*100.0, 2)
	}

	periodName := req.PeriodName
	if periodName == "" {
		periodName = "Periode Berjalan"
	}

	return &models.RetailHppResult{
		PeriodName:               periodName,
		BeginningInventory:       req.BeginningInventory,
		GrossPurchases:           req.Purchases,
		FreightIn:                req.FreightIn,
		PurchaseDeductions:       purchaseDeductions,
		NetPurchases:             netPurchases,
		GoodsAvailableForSale:    goodsAvailableForSale,
		EndingInventory:          req.EndingInventory,
		TotalHpp:                 totalHpp,
		GrossSales:               req.TotalSales,
		SalesReturns:             req.SalesReturns,
		NetSales:                 netSales,
		GrossProfit:              grossProfit,
		GrossProfitMarginPercent: grossProfitMarginPercent,
		OperatingExpenses:        req.OperatingExpenses,
		NetProfit:                netProfit,
		NetProfitMarginPercent:   netProfitMarginPercent,
		CalculatedAt:             models.CurrentTimestamp(),
	}, nil
}

// GetPresets returns curated preset scenarios
func (s *HppService) GetPresets() []models.SamplePreset {
	return []models.SamplePreset{
		{
			ID:          "kopi-susu-aren",
			Title:       "Kopi Susu Gula Aren (F&B / Kuliner)",
			Category:    "Kuliner",
			Type:        "manufacture",
			Description: "Simulasi perhitungan HPP produksi 100 cup Kopi Susu Gula Aren untuk coffee shop / kedai.",
			ManufactureData: &models.ManufactureHppRequest{
				ProductName:     "Es Kopi Susu Gula Aren (16oz)",
				ProductionUnits: 100,
				MarkupPercent:   60,
				PricingMethod:   "markup_on_cost",
				RawMaterials: []models.CostItem{
					{Name: "Espresso Roast Coffee Beans", Quantity: 1.2, Unit: "kg", PricePerUnit: 180000, Subtotal: 216000},
					{Name: "Fresh Milk UHT", Quantity: 10, Unit: "liter", PricePerUnit: 20000, Subtotal: 200000},
					{Name: "Sirup Gula Aren Cair", Quantity: 2, Unit: "liter", PricePerUnit: 40000, Subtotal: 80000},
					{Name: "Cup Plastik 16oz + Tutup + Sedotan", Quantity: 100, Unit: "set", PricePerUnit: 650, Subtotal: 65000},
					{Name: "Es Batu Tube Higienis", Quantity: 10, Unit: "kg", PricePerUnit: 2500, Subtotal: 25000},
				},
				LaborCosts: []models.LaborItem{
					{Role: "Barista Part-time", WorkerCount: 1, WageRate: 20000, HoursOrDays: 4, Subtotal: 80000},
				},
				OverheadCosts: []models.OverheadItem{
					{Name: "Listrik Mesin Espresso & Grinder", Category: "Utilitas", Cost: 30000},
					{Name: "Plastik Takeaway & Sealer", Category: "Kemasan", Cost: 15000},
					{Name: "Alokasi Sewa Tempat & Air", Category: "Sewa", Cost: 45000},
				},
			},
		},
		{
			ID:          "kaos-sablon-distro",
			Title:       "Kaos Sablon Distro Cotton Combed 30s",
			Category:    "Konveksi / Manufaktur",
			Type:        "manufacture",
			Description: "Simulasi HPP konveksi sablon kaos distro katun combed 30s sebanyak 50 pcs per batch.",
			ManufactureData: &models.ManufactureHppRequest{
				ProductName:     "Kaos Distro Cotton Combed 30s (Black)",
				ProductionUnits: 50,
				MarkupPercent:   50,
				PricingMethod:   "markup_on_cost",
				RawMaterials: []models.CostItem{
					{Name: "Kain Cotton Combed 30s Reaktif", Quantity: 12.5, Unit: "kg", PricePerUnit: 120000, Subtotal: 1500000},
					{Name: "Kain Rib Leher", Quantity: 1, Unit: "kg", PricePerUnit: 80000, Subtotal: 80000},
					{Name: "Tinta Plastisol + Thinner + Film Sablon", Quantity: 1, Unit: "paket", PricePerUnit: 220000, Subtotal: 220000},
					{Name: "Plastik OPP Tebal + Hangtag + Label Woven", Quantity: 50, Unit: "set", PricePerUnit: 1200, Subtotal: 60000},
				},
				LaborCosts: []models.LaborItem{
					{Role: "Penjahit & Obras", WorkerCount: 1, WageRate: 7000, HoursOrDays: 50, Subtotal: 350000},
					{Role: "Tukang Sablon Manual", WorkerCount: 1, WageRate: 6000, HoursOrDays: 50, Subtotal: 300000},
				},
				OverheadCosts: []models.OverheadItem{
					{Name: "Listrik Heatpress & Mesin Jahit", Category: "Utilitas", Cost: 50000},
					{Name: "Penyusutan Alat Sablon & Jahit", Category: "Penyusutan", Cost: 40000},
					{Name: "Benang & Jarum Jahit", Category: "Bahan Pembantu", Cost: 25000},
				},
			},
		},
		{
			ID:          "fudgy-brownies",
			Title:       "Brownies Panggang Fudgy Sekat (Bakery)",
			Category:    "Bakery",
			Type:        "manufacture",
			Description: "Simulasi HPP pembuatan 20 loyang Brownies Fudgy Sekat ukuran 20x20 cm.",
			ManufactureData: &models.ManufactureHppRequest{
				ProductName:     "Fudgy Brownies Panggang Sekat 20x20",
				ProductionUnits: 20,
				MarkupPercent:   45,
				PricingMethod:   "margin_on_sales",
				RawMaterials: []models.CostItem{
					{Name: "Dark Cooking Chocolate (DCC)", Quantity: 2.5, Unit: "kg", PricePerUnit: 65000, Subtotal: 162500},
					{Name: "Butter & Margarin Blend", Quantity: 1.5, Unit: "kg", PricePerUnit: 70000, Subtotal: 105000},
					{Name: "Telur Ayam Segar", Quantity: 2, Unit: "kg", PricePerUnit: 28000, Subtotal: 56000},
					{Name: "Tepung Terigu Protein Sedang", Quantity: 2, Unit: "kg", PricePerUnit: 14000, Subtotal: 28000},
					{Name: "Gula Halus / Kastor", Quantity: 2.5, Unit: "kg", PricePerUnit: 18000, Subtotal: 45000},
					{Name: "Cokelat Bubuk Premium", Quantity: 0.8, Unit: "kg", PricePerUnit: 80000, Subtotal: 64000},
					{Name: "Topping Almond Slice, Chocochip, Keju", Quantity: 1, Unit: "paket", PricePerUnit: 75000, Subtotal: 75000},
					{Name: "Box Brownies Kraft + Sekat 25 + Pita", Quantity: 20, Unit: "set", PricePerUnit: 4500, Subtotal: 90000},
				},
				LaborCosts: []models.LaborItem{
					{Role: "Baker & Pengemasan", WorkerCount: 1, WageRate: 25000, HoursOrDays: 5, Subtotal: 125000},
				},
				OverheadCosts: []models.OverheadItem{
					{Name: "Gas Oven LPG", Category: "Bahan Bakar", Cost: 40000},
					{Name: "Listrik Mixer & Lampu", Category: "Utilitas", Cost: 20000},
					{Name: "Baking Paper & Sarung Tangan", Category: "Perlengkapan", Cost: 25000},
				},
			},
		},
		{
			ID:          "retail-sembako",
			Title:       "Toko Kelontong & Sembako (Bisnis Dagang/Retail)",
			Category:    "Retail",
			Type:        "retail",
			Description: "Perhitungan HPP usaha dagang retail sembako untuk laporan laba kotor 1 bulan operasional.",
			RetailData: &models.RetailHppRequest{
				PeriodName:         "Laporan Operasional Bulan Oktober 2026",
				BeginningInventory: 35000000,
				Purchases:          55000000,
				FreightIn:          1500000,
				PurchaseReturns:    750000,
				PurchaseDiscounts:  850000,
				EndingInventory:    32000000,
				TotalSales:         74000000,
				SalesReturns:       400000,
				OperatingExpenses:  6200000,
			},
		},
	}
}
