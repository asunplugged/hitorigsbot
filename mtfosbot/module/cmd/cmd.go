package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"/golang/mtfosbot/module/schema"

	"/golang/mtfosbot/model"
)

// DBTool - deploy database schemas
func DBTool() {
	db := model.GetDB()
	if db == nil {
		log.Fatal(errors.New("database object is nil"))
	}

	dbver, err := schema.ReadVersions()
	if err != nil {
		log.Fatal(err)
	}

	version := -1

	vcExists := false
	// check version_ctrl table exists
	row := db.QueryRowx(`select exists(select 1 from "information_schema"."tables" where "table_schema" = $1 and "table_name" = $2) as exists`, "public", "version_ctrl")
	err = row.Scan(&vcExists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	// read max version
	if vcExists {
		row := db.QueryRowx(`select max(version) as version from "public"."version_ctrl"`)
		err := row.Scan(&version)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}

	fmt.Println("Database Schema Version is ", version)
	var vers []schema.VersionInfo
	for _, v := range dbver.Versions {
		if v.Version > version {
			vers = append(vers, v)
		}
	}
	if len(vers) == 0 {
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range vers {
		fmt.Printf("Run Version: %d, FileName: %s\n", v.Version, v.File)
		query, err := schema.ReadSchema(v.File)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		// run schema file
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		// insert version ctrl
		_, err = tx.Exec(`insert into "public"."version_ctrl" ("version", "ctime", "querystr") values ($1, now(), $2)`, v.Version, query)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	tx.Commit()
}
