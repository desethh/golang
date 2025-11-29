package dbs

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DBopen() {
	var err error
	DB, err = sql.Open("mysql", "tauren91_itastan:9pV*taaN%baU@tcp(tauren91.beget.tech:3306)/tauren91_itastan?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("База недоступна:", err)
	}
}
