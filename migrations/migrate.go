package main

import (
	"database/sql"
	"embed"
	"fmt"
	"path"
	"runtime"

	"github.com/pressly/goose/v3"
	"github.com/url-shortner/config"
	"github.com/url-shortner/connection"
)

//go:embed scripts/*.sql
var embedMigrations embed.FS

type Migration interface {
	Up() error
	Down() error
}

type migration struct {
	db *sql.DB
}

var (
	_, b, _, _ = runtime.Caller(0)
	d          = path.Join(path.Dir(b))
)

func NewMigration(
	db *sql.DB,
) Migration {
	return &migration{
		db: db,
	}
}

func (s *migration) Up() error {
	goose.SetBaseFS(embedMigrations)
	return goose.Up(s.db, "scripts", goose.WithAllowMissing())
}

func (s *migration) Down() error {
	goose.SetBaseFS(embedMigrations)
	return goose.Down(s.db, "scripts", goose.WithAllowMissing())
}

func main() {
	config := config.GetAppConfigFromEnv()
	postgres := connection.NewPostgresConnection(config.DBConfig)
	db := postgres.CreateDB()
	sqlDB := postgres.GetRawDB(db)
	fmt.Println("Start migration ...")
	migration := NewMigration(sqlDB)
	err := migration.Up()
	if err != nil {
		panic(err)
	}
}
