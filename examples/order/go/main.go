package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// main Casts the Players, then Steps off the Stage.
func main() {
	store := GormOrders{DB: OpenOrderDatabase()}
	api := OrderAPI{Orders: store}

	log.Print("✅ Orders Listening on :8081")
	log.Fatal(ServeOrderRoutes(api.DeclareOrderRoutes(), ":8081"))
}

// OpenOrderDatabase Opens SQLite and Shapes its one Table.
// Reference: https://gorm.io/docs/connecting_to_the_database.html
func OpenOrderDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("order.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database Refused to Open: %v", err)
	}

	if err := db.AutoMigrate(&OrderRow{}); err != nil {
		log.Fatalf("❌ Schema Refused to Migrate: %v", err)
	}

	return db
}
