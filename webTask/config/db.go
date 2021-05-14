package config

import (
	"github.com/go-pg/pg"
)

var Database = pg.Options{
	User:     "postgres",
	Password: "root",
	Database: "WebTask",
	Addr:     "localhost:5432",
}
