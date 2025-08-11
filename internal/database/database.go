package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DBConfig struct {
	DSN string `env:"DSN"`
}

type ReadModel interface {
	CreateUser(ctx context.Context, user *User) error
}

type readModel struct {
	db *bun.DB
}

func New(dsn string) (ReadModel, error) {
	pgsql := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(pgsql, pgdialect.New())

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return &readModel{db: db}, err
}

func (r *readModel) GetUserBySurname(ctx context.Context, surname string) (*User, error) {
	usr := &User{}
	err := r.db.NewSelect().Model(usr).Where("surname = ?", surname).Scan(ctx, usr)
	return usr, err
}

func (r *readModel) GetUserByID(ctx context.Context, id int) (*User, error) {}

func (r *readModel) GetAllUsers(ctx context.Context, limit, offset int) ([]*User, int, error) {
	var users []*User
	total, err := r.db.NewSelect().Model(users).Limit(limit).Offset(offset).ScanAndCount(ctx, users)
	return users, total, err
}
func (r *readModel) GetUsersFriends(ctx context.Context, id uuid.UUID) ([]*User, error) {
	return nil, nil
}

func (r *readModel) CreateUser(ctx context.Context, user *User) error {
	return nil
}

func (r *readModel) CreateFriendship(ctx context.Context, firstID, secondID uuid.UUID) error {
	return nil
}

func (r *readModel) PatchUser(ctx context.Context, data *User, id int) error {
	return nil
}
