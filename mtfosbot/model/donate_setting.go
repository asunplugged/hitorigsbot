package model

import (
	"database/sql"
	"time"
)

// DonateSetting - struct
type DonateSetting struct {
	Twitch       string    `db:"twitch" cc:"twitch"`
	StartDate    time.Time `db:"start_date" cc:"start_date"`
	EndDate      time.Time `db:"end_date" cc:"end_date"`
	TargetAmount int       `db:"target_amount" cc:"target_amount"`
	StartAmount  int       `db:"start_amount" cc:"start_amount"`
	Title        string    `db:"title" cc:"title"`
	Type         string    `db:"type" cc:"type"`
	Ctime        time.Time `db:"ctime" cc:"ctime"`
	Mtime        time.Time `db:"mtime" cc:"mtime"`
}

// GetDonateSettingByChannel -
func GetDonateSettingByChannel(id string) (ds *DonateSetting, err error) {
	query := `select ds.*
		from "public"."donate_setting" ds
		left join "public"."twitch_channel" ch
		on ch.id = ds.twitch 
		where 
		ch.id = $1`
	row := x.QueryRowx(query, id)
	ds = &DonateSetting{}
	err = row.StructScan(ds)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// InsertOrUpdate -
func (p *DonateSetting) InsertOrUpdate() (err error) {
	query := `insert into "public"."donate_setting" 
	("twitch", "start_date", "end_date", "target_amount", "title", "start_amount") values
	(:twitch, now(), :end_date, :target_amount, :title, :start_amount)
	on CONFLICT ("twitch") DO UPDATE 
	set 
	"start_date" = now(),
	"end_date" = :end_date,
	"target_amount" = :target_amount,
	"title" = :title,
	"start_amount" = :start_amount, "mtime" = now()`
	_, err = x.NamedExec(query, p)
	return
}
