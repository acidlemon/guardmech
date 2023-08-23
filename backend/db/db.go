package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var pool *sql.DB

func init() {
	initDB()
}

func initDB() {
	addr := os.Getenv("GUARDMECH_DB_HOST")
	if os.Getenv("GUARDMECH_DB_PORT") != "" {
		addr += ":" + os.Getenv("GUARDMECH_DB_PORT")
	}

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("GUARDMECH_DB_USER")
	cfg.Passwd = os.Getenv("GUARDMECH_DB_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = addr
	cfg.DBName = os.Getenv("GUARDMECH_DB_NAME")
	dsn := cfg.FormatDSN()
	log.Println("connecing to", strings.Replace(dsn, cfg.Passwd, strings.Repeat("*", len(cfg.Passwd)), -1))
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	log.Println("connecting db OK")

	// avoid recycle of too idle connection
	d.SetConnMaxLifetime(time.Minute * 5)
	d.SetMaxOpenConns(20)

	pool = d
}

func GetConn(ctx context.Context) (*sql.Conn, error) {
	conn, err := pool.Conn(ctx)
	if err != nil {
		log.Println("Could Not Get Conn:", err)
		return nil, err
	}

	err = conn.PingContext(ctx)
	if err != nil {
		log.Println("Failed to Ping:", err)
		defer conn.Close()
		return nil, err
	}

	return conn, nil
}

func Begin(ctx context.Context, conn *sql.Conn) (*Tx, error) {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Tx{Tx: tx}, nil
}

func GetTxConn(ctx context.Context) (*sql.Conn, *Tx, error) {
	conn, err := GetConn(ctx)
	if err != nil {
		return nil, nil, err
	}

	tx, err := Begin(ctx, conn)
	if err != nil {
		return nil, nil, err
	}

	return conn, tx, nil
}
