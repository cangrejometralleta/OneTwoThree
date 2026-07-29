package main

import (
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// main Casts the Players, then Steps off the Stage.
func main() {
	api := PeopleAPI{Store: GormPeople{DB: OpenPeopleDatabase()}}

	http.HandleFunc("GET /people", api.ListPeopleRecords)
	http.HandleFunc("POST /people", api.AddPersonRecord)
	http.HandleFunc("GET /people/{id}", api.ShowPersonRecord)
	http.HandleFunc("PUT /people/{id}", api.SavePersonRecord)
	http.HandleFunc("DELETE /people/{id}", api.DropPersonRecord)

	log.Println("✅ People Service Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// OpenPeopleDatabase Opens SQLite and Shapes the Table.
func OpenPeopleDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("people.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database Refused to Open: %v", err)
	}

	if err := db.AutoMigrate(&PersonRow{}); err != nil {
		log.Fatalf("❌ Schema Refused to Migrate: %v", err)
	}

	return db
}
