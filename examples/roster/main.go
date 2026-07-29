package main

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// main Casts the Players, then Steps off the Stage.
func main() {
	api := RosterAPI{Store: GormPeople{DB: OpenPeopleDatabase()}}

	choice := os.Getenv("SERVER")
	server := SelectServerAdapter(choice)

	log.Printf("✅ Roster Listening on :8080 through %q", choice)
	log.Fatal(server.ServeRoutes(api.DeclareRosterRoutes(), ":8080"))
}

// OpenPeopleDatabase Opens SQLite and Shapes the Table.
func OpenPeopleDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("roster.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database Refused to Open: %v", err)
	}

	if err := db.AutoMigrate(&PersonRow{}); err != nil {
		log.Fatalf("❌ Schema Refused to Migrate: %v", err)
	}

	return db
}
