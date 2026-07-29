package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// main Casts the Players, then Steps off the Stage.
func main() {
	school := GormSchool{DB: OpenSchoolDatabase()}
	api := SchoolAPI{Students: school, Courses: school, Tokens: BuildTokenIssuer()}

	choice := os.Getenv("SERVER")
	server := SelectServerAdapter(choice)

	log.Printf("✅ School Listening on :8080 through %q", choice)
	log.Fatal(server.ServeRoutes(api.DeclareSchoolRoutes(), ":8080"))
}

// BuildTokenIssuer Refuses to Start without a Secret.
// A silent Default Secret is a public Key with extra Steps.
func BuildTokenIssuer() TokenIssuer {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		log.Fatal("❌ TOKEN_SECRET is Missing")
	}

	return HmacTokens{Secret: []byte(secret), Life: 10 * time.Minute, Now: time.Now}
}

// OpenSchoolDatabase Opens SQLite and Shapes both Tables.
// Reference: https://gorm.io/docs/connecting_to_the_database.html
func OpenSchoolDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("school.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database Refused to Open: %v", err)
	}

	if err := db.AutoMigrate(&StudentRow{}, &CourseRow{}); err != nil {
		log.Fatalf("❌ Schema Refused to Migrate: %v", err)
	}

	return db
}
