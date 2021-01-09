package model

import (
	"database/sql"
	"time"
)

// LineUser -
type LineUser struct {
	ID    string    `db:"id" cc:"id"`
	Name  string    `db:"name" cc:"name"`
	Ctime time.Time `db:"ctime" cc:"ctime"`
	Mtime time.Time `db:"mtime" cc:"mtime"`
}

// GetLineUserByID -
func GetLineUserByID(id string) (u *LineUser, err error) {
	u = &LineUser{}
	err = x.Get(u, `select * from "public"."line_user" where "id" = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// Add -
func (p *LineUser) Add() (err error) {
	_, err = x.NamedExec(`insert into "public"."line_user" ("id", "name") values (:id, :name)`, p)
	return
}

// UpdateName -
func (p *LineUser) UpdateName() (err error) {
	_, err = x.NamedExec(`update "public"."line_user" set "name" = :name, "mtime" = now() where "id" = :id`, p)
	return
}
