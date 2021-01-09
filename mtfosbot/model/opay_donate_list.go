package model

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// OpayDonateList - struct
type OpayDonateList struct {
	OpayID   string    `db:"opayid" cc:"opayid"`
	DonateID string    `db:"donate_id" cc:"donate_id"`
	Price    int       `db:"price" cc:"price"`
	Text     string    `db:"text" cc:"text"`
	Name     string    `db:"name" cc:"name"`
	Ctime    time.Time `db:"ctime" cc:"ctime"`
	Mtime    time.Time `db:"mtime" cc:"ctime"`
}

// SumChannelDonatePriceSinceTime -
func SumChannelDonatePriceSinceTime(twitch string, t time.Time) (i int, err error) {
	query := `select sum(list.price) as total from "public"."opay_donate_list" list 
	left join "public"."twitch_channel" tw
	on tw.opayid = list.opayid
	where
	tw.opayid != ''
	and tw.id = $1
	and list.ctime >= $2
	group by list.opayid`
	row := x.QueryRowx(query, twitch, t)
	err = row.Scan(&i)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return
}

// GetDonateListWithIDs -
func GetDonateListWithIDs(ids []string) (ls []*OpayDonateList, err error) {
	if len(ids) == 0 {
		return
	}
	s, i, err := sqlx.In(`select * from "public"."opay_donate_list" where "donate_id" in (?)`, ids)
	if err != nil {
		return
	}
	s = x.Rebind(s)
	err = x.Select(&ls, s, i...)
	return
}

// InsertData -
func (p *OpayDonateList) InsertData() (err error) {
	_, err = x.NamedExec(`insert into "public"."opay_donate_list" ("opayid", "donate_id", "price", "text", "name") values (:opayid, :donate_id, :price, :text, :name)`, p)
	return
}
