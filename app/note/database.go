package note

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNotOwner = errors.New("account is not note's owner")

type SessionRunner interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type loader struct {
	runner SessionRunner
}

func (l *loader) BeginTransaction() (*loader, error) {
	if tx, ok := l.runner.(*sql.Tx); ok {
		return &loader{runner: tx}, nil
	}
	if _, ok := l.runner.(*sql.DB); !ok {
		return nil, errors.New("Invalid SessionRunner")
	}
	tx, err := l.runner.(*sql.DB).Begin()
	if err != nil {
		return nil, err
	}

	return &loader{runner: tx}, nil
}

func (l *loader) getAccountID(username string) (int64, error) {
	query := "SELECT `id` FROM `account` WHERE username=?"
	row := l.runner.QueryRow(query, username)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (l *loader) getNoteID(SID string) (int64, error) {
	query := "SELECT `id` FROM `note` WHERE sid=?"
	row := l.runner.QueryRow(query, SID)
	if row.Err() != nil {
		return 0, row.Err()
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (l *loader) createNote(account *Account, note *Note) (int64, error) {
	create_time := time.Now().UnixNano() / 1e6
	update_time := create_time

	txl, err := l.BeginTransaction()
	if err != nil {
		return 0, err
	}
	defer txl.runner.(*sql.Tx).Rollback()

	accountID, err := txl.getAccountID(account.Username)
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO `note` (`sid`, `topic`, `content`, `create_time`, `update_time`) VALUES (?, ?, ?, ?, ?);"
	result, err := txl.runner.Exec(query, note.SID, note.Topic, note.Content, create_time, update_time)
	if err != nil {
		return 0, err
	}

	noteID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if _, err = txl.createAccountNote(&AccountNote{
		AccountID: accountID,
		NoteID:    noteID,
		IsOwner:   int64(1),
	}); err != nil {
		return 0, err
	}

	if err = txl.runner.(*sql.Tx).Commit(); err != nil {
		return 0, err
	}

	return noteID, nil
}

func (l *loader) getNote(account *Account, SID string) (*Note, error) {
	query := "SELECT `note`.`sid`, `note`.`topic`, `note`.`content` FROM `account_note` " +
		"INNER JOIN `account` ON `account`.`id`=`account_note`.`account_id` " +
		"INNER JOIN `note` ON `note`.`id`=`account_note`.`note_id` " +
		"WHERE `account`.`username`=? AND `note`.`sid`=?;"
	row := l.runner.QueryRow(query, account.Username, SID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var note Note
	if err := row.Scan(&note.SID, &note.Topic, &note.Content); err != nil {
		return nil, err
	}

	return &note, nil
}

func (l *loader) getNoteList(account *Account) (*NoteList, error) {
	query := "SELECT `note`.`sid`, `note`.`topic`, `account_note`.`is_owner` FROM `account_note` " +
		"INNER JOIN `account` ON `account`.`id`=`account_note`.`account_id` " +
		"INNER JOIN `note` ON `note`.`id`=`account_note`.`note_id` " +
		"WHERE `account`.`username`=? ORDER BY `account_note`.`is_owner`;"
	rows, err := l.runner.Query(query, account.Username)
	if err != nil {
		return nil, err
	}

	var noteList NoteList
	var is_owner int64
	for rows.Next() {
		var note Note
		if err = rows.Scan(&note.SID, &note.Topic, &is_owner); err != nil {
			return nil, err
		}
		if is_owner == int64(1) {
			noteList.Owner = append(noteList.Owner, note)
		} else {
			noteList.Viewer = append(noteList.Viewer, note)
		}
	}

	return &noteList, nil
}

func (l *loader) isNoteOwner(account *Account, SID string) (bool, error) {
	query := "SELECT `account_note`.`is_owner` FROM `account_note` " +
		"INNER JOIN `account` ON `account`.`id`=`account_note`.`account_id` " +
		"INNER JOIN `note` ON `note`.`id`=`account_note`.`note_id` " +
		"WHERE `account`.`username`=? AND `note`.`sid`=? AND `account_note`.`is_owner`=1;"
	row := l.runner.QueryRow(query, account.Username, SID)
	if row.Err() != nil {
		return false, row.Err()
	}

	var is_owner int64
	err := row.Scan(&is_owner)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (l *loader) updateNote(account *Account, note *Note) (int64, error) {
	update_time := time.Now().UnixNano() / 1e6

	txl, err := l.BeginTransaction()
	if err != nil {
		return 0, err
	}
	defer txl.runner.(*sql.Tx).Rollback()

	is_owner, err := txl.isNoteOwner(account, note.SID)
	if err != nil {
		return 0, err
	}
	if !is_owner {
		return 0, ErrNotOwner
	}

	query := "UPDATE `note` SET `topic`=?, `content`=?, `update_time`=? WHERE `sid`=?;"
	result, err := txl.runner.Exec(query, note.Topic, note.Content, update_time, note.SID)
	if err != nil {
		return 0, err
	}

	updated, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if err = txl.runner.(*sql.Tx).Commit(); err != nil {
		return 0, err
	}

	return updated, nil
}

func (l *loader) removeNote(account *Account, SID string) (int64, error) {
	txl, err := l.BeginTransaction()
	if err != nil {
		return 0, err
	}
	defer txl.runner.(*sql.Tx).Rollback()

	noteID, err := txl.getNoteID(SID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if err == sql.ErrNoRows {
		return 0, ErrNotOwner
	}

	is_owner, err := txl.isNoteOwner(account, SID)
	if err != nil {
		return 0, err
	}
	if !is_owner {
		return 0, ErrNotOwner
	}

	if _, err = txl.removeAccountNote(noteID); err != nil {
		return 0, err
	}

	query := "DELETE FROM `note` WHERE `sid`=?;"
	result, err := txl.runner.Exec(query, SID)
	if err != nil {
		return 0, err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if err = txl.runner.(*sql.Tx).Commit(); err != nil {
		return 0, err
	}

	return deleted, nil
}

func (l *loader) createAccountNote(accountNote *AccountNote) (int64, error) {
	create_time := time.Now().UnixNano() / 1e6
	update_time := create_time

	query := "INSERT INTO `account_note` (`account_id`, `note_id`, `is_owner`, `create_time`, `update_time`) VALUES (?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE `is_owner`=?, `update_time`=?;"
	result, err := l.runner.Exec(query, accountNote.AccountID, accountNote.NoteID, accountNote.IsOwner, create_time, update_time,
		accountNote.IsOwner, update_time)
	if err != nil {
		return 0, err
	}

	accountNoteID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return accountNoteID, nil
}

func (l *loader) removeAccountNote(noteID int64) (int64, error) {
	query := "DELETE FROM `account_note` WHERE `note_id`=?;"
	result, err := l.runner.Exec(query, noteID)
	if err != nil {
		return 0, err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return deleted, err
}
