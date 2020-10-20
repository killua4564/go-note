package note

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/killua4564/go-note/config"
	"github.com/killua4564/go-note/utils/hash"
)

type api struct {
	account *Account
	cfg     *config.Config
	loader  *loader
}

var url = struct {
	list   string
	view   string
	create string
	update string
	remove string
	viewer string
}{
	list:   "/note",
	view:   "/note/:sid",
	create: "/note",
	update: "/note/:sid",
	remove: "/note/:sid",
	viewer: "/note/:sid",
}

func NoteService(router *gin.RouterGroup, dbconn *sql.DB, cfg *config.Config) {
	api := &api{
		cfg: cfg,
		loader: &loader{
			runner: dbconn,
		},
	}

	router = router.Group("", func(c *gin.Context) {
		api.account = &Account{}

		auth := strings.Split(c.GetHeader("authorization"), " ")
		if len(auth) < 2 || auth[0] != "note-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   errors.New("Invalid note-token"),
				"message": "token error",
			})
			return
		}

		token := strings.Split(auth[1], ":")
		if len(token) < 2 || token[1] != hash.HS256(token[0]) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
				"error":   errors.New("Invalid note-token"),
				"message": "token error",
			})
			return
		}

		api.account = &Account{
			Username: token[0],
		}
	})

	router.GET(url.list, api.list)
	router.GET(url.view, api.view)
	router.POST(url.create, api.create)
	router.POST(url.viewer, api.viewer)
	router.PUT(url.update, api.update)
	router.DELETE(url.remove, api.remove)
}

func (api *api) list(c *gin.Context) {
	noteList, err := api.loader.getNoteList(api.account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "get note list error",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"error":   nil,
		"message": noteList,
	})
}

func (api *api) view(c *gin.Context) {
	var validate validateNoteUri
	if err := c.ShouldBindUri(&validate); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "validate failed",
		})
		return
	}

	note, err := api.loader.getNote(api.account, validate.SID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "get note error",
		})
		return
	}
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   nil,
			"message": "note not found or not viewer",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"error":   nil,
		"message": note,
	})
}

func (api *api) create(c *gin.Context) {
	var validate validateNote
	if err := c.ShouldBindJSON(&validate); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "validate failed",
		})
		return
	}

	sid := uuid.New().String()
	if _, err := api.loader.createNote(api.account, &Note{
		SID:     sid,
		Topic:   validate.Topic,
		Content: validate.Content,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "create note error",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"error":   nil,
		"message": sid,
	})
	return
}

func (api *api) viewer(c *gin.Context) {
	var validateUri validateNoteUri
	if err := c.ShouldBindUri(&validateUri); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "validate failed",
		})
		return
	}

	var validate validateAccountNote
	if err := c.ShouldBindJSON(&validate); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "validate failed",
		})
		return
	}

	accountID, err := api.loader.getAccountID(validate.Username)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "get account error",
		})
		return
	}
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   nil,
			"message": "account not found",
		})
		return
	}

	noteID, err := api.loader.getNoteID(validateUri.SID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "get note error",
		})
		return
	}
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   nil,
			"message": "note not found",
		})
		return
	}

	is_owner, err := api.loader.isNoteOwner(api.account, validateUri.SID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "get account note error",
		})
		return
	}
	if !is_owner {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   nil,
			"message": "note not owner",
		})
		return
	}

	if _, err = api.loader.createAccountNote(&AccountNote{
		AccountID: accountID,
		NoteID:    noteID,
		IsOwner:   *validate.IsOwner,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "create account note error",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"error":   nil,
		"message": "success",
	})
	return
}

func (api *api) update(c *gin.Context) {
	var validateUri validateNoteUri
	if err := c.ShouldBindUri(&validateUri); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "validate failed",
		})
		return
	}

	var validate validateNote
	if err := c.ShouldBindJSON(&validate); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "validate failed",
		})
		return
	}

	updated, err := api.loader.updateNote(api.account, &Note{
		SID:     validateUri.SID,
		Topic:   validate.Topic,
		Content: validate.Content,
	})
	if err != nil && err != ErrNotOwner {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "update note error",
		})
		return
	}
	if updated <= 0 || err == ErrNotOwner {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   nil,
			"message": "note not found or not owner",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"error":   nil,
		"message": "success",
	})
}

func (api *api) remove(c *gin.Context) {
	var validate validateNoteUri
	if err := c.ShouldBindUri(&validate); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   err,
			"message": "validate failed",
		})
		return
	}

	deleted, err := api.loader.removeNote(api.account, validate.SID)
	if err != nil && err != ErrNotOwner {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   err,
			"message": "delete note error",
		})
		return
	}
	if deleted <= 0 || err == ErrNotOwner {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   nil,
			"message": "note not found or not owner",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"error":   nil,
		"message": "success",
	})
}
