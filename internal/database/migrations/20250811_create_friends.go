package migrations

import (
	"context"
	"redsoft-test-task/internal/database"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(createUsersTable, dropUsersTable)
}

func createFriendsTable(ctx context.Context, db *bun.DB) error {
	// Enable uuid-ossp extension for UUID generation.
	_, err := db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	if err != nil {
		return err
	}

	// Create users table based on the User model.
	_, err = db.NewCreateTable().
		Model((*database.User)(nil)).
		IfNotExists().
		Exec(ctx)
	return err
}

func dropFriendsTable(ctx context.Context, db *bun.DB) error {
	// Drop users table if it exists.
	_, err := db.NewDropTable().
		Model((*database.User)(nil)).
		IfExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	// Optionally drop uuid-ossp extension (be cautious if other tables use it).
	_, err = db.ExecContext(ctx, "DROP EXTENSION IF EXISTS \"uuid-ossp\"")
	return err
}
