package model

import (
	"database/sql"
	"time"
)

// Lottery -
type Lottery struct {
	ID      string    `db:"id" cc:"id"`
	Type    string    `db:"type" cc:"type"`
	Message string    `db:"message" cc:"message"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"mtime"`
}

// GetRandomLotteryByType -
func GetRandomLotteryByType(t string) (p *Lottery, err error) {
	parr, err := GetRandomLotteryByTypeAndLimit(t, 1)
	if err != nil {
		return nil, err
	}
	if parr == nil {
		return nil, nil
	}
	return parr[0], err
}

// GetRandomLotteryByTypeAndLimit -
func GetRandomLotteryByTypeAndLimit(t string, limit int) (p []*Lottery, err error) {
	err = x.Select(&p, `select * from "public"."lottery" where "type" = $1 offset floor(random() * (select count(*) as c from "public"."lottery" where "type" = $2)) limit $3`, t, t, limit)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
