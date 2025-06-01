package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	dsn := "web_dev:WebDev@123@tcp(103.127.99.152:3306)/webdmf"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak bisa diakses:", err)
	}
	fmt.Println("Berhasil koneksi ke database")
}
