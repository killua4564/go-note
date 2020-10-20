package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/killua4564/go-note/app/account"
	"github.com/killua4564/go-note/app/note"
	"github.com/killua4564/go-note/config"

	_ "github.com/go-sql-driver/mysql"
)

func initDB() (db *sql.DB) {
	dbcfg := config.Database{}
	if err := env.Parse(&dbcfg); err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbcfg.Username, dbcfg.Password, dbcfg.Hostname, dbcfg.DBname))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return
}

func main() {
	db := initDB()

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	router := engine.Group("/api", func(c *gin.Context) {})

	account.AccountService(router, db, &cfg)
	note.NoteService(router, db, &cfg)

	engine.Run(":8080")
}
