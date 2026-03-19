package db

import (
	"database/sql"

	"github.com/katy248/ship-cloud-auth/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func Setup() {
	dsn := config.Config.GetString("database-url")
	if dsn == "" {
		panic("database-url not configured")
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
	))
	DB = bun.NewDB(sqldb, pgdialect.New())
}
