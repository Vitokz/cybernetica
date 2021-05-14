package repository

import (
	"main/config"

	"github.com/go-pg/pg"
)

type PgClient interface { //Интерфеис с функциями подключения и закрытия подключение с бд
	GetConnection() *pg.DB
	Close() error
}

type Pg struct { //Структура с бд
	Db *pg.DB
}

func (p *Pg) GetConnection() *pg.DB {
	return p.Db
}

func (p *Pg) Close() error {
	return p.Db.Close()
}

func NewPgClient() PgClient {
	db := pg.Connect(&config.Database)
	return &Pg{
		Db: db,
	}
}
