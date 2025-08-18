package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"user-service/database/seeds"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PsqlDB.User,
		cfg.PsqlDB.Password,
		cfg.PsqlDB.Host,
		cfg.PsqlDB.Port,
		cfg.PsqlDB.DbName,
	)
	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-1] failed to coonect to database" + cfg.PsqlDB.Host)
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-2] failed to get database instance")
		return nil, err
	}

	seeds.SeedRole(db)
	seeds.SeedAdmin(db)

	sqlDb.SetMaxOpenConns(cfg.PsqlDB.DBMaxOpen)
	sqlDb.SetMaxIdleConns(cfg.PsqlDB.DBMaxIdle)

	return &Postgres{DB: db}, nil
}
