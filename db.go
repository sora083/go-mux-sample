package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func Getenv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	} else {
		return val
	}
}

func init() {
	host := Getenv("DB_HOST", "127.0.0.1")
	port := Getenv("DB_PORT", "23306")
	user := Getenv("DB_USER", "mao")
	pass := Getenv("DB_PASS", "mao")
	name := Getenv("DB_NAME", "chihuahua")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&interpolateParams=true", user, pass, host, port, name)

	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(10 * time.Second)
}
