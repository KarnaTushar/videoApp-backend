package migrations

import (
	"database/sql"

	"flag"
	"log"

	"github.com/KarnaTushar/videoApp-backend/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

// RunMigration runs the schema migrations
func RunMigration(config *string) {
	if config == nil {
		configDir := flag.String("config", "..", "Directory which contains the config.json")
		utils.SetupConfig(configDir)
	}

	println("Running Migrations...")

	db, _ := sql.Open("mysql", viper.GetString("DATABASE_URL"))
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Print(err)
	}
	// Uncomment the below lines for migrate down!
	// if err := m.Down(); err != nil {
	// 	log.Fatal(err)
	// }
}
