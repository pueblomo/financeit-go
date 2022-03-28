package api

import (
	"github.com/gofiber/fiber/v2"
	"main.go/categories"
	"main.go/expenses"
)

func InitRoutes(app *fiber.App) {
	catGroup := app.Group("/categories")
	expGroup := app.Group("/expenses")

	initCatRoutes(catGroup)
	initExpRoutes(expGroup)
}

func initCatRoutes(grp fiber.Router) {
	grp.Get("/", func (c *fiber.Ctx) error {
		return categories.GetCategories(c)
	})

	grp.Post("/", func (c * fiber.Ctx) error {
		return categories.PostCategorie(c)
	})
}

func initExpRoutes(grp fiber.Router) {
	grp.Post("/", func (c * fiber.Ctx) error {
		return expenses.PostExpense(c)
	})
}