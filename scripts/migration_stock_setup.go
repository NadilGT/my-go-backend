package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ProductId  string     `bson:"productId"`
	Name       string     `bson:"name"`
	StockQty   int        `bson:"stockQty"`
	ExpiryDate *time.Time `bson:"expiry_date,omitempty"`
}

type Stock struct {
	ProductId  string     `bson:"productId"`
	Name       string     `bson:"name"`
	StockQty   int        `bson:"stockQty"`
	ExpiryDate *time.Time `bson:"expiry_date,omitempty"`
	UpdatedAt  time.Time  `bson:"updated_at"`
}

// This script should be run ONCE after deploying the new code
// It will:
// 1. Create necessary indexes on the Stocks collection
// 2. Perform initial sync of all products to Stocks collection
//
// Usage: go run scripts/migration_stock_setup.go

func main() {
	fmt.Println("=== Stock Management Migration Script ===")
	fmt.Println("Starting migration at:", time.Now().Format(time.RFC3339))

	// Get MongoDB connection string from environment
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default for local development
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "your_database_name" // Change this to your actual database name
		fmt.Printf("\n⚠ DB_NAME not set, using default: %s\n", dbName)
	}

	// Connect to MongoDB
	fmt.Println("\nConnecting to MongoDB...")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	// Ping to verify connection
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	fmt.Println("✓ Connected to MongoDB successfully")

	database := client.Database(dbName)
	stocksCollection := database.Collection("Stocks")
	productsCollection := database.Collection("Products")
	ctx := context.Background()

	// Step 1: Create indexes
	fmt.Println("\n[Step 1/3] Creating indexes on Stocks collection...")
	if err := createIndexes(ctx, stocksCollection); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}
	fmt.Println("✓ Indexes created successfully")

	// Step 2: Count existing products
	fmt.Println("\n[Step 2/3] Counting products to sync...")
	filter := bson.M{"deleted": false}
	totalProducts, err := productsCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal("Failed to count products:", err)
	}
	fmt.Printf("✓ Found %d products to sync\n", totalProducts)

	if totalProducts == 0 {
		fmt.Println("\n⚠ No products found. Nothing to sync.")
		fmt.Println("\nMigration script completed (no data to migrate)")
		return
	}

	// Step 3: Sync all products to Stocks
	fmt.Println("\n[Step 3/3] Syncing products to Stocks collection...")
	fmt.Println("This may take a few minutes for large datasets...")

	startTime := time.Now()
	if err := syncAllProducts(ctx, productsCollection, stocksCollection, totalProducts); err != nil {
		log.Fatal("Failed to sync products:", err)
	}
	duration := time.Since(startTime)

	// Verify sync
	stockCount, _ := stocksCollection.CountDocuments(ctx, bson.M{})

	fmt.Println("\n=== Migration Complete ===")
	fmt.Printf("✓ Synced %d products to Stocks collection\n", stockCount)
	fmt.Printf("✓ Time taken: %v\n", duration)
	fmt.Printf("✓ Completed at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Println("\n🎉 Your Stock Management system is now ready to use!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Test the API: POST http://your-api/SyncStocks")
	fmt.Println("  2. Fetch stocks: GET http://your-api/FindAllStocks?per_page=15")
	fmt.Println("  3. See STOCK_MANAGEMENT_README.md for full documentation")
}

func createIndexes(ctx context.Context, collection *mongo.Collection) error {
	// Index 1: Compound index (updated_at, _id) for efficient cursor pagination
	// CRITICAL: Must be compound index to support sorting by both fields
	index1 := mongo.IndexModel{
		Keys: bson.D{
			{Key: "updated_at", Value: -1},
			{Key: "_id", Value: -1},
		},
		Options: options.Index().SetName("updated_at_id_desc"),
	}

	// Index 2: productId unique for preventing duplicates
	index2 := mongo.IndexModel{
		Keys: bson.D{{Key: "productId", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("productId_unique"),
	}

	// Index 3: stockQty for efficient filtering by status (low/average/good)
	index3 := mongo.IndexModel{
		Keys:    bson.D{{Key: "stockQty", Value: 1}},
		Options: options.Index().SetName("stockQty_idx"),
	}

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{index1, index2, index3})
	return err
}

func syncAllProducts(ctx context.Context, productsCollection, stocksCollection *mongo.Collection, total int64) error {
	filter := bson.M{"deleted": false}
	batchSize := 500
	skip := 0
	syncedCount := 0

	for {
		// Find products in batches
		findOptions := options.Find()
		findOptions.SetSkip(int64(skip))
		findOptions.SetLimit(int64(batchSize))

		cursor, err := productsCollection.Find(ctx, filter, findOptions)
		if err != nil {
			return err
		}

		var products []Product
		if err := cursor.All(ctx, &products); err != nil {
			cursor.Close(ctx)
			return err
		}
		cursor.Close(ctx)

		if len(products) == 0 {
			break
		}

		// Prepare bulk operations
		var operations []mongo.WriteModel
		currentTime := time.Now()

		for _, product := range products {
			stock := Stock{
				ProductId:  product.ProductId,
				Name:       product.Name,
				StockQty:   product.StockQty,
				ExpiryDate: product.ExpiryDate,
				UpdatedAt:  currentTime,
			}

			filter := bson.M{"productId": product.ProductId}
			update := bson.M{
				"$set": stock,
				"$setOnInsert": bson.M{
					"created_at": currentTime,
				},
			}

			operation := mongo.NewUpdateOneModel()
			operation.SetFilter(filter)
			operation.SetUpdate(update)
			operation.SetUpsert(true)

			operations = append(operations, operation)
		}

		// Execute bulk write
		if len(operations) > 0 {
			_, err := stocksCollection.BulkWrite(ctx, operations)
			if err != nil {
				return err
			}
		}

		syncedCount += len(products)
		fmt.Printf("  Progress: %d/%d (%.1f%%)\n", syncedCount, total, float64(syncedCount)/float64(total)*100)

		skip += batchSize
	}

	return nil
}
