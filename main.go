package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	// "app"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/killua4564/go-note/utils/config"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbcfg := config.Database{}
	if err := env.Parse(&dbcfg); err != nil {

	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", ))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})
	r.Run()
}
