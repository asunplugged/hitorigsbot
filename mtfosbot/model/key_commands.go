package model

import (
	"errors"
	"time"
)

// KeyCommands - struct
type KeyCommands struct {
	Key     string    `db:"key" cc:"key"`
	Group   string    `db:"group" cc:"group"`
	Message string    `db:"message" cc:"message"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"ctime"`
}

// GetGroupKeyCommand -
func GetGroupKeyCommand(c, g string) (cmd *KeyCommands, err error) {
	if len(c) == 0 {
		return nil, errors.New("command is empty")
	}

	query := `select c.* from "public"."key_commands" c
		where c."key" = $1
		and (c."group" = '' or c."group" = $2)
		order by c."group" desc 
		limit 1`
	row := x.QueryRowx(query, c, g)
	// err = x.Get(&cmd, query, c, g)
	cmd = &KeyCommands{}
	err = row.StructScan(cmd)
	if err != nil {
		return nil, err
	}

	return
}
