package db

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"sourcecraft.dev/organization-shipmonitor/ship-cloud-auth/config"
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
