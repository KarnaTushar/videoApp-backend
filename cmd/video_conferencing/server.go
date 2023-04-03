package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

	"github.com/KarnaTushar/videoApp-backend/internal/generated"
	"github.com/KarnaTushar/videoApp-backend/migrations"
	"github.com/KarnaTushar/videoApp-backend/pkg/graph"
	"github.com/KarnaTushar/videoApp-backend/pkg/middleware"
	"github.com/KarnaTushar/videoApp-backend/pkg/models"
	"github.com/KarnaTushar/videoApp-backend/services"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
	"github.com/rs/zerolog/hlog"

	"github.com/KarnaTushar/videoApp-backend/utils"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"

	_ "github.com/go-sql-driver/mysql"
)

// const defaultPort = "8080"

func main() {
	configDir := flag.String("config", ".", "Directory which contains the config.json")
	if configDir != nil {
		println("Config Path", *configDir)
	} else {
		println(fmt.Errorf("config dir is nil"))
	}

	utils.SetupConfig(configDir)

	logger := utils.Configure(utils.Config{
		ConsoleLoggingEnabled: viper.GetBool("ENABLE_CONSOLE_LOGGING"),
		FileLoggingEnabled:    viper.GetBool("ENABLE_FILE_LOGGING"),
		Directory:             viper.GetString("LOG_DIR"),
		MaxSize:               viper.GetInt("MAX_LOG_SIZE"),
		MaxBackups:            viper.GetInt("MAX_LOG_BACKUP"),
		MaxAge:                viper.GetInt("MAX_LOG_AGE"),
		Filename:              "app-builder-logs",
	})

	port := viper.GetString("PORT")

	database, err := models.CreateDB(viper.GetString("DATABASE_URL"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Error initializing database")
		return
	}

	defer database.Close()

	if viper.GetBool("RUN_MIGRATION") {
		migrations.RunMigration(configDir)
	}

	router := mux.NewRouter()

	config := generated.Config{
		Resolvers: &graph.Resolver{
			DB:     database,
			Logger: logger,
		},
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	requestHandler := services.ServiceRouter{
		DB:     database,
		Logger: logger,
	}

	router.HandleFunc("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.HandleFunc("/oauth", http.HandlerFunc(requestHandler.OAuth))
	router.HandleFunc("/pstn", http.HandlerFunc(requestHandler.PSTN))

	router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		logger.Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	logger.Info().Str("origin", viper.GetString("ALLOWED_ORIGIN")).Msg("")
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{viper.GetString("ALLOWED_ORIGIN")},
		AllowCredentials: true,
		AllowedHeaders:   []string{"authorization", "content-type"},
		Debug:            false,
	}).Handler)
	router.Use(handlers.RecoveryHandler())

	router.Use(middleware.AuthHandler(database, logger))

	if viper.GetBool("ENABLE_NEWRELIC_MONITORING") {
		nrAgent, err := newrelic.NewApplication(
			newrelic.ConfigAppName(viper.GetString("NEWRELIC_APPNAME")),
			newrelic.ConfigLicense(viper.GetString("NEWRELIC_LICENSE")),
			newrelic.ConfigDebugLogger(os.Stdout),
		)

		if err != nil {
			logger.Fatal().Err(err).Msg("Error initializing New Relic Agent")
			return
		}

		router.Use(nrgorilla.Middleware(nrAgent))
	}

	logger.Debug().Str("PORT", port)
	logger.Fatal().Err(http.ListenAndServe(":"+port, router))
}
