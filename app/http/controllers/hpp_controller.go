package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support"

	"goravel/app/models"
	"goravel/app/services"
)

type HppController struct {
	hppService *services.HppService
}

func NewHppController() *HppController {
	return &HppController{
		hppService: services.NewHppService(),
	}
}

// Index renders the main HPP calculator interactive web interface
func (c *HppController) Index(ctx http.Context) http.Response {
	presets := c.hppService.GetPresets()
	return ctx.Response().View().Make("hpp.tmpl", map[string]any{
		"version": support.Version,
		"presets": presets,
	})
}

// CalculateManufacture calculates HPP for manufacturing / food & beverage / craft products
func (c *HppController) CalculateManufacture(ctx http.Context) http.Response {
	var req models.ManufactureHppRequest
	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"status":  "error",
			"message": "Format data tidak valid: " + err.Error(),
		})
	}

	result, err := c.hppService.CalculateManufacture(req)
	if err != nil {
		return ctx.Response().Status(422).Json(http.Json{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"status":  "success",
		"message": "Perhitungan HPP Manufaktur berhasil",
		"data":    result,
	})
}

// CalculateRetail calculates HPP for retail / trade businesses
func (c *HppController) CalculateRetail(ctx http.Context) http.Response {
	var req models.RetailHppRequest
	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"status":  "error",
			"message": "Format data tidak valid: " + err.Error(),
		})
	}

	result, err := c.hppService.CalculateRetail(req)
	if err != nil {
		return ctx.Response().Status(422).Json(http.Json{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"status":  "success",
		"message": "Perhitungan HPP Dagang/Retail berhasil",
		"data":    result,
	})
}

// Presets returns sample presets for quick demonstration
func (c *HppController) Presets(ctx http.Context) http.Response {
	presets := c.hppService.GetPresets()
	return ctx.Response().Success().Json(http.Json{
		"status": "success",
		"data":   presets,
	})
}
