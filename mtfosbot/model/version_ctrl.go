package model

import "time"

// VersionCTRL -
type VersionCTRL struct {
	Version  int       `db:"version" cc:"version"`
	Ctime    time.Time `db:"ctime" cc:"ctime"`
	QueryStr string    `db:"querystr" cc:"querystr"`
}

// SaveVersionLog -
func SaveVersionLog(version int, qstr string) error {
	query := `insert into "public"."version_ctrl" ("version", "querystr") values (?, ?)`
	_, err := x.Exec(query, version, qstr)
	return err
}
