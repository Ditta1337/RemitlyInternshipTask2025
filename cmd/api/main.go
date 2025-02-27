package main

import (
	dbPkg "github.com/Ditta1337/RemitlyInternshipTask2025/internal/db"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/env"
	storePkg "github.com/Ditta1337/RemitlyInternshipTask2025/internal/store"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Remitly SWIFT API 2025
//	@description	Remitly 2025 internship task

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// loading .env
	if err := godotenv.Load(); err != nil {
		logger.Warnf("error loading .env file: %s", err.Error())
	}

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:remitly2025@localhost/swift?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:        env.GetString("ENV", "development"),
		apiVersion: env.GetString("API_VERSION", "v1"),
	}

	// db connection
	db, err := dbPkg.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("db connection established")

	// migrations
	migrationsDir := env.GetString("GOOSE_MIGRATION_DIR", "./cmd/migrations")
	if err := goose.Up(db, migrationsDir); err != nil {
		logger.Fatalf("failed to run migrations: %s", err.Error())
	}

	store := storePkg.NewPostgresStorage(db)

	// seed db if its empty
	if err := dbPkg.SeedDBIfEmpty(db, store); err != nil {
		logger.Infof("failed to seed db: %s", err.Error())
	} else {
		logger.Info("finished seeding db")
	}

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
