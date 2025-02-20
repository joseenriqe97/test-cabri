package config

import (
	"database/sql"

	"github.com/ansel1/merry"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
)

var (
	SQLDB *sql.DB
	// Guarda el dialecto(Lenguaje sql) con que se trabajara
	SQLDBGoqu goqu.Database
)

func NewDBConnection() error {

	psqlInfo := PgConn()
	err := newDB(psqlInfo)
	if err != nil {
		return merry.Wrap(err)
	}
	return nil
}

func newDB(psqlInfo string) error {

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return merry.Wrap(err)
	}
	err = db.Ping()
	if err != nil {
		return merry.Wrap(err)
	}

	SQLDBGoqu = *goqu.Dialect("postgres").DB(db)
	SQLDB = db

	return nil
}

func CloseConnections() {
	err := SQLDB.Close()
	if err != nil {
		merry.New("can not close database connection")
	}

}
