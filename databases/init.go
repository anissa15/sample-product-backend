package databases

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	db *gorm.DB
}

func New(dsn string) *PostgreSQL {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatal("cannot connect to postgres driver", err)
	}
	err = db.AutoMigrate(
		&Product{},
	)
	if err != nil {
		log.Fatal("cannnot auto migrate", err)
	}
	return &PostgreSQL{db: db}
}
