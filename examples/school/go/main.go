package main

import (
	"log"
	"log/slog"
	"strconv"

	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/app"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/campus"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/handler"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving/chiserver"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/serving/ginserver"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/settings"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/store"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/tokens"
	"github.com/cangrejometralleta/OneTwoThree/examples/school/go/transport"
)

// main Casts the Players, then Steps off the Stage.
// It is the only File that Knows every Package by Name.
func main() {
	config, err := settings.LoadSchoolConfig(settings.FindSchoolDataRoot(), settings.ReadProcessEnvironment())
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

	schoolHandler := handler.Handler{
		School: app.SchoolService{
			Students: database, Courses: database,
			Registry: office, Notices: office, Logger: slog.Default(),
		},
		Tokens: tokens.BuildAccessTokens(config.TokenSecret, config.TokenLifeSeconds),
	}
	server := SelectServerAdapter(config.ServerAdapter)
	address := ":" + strconv.Itoa(config.Port)

	log.Printf("✅ School Listening on %s through %q", address, config.ServerAdapter)
	log.Fatal(server.ServeRoutes(schoolHandler.DeclareSchoolRoutes(), address))
}

// SelectServerAdapter Picks the Framework at Startup, never at Compile Time.
// It Names every Adapter, so it Lives in the Composition.
func SelectServerAdapter(name string) transport.Server {
	switch name {
	case "gin":
		return ginserver.GinServer{}
	case "chi":
		return chiserver.ChiServer{}
	default:
		return serving.StdlibServer{}
	}
}
