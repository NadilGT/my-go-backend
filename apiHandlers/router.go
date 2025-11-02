package apiHandlers

import (
	"employee-crud/api"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Returns Monthly PDF Report
	app.Get("/GetMonthlyReturnsPDF", api.GetMonthlyReturnsReportPDF)
	// Expiring Stocks Report Route
	app.Get("/GetExpiringStocksReportPDF", api.GetExpiringStocksReportPDF)
	// Top 10 Expiring Stocks in 7 Days (JSON)
	app.Get("/GetExpiringStocksNext7Days", api.GetExpiringStocksNext7Days)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Fiber on Render!")
	})

	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	app.Post("/CreateCategory", api.CreateCategoryApi)
	app.Get("/FindAllCategory", api.FindAllCategoriesApi)
	app.Delete("/DeleteCategory", api.DeleteCategoryApi)

	// Category Enhancement Routes
	app.Get("/api/products/count", api.GetCategorizedProductsCountApi)                       // Get total categorized products count
	app.Get("/api/categories/:categoryId/products/count", api.GetProductsCountByCategoryApi) // Get product count by category
	app.Get("/api/categories/search", api.SearchCategoriesApi)                               // Search categories

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

	// Update supplier status (active/inactive)
	app.Put("/UpdateSupplierStatus", api.UpdateSupplierStatus)

	// Get total number of active and inactive suppliers
	app.Get("/GetSupplierStatusCounts", api.GetSupplierStatusCounts)
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

	// Total Products Count API
	app.Get("/GetTotalProducts", api.GetTotalProducts)
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
	app.Post("/SyncStocks", api.SyncStocksApi)                                     // Sync all product stocks to Stocks collection
	app.Get("/FindAllStocks", api.FindAllStocksApi)                                // Get all stocks with pagination (includes total count)
	app.Get("/FindAllStocksLite", api.FindAllStocksLightweightApi)                 // Get all stocks with pagination (lightweight, no total count)
	app.Get("/FindAllProductsWithStock", api.FindAllProductsWithStockApi)          // Get all products with stock info (includes products without batches)
	app.Get("/FindAllProductsWithStockLite", api.FindAllProductsWithStockLiteApi)  // Lightweight version - all products with stock info
	app.Get("/FindAllStocksFiltered", api.FindAllStocksFilteredApi)                // Get filtered stocks by status with pagination (includes total count)
	app.Get("/FindAllStocksFilteredLite", api.FindAllStocksFilteredLightweightApi) // Get filtered stocks by status with pagination (lightweight)
	app.Get("/GetTotalStockQuantity", api.GetTotalStockQuantityApi)                // Get sum of all stockQty (total quantity in inventory)
	app.Get("/GetStockStatusCounts", api.GetStockStatusCountsApi)                  // Get count of stocks by status (Low/Average/Good)

	// Low Stock Products API
	app.Get("/GetLowStockProducts", api.GetLowStockProductsHandler) // Get top 10 lowest stock products

	// Stock Maintenance Routes
	app.Delete("/CleanupOrphanedStocks", api.CleanupOrphanedStocksApi) // Remove stock entries with null/empty batchId
	app.Get("/ValidateStockIntegrity", api.ValidateStockIntegrityApi)  // Check for stock data inconsistencies

	// Batch Stock Management Routes
	app.Post("/AddStock", api.AddStock)                // Add stock to existing product (adds to existing batch or creates new batch based on expiry date)
	app.Put("/EditBatchStock", api.EditBatchStock)     // Edit stock quantity of a specific batch
	app.Put("/EditBatchDetails", api.EditBatchDetails) // Edit batch details (expiry date, prices)
	app.Put("/RemoveStock", api.RemoveStockFromBatch)  // Remove/reduce stock from a specific batch
	app.Delete("/DeleteBatch", api.DeleteBatch)        // Delete a batch completely

	// Return APIs
	app.Post("/returns", api.CreateReturnApi)
	app.Get("/returns", api.FindAllReturnsApi)
	app.Get("/returns/:id", api.FindReturnByIdApi)
}
