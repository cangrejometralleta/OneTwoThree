package main

import (
	"log"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// main Casts the Players, then Steps off the Stage.
func main() {
	config, err := LoadSchoolConfig("..", readProcessEnvironment())
	if err != nil {
		log.Fatalf("❌ School Configuration Failed: %v", err)
	}

	school := GormSchool{DB: OpenSchoolDatabase(config.DatabasePath)}
	api := SchoolAPI{Students: school, Courses: school, Tokens: BuildTokenIssuer(config)}
	server := SelectServerAdapter(config.ServerAdapter)
	address := ":" + strconv.Itoa(config.Port)

	log.Printf("✅ School Listening on %s through %q", address, config.ServerAdapter)
	log.Fatal(server.ServeRoutes(api.DeclareSchoolRoutes(), address))
}

// BuildTokenIssuer Uses the Secret and Lifetime Validated at Startup.
func BuildTokenIssuer(config SchoolConfig) TokenIssuer {
	return HmacTokens{
		Secret: []byte(config.TokenSecret),
		Life:   time.Duration(config.TokenLifeSeconds) * time.Second,
		Now:    time.Now,
	}
}

// OpenSchoolDatabase Opens SQLite and Shapes both Tables.
// Reference: https://gorm.io/docs/connecting_to_the_database.html
func OpenSchoolDatabase(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database Refused to Open: %v", err)
	}

	if err := db.AutoMigrate(&StudentRow{}, &CourseRow{}); err != nil {
		log.Fatalf("❌ Schema Refused to Migrate: %v", err)
	}

	return db
}
