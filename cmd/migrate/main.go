package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"redsoft-test-task/internal/database/migrations"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

func main() {
	// Initialize database connection.
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DSN"))))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	// Initialize migrator.
	migrator := migrate.NewMigrator(db, migrations.Migrations)

	// Initialize migrations.
	err := migrator.Init(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Apply migrations.
	group, err := migrator.Migrate(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if group == nil {
		log.Fatal(errors.New("nothing to migrate"))
	}

	log.Println("Migrations applied successfully")
}
