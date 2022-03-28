package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"main.go/conf"
)

var Conn *pgxpool.Pool

var conString = "postgresql://postgres:example@db:5432/?sslmode=disable"

var createCat = "CREATE TABLE IF NOT EXISTS category (id serial PRIMARY KEY, name text, icon text)"
var createExp = "CREATE TABLE IF NOT EXISTS expense (id serial PRIMARY KEY, description text, date timestamp, price int, categoryid int)"

func Connect() {
	var err error
	Conn,err = pgxpool.Connect(context.Background(), conString)
	conf.CheckErrorFatal(err)

	initializeDb()
}

func initializeDb() {
	_, err := Conn.Exec(context.Background(),createCat)
	conf.CheckErrorFatal(err)

	_, err = Conn.Exec(context.Background(), createExp)
	conf.CheckErrorFatal(err)
}

