package model

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// Commands - struct
type Commands struct {
	Cmd     string    `db:"cmd" cc:"cmd"`
	Message string    `db:"message" cc:"message"`
	Group   string    `db:"group" cc:"group"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"mtime"`
}

// CommandsWithGroup -
type CommandsWithGroup struct {
	Commands  `cc:"-,<<"`
	GroupName string `db:"group_name" cc:"group_name"`
}

// GetAllCommands -
func GetAllCommands() (cmds []*Commands, err error) {
	err = x.Select(&cmds, `select * from "public"."commands"`)
	return
}

// GetGroupCommand -
func GetGroupCommand(c, g string) (cmd *Commands, err error) {
	if len(c) == 0 {
		return nil, errors.New("command is empty")
	}
	query := `select c.* from "public"."commands" c
		where c."cmd" = $1
		and (c."group" = '' or c."group" = $2)
		order by c."group" desc
		limit 1`
	row := x.QueryRowx(query, c, g)
	// err = x.Get(&cmd, query, c, g)
	cmd = &Commands{}
	err = row.StructScan(cmd)
	if err != nil {
		return nil, err
	}
	return
}

// GetCommandCount -
func GetCommandCount(where ...interface{}) (c int, err error) {
	query := `select count(*) as c from "public"."commands"`
	values := make([]interface{}, 0)
	if len(where) > 0 {
		if whereMap, ok := where[0].(map[string]string); ok {
			str := ""
			idx := 1
			for k, v := range whereMap {
				if len(str) > 0 {
					str += " and "
				}
				str += fmt.Sprintf(` "%s" = $%d `, k, idx)
				idx++
				values = append(values, v)
			}
			if len(str) > 0 {
				query += ` where ` + str
			}
		}
	}
	err = x.Get(&c, query, values...)
	return
}

// GetCommands -
func GetCommands(where map[string]string, offset, limit int, order map[string]string) (cmds []*CommandsWithGroup, err error) {
	query := `select c.*, (case when g.name is null then '' else g.name end) as group_name from "public"."commands" c
		left join "public"."line_group" g
		on g.id = c.group `
	values := make([]interface{}, (len(where) + len(order)))
	idx := 1

	if len(where) > 0 {
		str := ""
		for k, v := range where {
			if len(str) > 0 {
				str += " and "
			}
			str += fmt.Sprintf(` "%s" = $%d `, k, idx)
			idx++
			values = append(values, v)
		}
		if len(str) > 0 {
			query += ` where ` + str
		}
	}
	if offset >= 0 {
		query += fmt.Sprintf(" offset %d ", offset)
	}
	if limit > 0 {
		query += fmt.Sprintf(" limit %d ", limit)
	}

	if len(order) > 0 {
		regex, err := regexp.Compile("(?i)(desc|asc)")
		if err != nil {
			return nil, err
		}
		str := ""
		for k, v := range order {
			if !regex.Match([]byte(v)) {
				continue
			}
			if len(str) > 0 {
				str += " , "
			}
			str += fmt.Sprintf(` "%s" %s `, k, v)
		}
	}

	err = x.Select(&cmds, query, values...)

	return
}

// AddCommand -
func AddCommand(cmdkey, message, group string) (cmd *Commands, err error) {
	if len(cmdkey) == 0 || len(message) == 0 {
		return nil, errors.New("cmd or message is empty")
	}
	query := `insert into "public"."commands" ("cmd", "message", "group") values ($1, $2, $3) returning *`
	cmd = &Commands{}
	err = x.Get(cmd, query, cmdkey, message, group)
	if err != nil {
		return nil, err
	}
	return
}

// CheckCommand -
func CheckCommand(cmd, group string) (exist bool, err error) {
	if len(cmd) == 0 {
		return false, errors.New("cmd is empty")
	}
	query := `select count(*) as c from "public"."commands" where "cmd" = $1 and "group" = $2`
	c := 0

	err = x.Get(&c, query, cmd, group)
	if err != nil {
		return false, err
	}

	return c > 0, nil
}

// DeleteCommand -
func DeleteCommand(cmd, group string) (err error) {
	if len(cmd) == 0 {
		return errors.New("cmd is empty")
	}
	query := `delete from "public"."commands" where "cmd" = $1 and "group" = $2`
	_, err = x.Exec(query, cmd, group)
	return
}

// UpdateCommand -
func UpdateCommand(cmd, group, message string) (err error) {
	if len(cmd) == 0 || len(message) == 0 {
		return errors.New("cmd or message is empty")
	}

	query := `update "public"."commands" set "message" = $1, "mtime" = now() where "cmd" = $2 and "group" = $3`
	_, err = x.Exec(query, message, cmd, group)
	return
}
