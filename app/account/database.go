package account

import (
	"database/sql"
	"time"

	"github.com/killua4564/go-note/utils/hash"
)

type loader struct {
	db *sql.DB
}

func (l *loader) createAccount(username string, password string) (int64, error) {
	password = hash.PBKDF2(password)
	create_time := time.Now().UnixNano() / 1e6
	var args []interface{} = []interface{}{
		username, password, create_time,
	}

	query := "INSERT INTO `account` (`username`, `password`, `create_time`) VALUES (?, ?, ?);"
	result, err := l.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (l *loader) getAccount(username string, password string) (*Account, error) {
	password = hash.PBKDF2(password)
	var args []interface{} = []interface{}{
		username, password,
	}

	query := "SELECT `username` FROM `account` WHERE username=? AND password=?;"
	row := l.db.QueryRow(query, args...)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var account Account
	if err := row.Scan(&account.Username); err != nil {
		return nil, err
	}

	return &account, nil
}
