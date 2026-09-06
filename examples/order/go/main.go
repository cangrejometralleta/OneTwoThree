package main

import (
	"log"
	"strconv"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// main Casts the Players, then Steps off the Stage.
func main() {
	config, err := LoadOrderConfig("..", readProcessEnvironment())
	if err != nil {
		log.Fatalf("❌ Order Configuration Failed: %v", err)
	}

	store := GormOrders{DB: OpenOrderDatabase(config.DatabasePath)}
	api := OrderAPI{Orders: store}

	address := ":" + strconv.Itoa(config.Port)
	log.Printf("✅ Orders Listening on %s", address)
	log.Fatal(ServeOrderRoutes(api.DeclareOrderRoutes(), address))
}

// OpenOrderDatabase Opens SQLite and Shapes its one Table.
// Reference: https://gorm.io/docs/connecting_to_the_database.html
func OpenOrderDatabase(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database Refused to Open: %v", err)
	}

	if err := db.AutoMigrate(&OrderRow{}); err != nil {
		log.Fatalf("❌ Schema Refused to Migrate: %v", err)
	}

	return db
}
