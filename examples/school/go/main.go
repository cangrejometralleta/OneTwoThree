package main

import (
	"log"
	"strconv"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/api"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/settings"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/store"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/tokens"
)

// main Casts the Players, then Steps off the Stage.
// It is the only File that Knows every Package by Name.
func main() {
	config, err := settings.LoadSchoolConfig(settings.FindSchoolDataRoot(), settings.ReadProcessEnvironment())
	if err != nil {
		log.Fatalf("❌ School Configuration Failed: %v", err)
	}

	registry, err := store.OpenSchoolStore(config.DatabasePath)
	if err != nil {
		log.Fatalf("❌ School Database Refused to Open: %v", err)
	}

	service := api.SchoolAPI{
		Students: registry,
		Courses:  registry,
		Tokens:   tokens.BuildAccessTokens(config.TokenSecret, config.TokenLifeSeconds),
	}
	server := serving.SelectServerAdapter(config.ServerAdapter)
	address := ":" + strconv.Itoa(config.Port)

	log.Printf("✅ School Listening on %s through %q", address, config.ServerAdapter)
	log.Fatal(server.ServeRoutes(service.DeclareSchoolRoutes(), address))
}
