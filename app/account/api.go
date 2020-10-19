package account

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/killua4564/go-note/config"
	"github.com/killua4564/go-note/utils/hash"
)

type api struct {
	cfg    *config.Config
	loader *loader
}

var url = struct {
	create string
	login  string
}{
	create: "/account",
	login:  "/account/login",
}

func AccountSerivce(router *gin.RouterGroup, dbconn *sql.DB, cfg *config.Config) {
	api := &api{
		cfg: cfg,
		loader: &loader{
			db: dbconn,
		},
	}

	router.POST(url.create, api.create)
	router.POST(url.login, api.login)
}

func (api *api) create(c *gin.Context) {
	var validate validateAccount
	if err := c.ShouldBindJSON(&validate); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "validate failed",
			"error":   err,
		})
		return
	}

	if _, err := api.loader.createAccount(validate.Username, validate.Password); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "create account error",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"error":   nil,
	})
	return
}

func (api *api) login(c *gin.Context) {
	var validate validateAccount
	if err := c.ShouldBindJSON(&validate); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "validate failed",
			"error":   err,
		})
		return
	}

	account, err := api.loader.getAccount(validate.Username, validate.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "get account error",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("%s:%s", account.username, hash.HS256(account.username)),
		"error":   nil,
	})
	return
}
