package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	loadenv "go-script-moph-appoint/src/loadenv"
)

var Db *sql.DB

func ConnectDB() *sql.DB {
	if Db != nil {
		return Db
	}

	var err error
	Db, err = sql.Open("postgres", loadenv.LoadDBconnec())
	if err != nil {
		log.Fatal(err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("connect database success")
	return Db
}
