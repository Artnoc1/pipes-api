package main

import (
	"database/sql"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/bugsnag/bugsnag-go"
	_ "github.com/lib/pq"

	"github.com/toggl/pipes-api/pkg/authorization"
	"github.com/toggl/pipes-api/pkg/autosync"
	"github.com/toggl/pipes-api/pkg/config"
	"github.com/toggl/pipes-api/pkg/connection"
	"github.com/toggl/pipes-api/pkg/pipe"
	"github.com/toggl/pipes-api/pkg/server"
	"github.com/toggl/pipes-api/pkg/toggl/client"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().Unix())

	envFlags := config.Flags{}
	config.ParseFlags(&envFlags)

	cfg := config.Load(envFlags.Environment, envFlags.WorkDir)
	corsWhitelist := cfg.Urls.CorsWhitelist[cfg.EnvType]
	pipesApiHost := cfg.Urls.PipesAPIHost[cfg.EnvType]
	togglApiHost := cfg.Urls.TogglAPIHost[cfg.EnvType]

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       envFlags.BugsnagAPIKey,
		ReleaseStage: envFlags.Environment,
		NotifyReleaseStages: []string{
			config.EnvTypeProduction,
			config.EnvTypeStaging,
		},
		// more configuration options
	})

	db, err := sql.Open("postgres", envFlags.DbConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api := client.NewTogglApiClient(togglApiHost)

	authStore := authorization.NewStorage(db, cfg)
	connStore := connection.NewStorage(db)
	pipesStore := pipe.NewStorage(db)

	pipesService := pipe.NewService(cfg, authStore, pipesStore, connStore, api, pipesApiHost, cfg.WorkDir)

	autosync.NewService(pipesService).Start()

	router := server.NewRouter(corsWhitelist).AttachHandlers(
		server.NewController(pipesService),
		server.NewMiddleware(api, pipesService),
	)
	server.Start(envFlags.Port, router)
}
