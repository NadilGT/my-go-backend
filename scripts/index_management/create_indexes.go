// Create/Update Indexes Script
// This is a standalone utility to create/update indexes on the Stocks collection
// Run with: go run scripts/create_indexes.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database configuration
const DATABASE_URL = "mongodb+srv://admin:W6ptbj7HPS3RJ4cU@cluster0.tgypip5.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const DATABASE_NAME = "POS"

func main() {
	fmt.Println("=== Create/Update Stocks Indexes ===")
	fmt.Println("Starting at:", time.Now().Format(time.RFC3339))

	// Connect to MongoDB
	fmt.Println("\nConnecting to MongoDB Atlas...")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(DATABASE_URL))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	fmt.Println("✓ Connected to MongoDB successfully")

	database := client.Database(DATABASE_NAME)
	collection := database.Collection("Stocks")

	// List existing indexes
	fmt.Println("\n[Step 1/3] Checking existing indexes...")
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		log.Fatal("Failed to list indexes:", err)
	}
	defer cursor.Close(ctx)

	var existingIndexes []bson.M
	if err := cursor.All(ctx, &existingIndexes); err != nil {
		log.Fatal("Failed to decode indexes:", err)
	}

	fmt.Println("Existing indexes:")
	for _, index := range existingIndexes {
		fmt.Printf("  - %v\n", index["name"])
	}

	// Drop old single-field index if it exists
	fmt.Println("\n[Step 2/3] Dropping old indexes if they exist...")
	indexNames := []string{"updated_at_desc", "updated_at_1", "updated_at_-1"}
	for _, name := range indexNames {
		_, err := collection.Indexes().DropOne(ctx, name)
		if err != nil {
			// Ignore errors (index might not exist)
			fmt.Printf("  - Index '%s' not found (ok)\n", name)
		} else {
			fmt.Printf("  ✓ Dropped old index: %s\n", name)
		}
	}

	// Create new indexes
	fmt.Println("\n[Step 3/3] Creating optimized indexes...")

	// Index 1: Compound index (updated_at, _id) for efficient cursor pagination
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

	createdIndexes, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{index1, index2, index3})
	if err != nil {
		// Try creating individually if batch fails
		fmt.Println("Batch creation failed, trying individually...")

		// Try index 1
		if name, err := collection.Indexes().CreateOne(ctx, index1); err == nil {
			fmt.Printf("  ✓ Created: %s\n", name)
		} else {
			fmt.Printf("  - Index 1 already exists or error: %v\n", err)
		}

		// Try index 2
		if name, err := collection.Indexes().CreateOne(ctx, index2); err == nil {
			fmt.Printf("  ✓ Created: %s\n", name)
		} else {
			fmt.Printf("  - Index 2 already exists or error: %v\n", err)
		}

		// Try index 3
		if name, err := collection.Indexes().CreateOne(ctx, index3); err == nil {
			fmt.Printf("  ✓ Created: %s\n", name)
		} else {
			fmt.Printf("  - Index 3 already exists or error: %v\n", err)
		}
	} else {
		for _, name := range createdIndexes {
			fmt.Printf("  ✓ Created: %s\n", name)
		}
	}

	// Verify final indexes
	fmt.Println("\n[Verification] Final index list:")
	cursor, err = collection.Indexes().List(ctx)
	if err != nil {
		log.Fatal("Failed to list indexes:", err)
	}
	defer cursor.Close(ctx)

	var finalIndexes []bson.M
	if err := cursor.All(ctx, &finalIndexes); err != nil {
		log.Fatal("Failed to decode indexes:", err)
	}

	for _, index := range finalIndexes {
		fmt.Printf("  ✓ %v: %v\n", index["name"], index["key"])
	}

	// Get collection stats
	var stats bson.M
	if err := database.RunCommand(ctx, bson.D{{Key: "collStats", Value: "Stocks"}}).Decode(&stats); err == nil {
		if count, ok := stats["count"].(int32); ok {
			fmt.Printf("\n✅ Success! Collection has %d documents\n", count)
		}
	}

	fmt.Println("\n🎉 Index optimization complete!")
	fmt.Println("Your Stocks collection is now ready for 10k+ products!")
}
