package model

import (
	"fmt"
	"log"

	"/golang/mtfosbot/module/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var x *sqlx.DB

// NewDB - connect to database
func NewDB() (*sqlx.DB, error) {
	var err error
	conf := config.GetConf()
	userPassStr := conf.Database.User
	if len(conf.Database.Pass) > 0 {
		userPassStr += ":" + conf.Database.Pass
	}
	connStr := fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable", userPassStr, conf.Database.Host, conf.Database.DB)
	x, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	x.SetMaxIdleConns(10)
	x.SetMaxOpenConns(200)
	err = x.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return x, err
}

// GetDB -
func GetDB() *sqlx.DB {
	return x
}
