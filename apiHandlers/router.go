package apiHandlers

import (
	"employee-crud/api"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber on Render!")
	})

	app.Post("/CreateCategory", api.CreateCategoryApi)
	app.Get("/FindAllCategory", api.FindAllCategoriesApi)
	app.Delete("/DeleteCategory", api.DeleteCategoryApi)
	app.Post("/CreateBrands", api.CreateBrand)
	app.Get("/FindAllBrands", api.FindAllBrands)
	app.Delete("/DeleteBrand", api.DeleteBrandApi)
	app.Post("/CreateSubCategory", api.CreateSubCategory)
	app.Get("/FindAllSubCategory", api.FindAllSubCategory)
	app.Delete("/DeleteSubCategory", api.DeleteSubCategoryApi)
	app.Post("/CreateProduct", api.CreateProduct)
	app.Get("/FindAllProducts", api.FindAllProducts)
	app.Delete("/DeleteProducts", api.DeleteProductApi)
	app.Post("/CreateSupplier", api.CreateSupplier)
	app.Get("/FindAllSuppliers", api.FindAllSuppliers)
	app.Delete("/DeleteSupplierById", api.DeleteSupplierApi)
	app.Post("/AssignProductToSupplier", api.AssignProductToSupplierApi)
	app.Get("/FindProductsBySupplierID", api.GetProductsBySupplierApi)
	app.Get("/FindProductsByCategoryId", api.GetProductsByCategoryApi)
	app.Get("/FindProductsByBrandId", api.GetProductsByBrandApi)
	app.Get("/FindProductsBySearch", api.FindAllProductsSearch)
	app.Get("/FindAllCategoriesSearchApi", api.FindAllCategoriesSearchApi)
	app.Post("/RestoreProduct", api.RestoreProductApi)
	app.Get("/FindProductByProductId", api.FindProductByID)
	app.Get("/FindBrandsBySearch", api.FindAllBrandsSearch)
	app.Get("/FindAllSuppliersSearch", api.FindAllSuppliersSearch)
	app.Put("/UpdateProduct", api.UpdateProductApi)
}
