package model

import (
	"database/sql"
	"time"
)

type IGGroup struct {
	*LineGroup
	Tmpl string `db:"tmpl"`
}

//Instagram -
type Instagram struct {
	ID       string     `db:"id" cc:"id"`
	LastPost string     `db:"lastpost" cc:"lastpost"`
	Ctime    time.Time  `db:"ctime" cc:"ctime"`
	Mtime    time.Time  `db:"mtime" cc:"ctime"`
	Groups   []*IGGroup `db:"-"`
}

// GetAllInstagram -
func GetAllInstagram() (igs []*Instagram, err error) {
	err = x.Select(&igs, `select * from "public"."instagram"`)
	if err != nil {
		return nil, err
	}
	return
}

// GetInstagram -
func GetInstagram(id string) (ig *Instagram, err error) {
	ig = &Instagram{}
	err = x.Get(ig, `select * from "public"."instagram" where "id" = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// AddIG -
func (p *Instagram) AddIG() (err error) {
	stmt, err := x.PrepareNamed(`insert into "public"."instagram" ("id", "lastpost") values (:id, :lastpost) returning *`)
	if err != nil {
		return err
	}
	err = stmt.Get(p, p)
	return
}

// UpdatePost -
func (p *Instagram) UpdatePost(postID string) (err error) {
	query := `update "public"."instagram" set "lastpost" = $1, "mtime" = now() where "id" = $2`
	_, err = x.Exec(query, postID, p.ID)
	if err != nil {
		return
	}
	p.LastPost = postID
	return
}

// GetGroups -
func (p *Instagram) GetGroups() (err error) {
	query := `select g.*, rt.tmpl as tmpl from "public"."instagram" p
	left join "public"."line_ig_rt" rt
	on rt."ig" = p.id
	left join "public"."line_group" g
	on g.id = rt."line"
	where 
	p.id = $1
	and rt.ig is not null`
	err = x.Select(&p.Groups, query, p.ID)
	return
}
