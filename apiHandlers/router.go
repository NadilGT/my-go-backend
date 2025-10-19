package apiHandlers

import (
	"employee-crud/api"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber on Render!")
	})

	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
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
	app.Put("/UpdateCategory", api.UpdateCategoryApi)
	app.Get("/FindAllDeletedProducts", api.FindAllDeletedProductsApi)
	app.Post("/CreateGRN", api.CreateGRN)
	app.Get("/FindAllGRNs", api.FindAllGRNs)
	app.Get("/FindGRNById", api.FindGRNByIdApi)
	app.Get("/GetGRNReport", api.GetGRNReportApi)
	app.Get("/GetGRNReportPDF", api.GetGRNReportPDFApi)
	app.Get("/GetTotalGRNsCount", api.GetTotalGRNsCount)
	app.Get("/GetCompletedGRNsCount", api.GetCompletedGRNsCount)
	app.Get("/GetPendingGRNsCount", api.GetPendingGRNsCount)
	app.Get("/GetPartialReceivedGRNsCount", api.GetPartialReceivedGRNsCount)
	app.Put("/UpdateGRNStatus", api.UpdateGRNStatusApi)
	app.Get("/FindAllProductsBySubCategory", api.GetAllProductsBySubCategoryApi)
	app.Put("/UpdateSupplier", api.UpdateSupplierApi)
	app.Get("/CalculateTotalCost", api.CalculateTotalAndExpectedCost)
	app.Put("/UpdateBrand", api.UpdateBrandApi)
	app.Get("/CalculateBrandCostSummary", api.GetBrandCostSummaryApi)
	app.Delete("/DeleteProductPermanent", api.DeleteProductPermanentApi)

	// Sales Management Routes
	app.Post("/CreateSale", api.CreateSaleApi)
	app.Get("/FindAllSales", api.FindAllSalesApi)
	app.Get("/FindSaleById", api.FindSaleByIdApi)
	app.Post("/CalculateOrderSummary", api.CalculateOrderSummaryApi)
	app.Post("/CalculateChange", api.CalculateChangeApi)
	app.Get("/GetDailySalesSummary", api.GetDailySalesSummaryApi)
	app.Get("/GetDailySalesSummaryPDF", api.GetDailySalesSummaryPDFApi)

	// Saved Daily Reports Routes
	app.Get("/GetSavedDailyReport", api.GetSavedDailyReportApi)
	app.Get("/GetMonthlyReports", api.GetMonthlyReportsApi)
	app.Get("/GetDateRangeReportsPDF", api.GetDateRangeReportsPDFApi)

	// Stock Management Routes
	app.Post("/SyncStocks", api.SyncStocksApi)                     // Sync all product stocks to Stocks collection
	app.Get("/FindAllStocks", api.FindAllStocksApi)                // Get all stocks with pagination (includes total count)
	app.Get("/FindAllStocksLite", api.FindAllStocksLightweightApi) // Get all stocks with pagination (lightweight, no total count)

}
