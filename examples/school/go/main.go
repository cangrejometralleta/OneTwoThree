package main

import (
	"log"
	"log/slog"
	"maps"
	"slices"
	"strconv"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/campus"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/handler"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/school"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving/chiserver"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving/ginserver"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/settings"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/store"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/tokens"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// serverAdapters is every Framework this Runtime Implements; the Loader only Receives the Names.
var serverAdapters = map[string]transport.Server{
	"stdlib": serving.StdlibServer{},
	"chi":    chiserver.ChiServer{},
	"gin":    ginserver.GinServer{},
}

// main Casts the Players, then Steps off the Stage.
// It is the only File that Knows every Package by Name.
func main() {
	root, err := settings.FindSchoolDataRoot(".")
	if err != nil {
		log.Fatalf("❌ School Data Not Found: %v", err)
	}

	constants, err := settings.ReadGlobalConstants(root)
	if err != nil {
		log.Fatalf("❌ School Constants Failed: %v", err)
	}

	names := slices.Sorted(maps.Keys(serverAdapters))
	config, err := settings.LoadSchoolConfig(root, settings.ReadProcessEnvironment(), names)
	if err != nil {
		log.Fatalf("❌ School Configuration Failed: %v", err)
	}

	database, err := store.OpenSchoolStore(config.DatabasePath)
	if err != nil {
		log.Fatalf("❌ School Database Refused to Open: %v", err)
	}

	office, err := campus.OpenOffice(config.RegistryPath, slog.Default())
	if err != nil {
		log.Fatalf("❌ School Registry Refused to Open: %v", err)
	}

	door := handler.Handler{
		School: app.SchoolService{
			Rules:    school.SchoolRules{MinimumAgeYears: constants.MinimumAgeYears},
			Students: database, Courses: database,
			Registry: office, Notices: office, Logger: slog.Default(),
		},
		Tokens: tokens.BuildAccessTokens(config.TokenSecret, config.TokenLifeSeconds),
	}
	address := ":" + strconv.Itoa(config.Port)

	log.Printf("✅ School Listening on %s through %q", address, config.ServerAdapter)
	log.Fatal(serverAdapters[config.ServerAdapter].ServeRoutes(door.DeclareSchoolRoutes(), address))
}
