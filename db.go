package guardmech

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

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

	// avoid recycle of too idle connection
	d.SetConnMaxLifetime(time.Second * 60)

	db = d
}

func GetConn(ctx context.Context) (*sql.Conn, error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Println("Could Not Get Conn")
		return nil, err
	}

	err = conn.PingContext(ctx)
	if err != nil {
		log.Println("Failed to Ping")
		defer conn.Close()
		return nil, err
	}

	return conn, nil
}
