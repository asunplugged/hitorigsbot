package background

import (
	"fmt"
	"log"
	"strings"
	"time"

	"/golang/mtfosbot/model"
	"/golang/mtfosbot/module/apis/line"
	"/golang/mtfosbot/module/apis/twitch"
)

func getStreamStatus() {
	channels, err := model.GetAllTwitchChannel()
	if err != nil {
		return
	}
	var ids []string
	for _, v := range channels {
		ids = append(ids, v.ID)
	}

	info := twitch.GetUserStreamStatus(ids)
	if len(info) == 0 {
		return
	}
	for _, v := range info {
		for _, ch := range channels {
			if v.UserID == ch.ID {
				go checkStream(ch, v)
			}
		}
	}
}

func checkStream(ch *model.TwitchChannel, info *twitch.StreamInfo) {
	if info.ID == ch.LastStream {
		return
	}
	err := ch.GetGroups()
	if err != nil {
		return
	}
	err = ch.UpdateStream(info.ID)
	if err != nil {
		return
	}

	// 開台間隔小於10分鐘不通知
	if time.Now().Unix()-ch.Mtime.Unix() < 600 {
		return
	}

	chData := twitch.GetUserDataByID(ch.ID)
	if chData != nil {
		if chData.Login != ch.Name {
			ch.UpdateName(chData.Login)
		}
	}

	link := fmt.Sprintf("https://twitch.tv/%s", ch.Name)
	for _, v := range ch.Groups {
		if v.Notify {
			bot, err := v.GetBot()
			if err != nil || bot == nil {
				log.Println("get group binding bot fail :: ", err)
				continue
			}
			tmpl := v.Tmpl
			if len(tmpl) > 0 {
				tmpl = strings.Replace(tmpl, "{txt}", info.Title, -1)
				tmpl = strings.Replace(tmpl, "{link}", link, -1)
			} else {
				tmpl = fmt.Sprintf("%s\n%s", info.Title, link)
			}
			msg := line.TextMessage{
				Text: tmpl,
			}
			line.PushMessage(bot.AccessToken, v.ID, msg)
		}
	}
}
