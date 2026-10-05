package models

import "time"

// CostItem represents a raw material or single input item
type CostItem struct {
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	PricePerUnit float64 `json:"price_per_unit"`
	Subtotal     float64 `json:"subtotal"`
}

// LaborItem represents a direct labor cost
type LaborItem struct {
	Role        string  `json:"role"`
	WorkerCount float64 `json:"worker_count"`
	WageRate    float64 `json:"wage_rate"`
	HoursOrDays float64 `json:"hours_or_days"`
	Subtotal    float64 `json:"subtotal"`
}

// OverheadItem represents overhead expenses (BOP)
type OverheadItem struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Cost     float64 `json:"cost"`
}

// ManufactureHppRequest payload for manufacturing / culinary / craft HPP
type ManufactureHppRequest struct {
	ProductName     string         `json:"product_name"`
	ProductionUnits float64        `json:"production_units"`
	MarkupPercent   float64        `json:"markup_percent"`
	PricingMethod   string         `json:"pricing_method"` // "markup_on_cost" or "margin_on_sales"
	RawMaterials    []CostItem     `json:"raw_materials"`
	LaborCosts      []LaborItem    `json:"labor_costs"`
	OverheadCosts   []OverheadItem `json:"overhead_costs"`
}

// ManufactureHppResult response after calculating manufacturing HPP
type ManufactureHppResult struct {
	ProductName           string         `json:"product_name"`
	ProductionUnits       float64        `json:"production_units"`
	RawMaterials          []CostItem     `json:"raw_materials"`
	LaborCosts            []LaborItem    `json:"labor_costs"`
	OverheadCosts         []OverheadItem `json:"overhead_costs"`
	TotalRawMaterialCost  float64        `json:"total_raw_material_cost"`
	TotalLaborCost        float64        `json:"total_labor_cost"`
	TotalOverheadCost     float64        `json:"total_overhead_cost"`
	TotalProductionCost   float64        `json:"total_production_cost"`
	HppPerUnit            float64        `json:"hpp_per_unit"`
	PricingMethod         string         `json:"pricing_method"`
	MarkupPercent         float64        `json:"markup_percent"`
	SellingPricePerUnit   float64        `json:"selling_price_per_unit"`
	ProfitPerUnit         float64        `json:"profit_per_unit"`
	TotalRevenue          float64        `json:"total_revenue"`
	TotalProfit           float64        `json:"total_profit"`
	ProfitMarginPercent   float64        `json:"profit_margin_percent"`
	RawMaterialPercentage float64        `json:"raw_material_percentage"`
	LaborPercentage       float64        `json:"labor_percentage"`
	OverheadPercentage    float64        `json:"overhead_percentage"`
	BreakEvenUnits        float64        `json:"break_even_units"`
	BreakEvenRevenue      float64        `json:"break_even_revenue"`
	CalculatedAt          string         `json:"calculated_at"`
}

// RetailHppRequest payload for retail / trade business HPP
type RetailHppRequest struct {
	PeriodName         string  `json:"period_name"`
	BeginningInventory float64 `json:"beginning_inventory"`
	Purchases          float64 `json:"purchases"`
	FreightIn          float64 `json:"freight_in"`
	PurchaseReturns    float64 `json:"purchase_returns"`
	PurchaseDiscounts  float64 `json:"purchase_discounts"`
	EndingInventory    float64 `json:"ending_inventory"`
	TotalSales         float64 `json:"total_sales"`
	SalesReturns       float64 `json:"sales_returns"`
	OperatingExpenses  float64 `json:"operating_expenses"`
}

// RetailHppResult response after calculating retail HPP
type RetailHppResult struct {
	PeriodName               string  `json:"period_name"`
	BeginningInventory       float64 `json:"beginning_inventory"`
	GrossPurchases           float64 `json:"gross_purchases"`
	FreightIn                float64 `json:"freight_in"`
	PurchaseDeductions       float64 `json:"purchase_deductions"`
	NetPurchases             float64 `json:"net_purchases"`
	GoodsAvailableForSale    float64 `json:"goods_available_for_sale"`
	EndingInventory          float64 `json:"ending_inventory"`
	TotalHpp                 float64 `json:"total_hpp"`
	GrossSales               float64 `json:"gross_sales"`
	SalesReturns             float64 `json:"sales_returns"`
	NetSales                 float64 `json:"net_sales"`
	GrossProfit              float64 `json:"gross_profit"`
	GrossProfitMarginPercent float64 `json:"gross_profit_margin_percent"`
	OperatingExpenses        float64 `json:"operating_expenses"`
	NetProfit                float64 `json:"net_profit"`
	NetProfitMarginPercent   float64 `json:"net_profit_margin_percent"`
	CalculatedAt             string  `json:"calculated_at"`
}

// SamplePreset holds preset sample data
type SamplePreset struct {
	ID              string                  `json:"id"`
	Title           string                  `json:"title"`
	Category        string                  `json:"category"`
	Type            string                  `json:"type"` // "manufacture" or "retail"
	Description     string                  `json:"description"`
	ManufactureData *ManufactureHppRequest `json:"manufacture_data,omitempty"`
	RetailData      *RetailHppRequest      `json:"retail_data,omitempty"`
}

// CurrentTimestamp returns formatted current timestamp
func CurrentTimestamp() string {
	return time.Now().Format("02-01-2006 15:04:05")
}
