package model

import (
	"database/sql"
	"errors"
	"time"
)

// TwitchGroup -
type TwitchGroup struct {
	*LineGroup
	Tmpl string `db:"tmpl"`
}

// TwitchChannel - struct
type TwitchChannel struct {
	ID         string         `db:"id" cc:"id"`
	Name       string         `db:"name" cc:"name"`
	LastStream string         `db:"laststream" cc:"laststream"`
	Join       bool           `db:"join" cc:"join"`
	OpayID     string         `db:"opayid" cc:"opayid"`
	Ctime      time.Time      `db:"ctime" cc:"ctime"`
	Mtime      time.Time      `db:"mtime" cc:"ctime"`
	Groups     []*TwitchGroup `db:"-" cc:"-"`
}

// GetAllTwitchChannel -
func GetAllTwitchChannel() (channels []*TwitchChannel, err error) {
	err = x.Select(&channels, `select * from "public"."twitch_channel" order by "name" desc`)
	return
}

// GetJoinChatChannel -
func GetJoinChatChannel() (channels []*TwitchChannel, err error) {
	err = x.Select(&channels, `select * from "public"."twitch_channel" where "join" = true`)
	return
}

// GetTwitchChannelWithName -
func GetTwitchChannelWithName(name string) (ch *TwitchChannel, err error) {
	if len(name) == 0 {
		return nil, errors.New("name empty")
	}
	ch = &TwitchChannel{}
	err = x.Get(ch, `select * from "public"."twitch_channel" where "name" = $1`, name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// GetTwitchChannelWithID -
func GetTwitchChannelWithID(id string) (ch *TwitchChannel, err error) {
	if len(id) == 0 {
		return nil, errors.New("id empty")
	}
	ch = &TwitchChannel{}
	err = x.Get(ch, `select * from "public"."twitch_channel" where "id" = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// GetWithName -
func (p *TwitchChannel) GetWithName() (err error) {
	stmt, err := x.PrepareNamed(`select * from "public"."twitch_channel" where "name" = :name`)
	if err != nil {
		return err
	}

	err = stmt.Get(p, p)
	return
}

// Add -
func (p *TwitchChannel) Add() (err error) {
	stmt, err := x.PrepareNamed(`insert into "public"."twitch_channel" ("id", "name", "laststream", "join", "opayid") values (:id, :name, :laststream, :join, :opayid) returning *`)
	if err != nil {
		return err
	}
	err = stmt.Get(p, p)
	return
}

// UpdateStream -
func (p *TwitchChannel) UpdateStream(streamID string) (err error) {
	query := `update "public"."twitch_channel" set "laststream" = $1, "mtime" = now() where "id" = $2`
	_, err = x.Exec(query, streamID, p.ID)
	if err != nil {
		return
	}
	p.LastStream = streamID
	return
}

// UpdateName -
func (p *TwitchChannel) UpdateName(name string) (err error) {
	_, err = x.Exec(`update "public"."twitch_channel" set "name" = $1, "mtime" = now() where "id" = $2`, name, p.ID)
	if err != nil {
		return
	}
	p.Name = name
	return
}

// UpdateJoin -
func (p *TwitchChannel) UpdateJoin(join bool) (err error) {
	p.Join = join
	_, err = x.NamedExec(`update "public"."twitch_channel" set "join" = :join, "mtime" = now() where "id" = :id`, p)
	return
}

// UpdateOpayID -
func (p *TwitchChannel) UpdateOpayID(id string) (err error) {
	p.OpayID = id
	_, err = x.NamedExec(`update "public"."twitch_channel" set "opayid" = :opayid, "mtime" = now() where "id" = :id`, p)
	return
}

// GetGroups -
func (p *TwitchChannel) GetGroups() (err error) {
	query := `select g.*, rt.tmpl as tmpl from "public"."twitch_channel" tw
		left join "public"."line_twitch_rt" rt
		on rt.twitch = tw.id
		left join "public"."line_group" g
		on g.id = rt.line
		where 
		tw.id = $1
		and rt.twitch is not null`
	err = x.Select(&p.Groups, query, p.ID)
	return
}
