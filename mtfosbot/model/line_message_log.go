package model

import (
	"fmt"
	"time"
)

// LineMessageLog -
type LineMessageLog struct {
	ID      string    `db:"id" cc:"id"`
	Group   string    `db:"group" cc:"group"`
	User    string    `db:"user" cc:"user"`
	Message string    `db:"message" cc:"message"`
	Type    string    `db:"type" cc:"type"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"mtime"`
}

// LineMessageLogWithUG -
type LineMessageLogWithUG struct {
	LineMessageLog `cc:"-,<<"`
	GroupName      string `db:"group_name" cc:"group_name"`
	UserName       string `db:"user_name" cc:"user_name"`
}

// AddLineMessageLog -
func AddLineMessageLog(g, u, msg, t string) (msglog *LineMessageLog, err error) {
	query := `insert into "public"."line_message_log" ("group", "user", "message", "type") values ($1, $2, $3, $4)`
	msglog = &LineMessageLog{}
	err = x.Get(msglog, query, g, u, msg, t)
	return
}

// GetLineMessageLogCount -
func GetLineMessageLogCount(where ...interface{}) (c int, err error) {
	query := `select count(*) as c from "public"."line_message_log"`
	values := make([]interface{}, 0)
	if len(where) > 0 {
		if whereMap, ok := where[0].(map[string]string); ok {
			str := ""
			idx := 1
			for k, v := range whereMap {
				if len(str) > 0 {
					str += " and "
				}
				str += fmt.Sprintf(` "%s" = $%d`, k, idx)
				idx++
				values = append(values, v)
			}
			if len(str) > 0 {
				query += " where " + str
			}
		}
	}
	err = x.Get(&c, query, values...)
	return
}

// GetLineMessageLogList -
func GetLineMessageLogList(g, u string, offset, limit int, orderArg string) (logs []*LineMessageLogWithUG, err error) {
	if orderArg != "asc" && orderArg != "desc" {
		orderArg = "desc"
	}
	params := struct {
		Group string `db:"group"`
		User  string `db:"user"`
	}{}
	query := `select m.*, g.name as group_name, u.name as user_name from "public"."line_message_log" m
		left join "public"."line_user" u
		on u.id = m.user
		left join "public"."line_group" g
		on g.id = m.group
		`
	where := ""
	if len(g) > 0 {
		where = ` where g.id = :group`
		params.Group = g
	}
	if len(u) > 0 {
		if len(where) > 0 {
			where += ` and u.id = :user`
		} else {
			where += ` where u.id = :user`
		}
		params.User = u
	}
	order := fmt.Sprintf(`order by m.ctime %s`, orderArg)
	pager := fmt.Sprintf("offset %d limit %d", offset, limit)

	stmt, err := x.PrepareNamed(fmt.Sprintf("%s %s %s %s", query, where, order, pager))
	if err != nil {
		return nil, err
	}

	err = stmt.Select(&logs, params)
	return
}
