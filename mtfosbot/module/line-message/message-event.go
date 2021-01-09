package linemsg

import (
	"fmt"
	"log"

	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
)

// MessageEvent -
func MessageEvent(botid string, e *lineobj.EventObject) {
	log.Println("proc message evt :: ", botid, *e)
	if len(botid) == 0 {
		return
	}
	switch e.Type {
	case "message":
		messageType(botid, e)
		break
	default:
		fmt.Println("line webhook type not match")
	}
}
