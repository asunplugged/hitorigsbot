package model

import (
	"database/sql"
	"errors"
	"time"
)

// LineBot -
type LineBot struct {
	ID          string    `db:"id" cc:"id"`
	Name        string    `db:"name" cc:"name"`
	AccessToken string    `db:"access_token" cc:"access_token"`
	Secret      string    `db:"secret" cc:"secret"`
	Ctime       time.Time `db:"ctime" cc:"ctime"`
	Mtime       time.Time `db:"mtime" cc:"mtime"`
}

// GetBotInfo -
func GetBotInfo(id string) (bot *LineBot, err error) {
	if len(id) == 0 {
		return nil, errors.New("id is emptu")
	}

	query := `select 
			"id", "name", "access_token", "secret", "ctime", "mtime"
		from public."line_bot"
		where
			"id" = $1`
	bot = &LineBot{}
	err = x.Get(bot, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return
}
