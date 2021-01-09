package linemsg

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"time"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/line"
	"git.trj.tw/golang/mtfosbot/module/config"
	"git.trj.tw/golang/mtfosbot/module/es"
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
	msgcmd "git.trj.tw/golang/mtfosbot/module/message-command"
)

func messageType(botid string, e *lineobj.EventObject) {
	log.Println("proc msg type text :: ", botid)
	msg := e.Message
	mtype, ok := msg["type"]
	if !ok {
		return
	}

	if t, ok := mtype.(string); ok {
		switch t {
		case "text":
			textMsg(botid, e)
			break
		case "image":
			imageMsg(botid, e)
			break
		}
	}
	return
}

func textMsg(botid string, e *lineobj.EventObject) {
	msg := e.Message
	mtxt, ok := msg["text"]
	if !ok {
		return
	}

	// group action
	if e.Source.Type == "group" {
		if txt, ok := mtxt.(string); ok {
			msgcmd.ParseLineMsg(botid, txt, e.ReplyToken, e.Source)
			saveTextMsgToLog(botid, txt, e.Source)
		}
	}
	return
}

func imageMsg(botid string, e *lineobj.EventObject) {
	msg := e.Message
	imgID, ok := msg["id"]
	if !ok {
		return
	}
	// group action
	if e.Source.Type == "group" {
		if id, ok := imgID.(string); ok {
			saveImageMsgToLog(botid, id, e.Source)
		}
	}
}

func getSourceUser(accessToken, uid, gid string) (u *model.LineUser, err error) {
	userData, err := model.GetLineUserByID(uid)
	if err != nil {
		return
	}

	if userData == nil {
		tmpu, err := line.GetUserInfo(accessToken, uid, gid)
		if err != nil || tmpu == nil {
			return nil, err
		}
		userData = &model.LineUser{}
		userData.ID = tmpu.UserID
		userData.Name = tmpu.DisplayName
		err = userData.Add()
		if err != nil {
			return nil, err
		}
	} else {
		if userData.Mtime.Unix() < (time.Now().Unix() - 86400) {
			tmpu, err := line.GetUserInfo(accessToken, uid, gid)
			if err != nil || tmpu == nil {
				return nil, err
			}
			userData.Name = tmpu.DisplayName
			err = userData.UpdateName()
			if err != nil {
				return nil, err
			}
		}
	}

	return userData, nil
}

func saveTextMsgToLog(botid, txt string, s *lineobj.SourceObject) {
	bot, err := model.GetBotInfo(botid)
	if err != nil || bot == nil {
		fmt.Println("get bot info fail :: ", err)
		return
	}
	u, err := getSourceUser(bot.AccessToken, s.UserID, s.GroupID)
	if err != nil || u == nil {
		return
	}

	// go saveLineMessageLogToES(s.GroupID, s.UserID, txt, "text")
	model.AddLineMessageLog(s.GroupID, s.UserID, txt, "text")
}

func saveImageMsgToLog(botid, id string, s *lineobj.SourceObject) {
	bot, err := model.GetBotInfo(botid)
	if err != nil || bot == nil {
		fmt.Println("get bot info fail :: ", err)
		return
	}
	u, err := getSourceUser(bot.AccessToken, s.UserID, s.GroupID)
	if err != nil || u == nil {
		return
	}

	mime, err := line.GetContentHead(bot.AccessToken, id)
	if err != nil || len(mime) == 0 {
		return
	}

	ext := ""
	switch mime {
	case "image/jpeg":
		ext = ".jpg"
		break
	case "image/jpg":
		ext = ".jpg"
		break
	case "image/png":
		ext = ".png"
		break
	default:
		return
	}

	conf := config.GetConf()

	fname := fmt.Sprintf("log_%s%s", id, ext)

	fullPath := path.Join(conf.LogImageRoot, fname)

	w, err := os.Create(fullPath)
	if err != nil {
		return
	}
	defer w.Close()

	err = line.DownloadContent(bot.AccessToken, id, w)
	if err != nil {
		return
	}

	furl, err := url.Parse(conf.URL)
	if err == nil {
		furl, err = furl.Parse(fmt.Sprintf("/image/line_log_image/%s", fname))
		if err == nil {
			fname = furl.String()
		}
	}

	// go saveLineMessageLogToES(s.GroupID, s.UserID, fname, "image")
	model.AddLineMessageLog(s.GroupID, s.UserID, fname, "image")
}

func saveLineMessageLogToES(gid, uid, content, msgType string) {
	lineGroup, err := model.GetLineGroup(gid)
	if err != nil {
		log.Println("get line group error :: ", err)
		return
	}
	lineUser, err := model.GetLineUserByID(uid)
	if err != nil {
		log.Println("get line user error :: ", err)
		return
	}

	logMsg := make(map[string]interface{})

	logMsg["message"] = content
	logMsg["type"] = msgType
	logMsg["group"] = lineGroup.ID
	logMsg["group_name"] = lineGroup.Name
	logMsg["user"] = lineUser.ID
	logMsg["user_name"] = lineUser.Name

	err = es.PutLog("log", logMsg)
	if err != nil {
		log.Println("put log fail :: ", err)
	}
}
