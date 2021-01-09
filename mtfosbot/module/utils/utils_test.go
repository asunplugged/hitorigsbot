package utils

import (
	"testing"
	"time"
)

// Commands - struct
type Commands struct {
	Cmd     string    `db:"cmd" cc:"cmd"`
	Message string    `db:"message" cc:"message"`
	Group   string    `db:"group" cc:"group"`
	Ctime   time.Time `db:"ctime" cc:"ctime"`
	Mtime   time.Time `db:"mtime" cc:"ctime"`
}

// CommandsWithGroup -
type CommandsWithGroup struct {
	Commands  `cc:"-,<<"`
	GroupName string `db:"group_name" cc:"group_name"`
}

func TestToMap(t *testing.T) {
	cmd := Commands{}
	cmd.Cmd = "test"
	cmd.Message = "test message"
	cmd.Group = ""
	cmd.Ctime = time.Now()
	cmd.Mtime = time.Now()

	cmdWGroup := CommandsWithGroup{}
	cmdWGroup.Commands = cmd
	cmdWGroup.GroupName = "asd"

	ToMap(cmdWGroup)
}

func TestCalcPage(t *testing.T) {
	page := CalcPage(10, 1, 10)
	if page.Page != 1 {
		t.Error("Page Calc fail")
	}
	if page.Total != 1 {
		t.Error("Page Calc fail")
	}
	if page.Limit != 10 {
		t.Error("limit calc fail")
	}
	if page.Offset != 0 {
		t.Error("offset calc fail")
	}
}

func BenchmarkCalcPage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalcPage(10000, 30, 10)
	}
}
