package guardmech

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	initDB()
}

func initDB() {
	addr := os.Getenv("DB_HOST")
	if os.Getenv("DB_PORT") != "" {
		addr += ":" + os.Getenv("DB_PORT")
	}

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = addr
	cfg.DBName = os.Getenv("DB_NAME")
	dsn := cfg.FormatDSN()
	log.Println("connecing to", dsn)
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	log.Println("connecting db OK")

	db = d
}
