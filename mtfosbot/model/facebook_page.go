package model

import (
	"database/sql"
	"time"
)

// FBGroup -
type FBGroup struct {
	*LineGroup
	Tmpl string `db:"tmpl"`
}

// FacebookPage - struct
type FacebookPage struct {
	ID       string     `db:"id" cc:"id"`
	LastPost string     `db:"lastpost" cc:"lastpost"`
	Ctime    time.Time  `db:"ctime" cc:"ctime"`
	Mtime    time.Time  `db:"mtime" cc:"ctime"`
	Groups   []*FBGroup `db:"-"`
}

// GetAllFacebookPage -
func GetAllFacebookPage() (pages []*FacebookPage, err error) {
	err = x.Select(&pages, `select * from "public"."facebook_page"`)
	if err != nil {
		return nil, err
	}
	return
}

// GetFacebookPage -
func GetFacebookPage(id string) (page *FacebookPage, err error) {
	page = &FacebookPage{}
	err = x.Get(page, `select * from "public"."facebook_page" where "id" = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}

// AddPage -
func (p *FacebookPage) AddPage() (err error) {
	stmt, err := x.PrepareNamed(`insert into "public"."facebook_page" ("id", "lastpost") values (:id, :lastpost) returning *`)
	if err != nil {
		return err
	}
	err = stmt.Get(p, p)
	return
}

// UpdatePost -
func (p *FacebookPage) UpdatePost(postID string) (err error) {
	query := `update "public"."facebook_page" set "lastpost" = $1, "mtime" = now() where id = $2`
	_, err = x.Exec(query, postID, p.ID)
	if err != nil {
		return
	}
	p.LastPost = postID
	return
}

// GetGroups -
func (p *FacebookPage) GetGroups() (err error) {
	query := `select g.*, rt.tmpl as tmpl from "public"."facebook_page" p
		left join "public"."line_fb_rt" rt
		on rt."facebook" = p.id
		left join "public"."line_group" g
		on g.id = rt."line"
		where 
		p.id = $1
		and rt.facebook is not null`
	err = x.Select(&p.Groups, query, p.ID)
	return
}
