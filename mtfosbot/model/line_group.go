package model

import (
	"database/sql"
	"errors"
	"time"
)

// LineGroup - struct
type LineGroup struct {
	ID     string         `db:"id" cc:"id"`
	Name   string         `db:"name" cc:"name"`
	Notify bool           `db:"notify" cc:"notify"`
	Owner  string         `db:"owner" cc:"owner"`
	Ctime  time.Time      `db:"ctime" cc:"ctime"`
	Mtime  time.Time      `db:"mtime" cc:"ctime"`
	BotID  sql.NullString `db:"bot_id" cc:"bot_id"`
}

// CheckGroup -
func CheckGroup(g string) (exists bool, err error) {
	ss := struct {
		C int `db:"c"`
	}{}

	err = x.Get(&ss, `select count(*) as c from "public"."line_group" where "id" = $1`, g)
	if err != nil {
		return false, err
	}
	return ss.C > 0, nil
}

// CheckGroupOwner -
func CheckGroupOwner(user, g string) (exists bool, err error) {
	ss := struct {
		C int `db:"c"`
	}{}

	err = x.Get(&ss, `select count(*) as c from "public"."line_group" where "id" = $1 and "owner" = $2`, g, user)
	if err != nil {
		return false, err
	}
	return ss.C > 0, nil
}

// GetLineGroup -
func GetLineGroup(id string) (g *LineGroup, err error) {
	g = &LineGroup{}
	err = x.Get(g, `select * from "public"."line_group" where "id" = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// GetLineGroupList -
func GetLineGroupList() (ls []*LineGroup, err error) {
	err = x.Select(&ls, `select * from "public"."line_group" order by "name"`)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// AddLineGroup -
func AddLineGroup(name, owner string, notify bool) (g *LineGroup, err error) {
	g = &LineGroup{}
	err = x.Get(g, `insert into "public"."line_group" ("name", "owner", "notify") values ($1, $2, $3)`, name, owner, notify)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// DeleteGroup -
func (p *LineGroup) DeleteGroup() (err error) {
	_, err = x.Exec(`delete from "public"."line_group" where "id" = $1`, p.ID)
	return
}

// GetBot - get group binding bot
func (p *LineGroup) GetBot() (bot *LineBot, err error) {
	id, err := p.BotID.Value()
	if err != nil {
		return nil, err
	}
	var botid string
	var ok bool
	if botid, ok = id.(string); !ok {
		return nil, errors.New("botid get fail")
	}
	bot, err = GetBotInfo(botid)
	return
}
