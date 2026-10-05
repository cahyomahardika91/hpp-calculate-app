package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/support"

	"goravel/app/facades"
	"goravel/app/http/controllers"
)

func Web() {
	hppController := controllers.NewHppController()
	userController := controllers.NewUserController()

	// Static Assets
	facades.Route().Static("public", "./public")

	// Main HPP Calculator Application Web Pages
	facades.Route().Get("/", hppController.Index)
	facades.Route().Get("/kalkulator-hpp", hppController.Index)

	// Goravel default welcome page preserved at /welcome
	facades.Route().Get("/welcome", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	// User sample endpoint
	facades.Route().Get("/users", userController.Index)

	// HPP REST API Endpoints
	facades.Route().Prefix("api/hpp").Group(func(router route.Router) {
		router.Get("/presets", hppController.Presets)
		router.Post("/calculate-manufacture", hppController.CalculateManufacture)
		router.Post("/calculate-retail", hppController.CalculateRetail)
	})
}
