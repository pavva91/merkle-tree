package db

import (
	"fmt"
	"github.com/pavva91/merkle-tree/server/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var (
	ORM      ORMer = orm{}
	database *gorm.DB
)

type ORMer interface {
	MustConnectToDB(cfg config.ServerConfig)
	GetDB() *gorm.DB
}

type orm struct{}

func (o orm) MustConnectToDB(cfg config.ServerConfig) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Println(err.Error())
		panic(fmt.Errorf("error connecting db: %w", err))
	}

	db.Logger.LogMode(logger.Info)
	database = db
}

func (o orm) GetDB() *gorm.DB {
	return database
}
