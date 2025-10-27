package connection

import (
	"database/sql"
	"fmt"

	"github.com/url-shortner/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type PostgresConnection interface {
	CreateDB() *gorm.DB
	GetRawDB(db *gorm.DB) *sql.DB
}

type postgresConnection struct {
	DBConfig config.DBConfig
	DB       *gorm.DB
	RawDB    *sql.DB
}

func NewPostgresConnection(
	DBConfig config.DBConfig,
) PostgresConnection {
	return postgresConnection{
		DBConfig: DBConfig,
	}
}

func (s postgresConnection) CreateDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(s.getGormPostgresUrl()), s.getGormConfig())
	if err != nil {
		panic(err)
	}
	s.DB = db
	s.DB = s.initPlugin()
	return s.DB
}

func (s postgresConnection) GetRawDB(db *gorm.DB) *sql.DB {
	rawDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	s.RawDB = rawDb
	return s.RawDB
}

func (s postgresConnection) getGormConfig() *gorm.Config {
	logMode := logger.Silent
	if s.DBConfig.DbLogEnable {
		logMode = logger.Info
	}
	return &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	}
}

func (s postgresConnection) getGormPostgresUrl() string {
	return fmt.Sprintf(
		"host=%s user=%s port=%d dbname=%s sslmode=disable password=%s",
		s.DBConfig.DbHost,
		s.DBConfig.DbUser,
		s.DBConfig.DbPort,
		s.DBConfig.DbName,
		s.DBConfig.DbPass,
	)
}

func (s postgresConnection) initPlugin() *gorm.DB {
	if err := s.DB.Use(
		dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{s.DB.Dialector},
			Policy:  dbresolver.RandomPolicy{},
		}).
			SetConnMaxIdleTime(s.DBConfig.DbMaxIdleConnectionLifetime).
			SetConnMaxLifetime(s.DBConfig.DbMaxConnectionLifetime).
			SetMaxIdleConns(s.DBConfig.DbMaxIdleConnection).
			SetMaxOpenConns(s.DBConfig.DbMaxIdleConnection),
	); err != nil {
		panic(err)
	}
	return s.DB
}
