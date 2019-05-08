package app

import (
	"github.com/jmoiron/sqlx"
)

func NewDB(log Logger, dbDriver string, dbUrl string) *sqlx.DB {
	db, err := sqlx.Open(dbDriver, dbUrl)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	return db
}
