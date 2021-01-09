package google

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"git.trj.tw/golang/mtfosbot/model"
	lineapi "git.trj.tw/golang/mtfosbot/module/apis/line"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/utils"
)

type feed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   []entry  `xml:"entry"`
}
type entry struct {
	XMLName xml.Name `xml:"entry"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Link    link     `xml:"link"`
	Author  []author `xml:"author"`
}
type link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
}
type author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	URI     string   `xml:"uri"`
}

// VerifyWebhook -
func VerifyWebhook(c *context.Context) {
	hubMode, ok := c.GetQuery("hub.mode")
	if !ok {
		c.DataFormat(nil)
		return
	}
	challenge, ok := c.GetQuery("hub.challenge")
	if !ok {
		c.DataFormat(nil)
		return
	}
	id, ok := c.GetQuery("id")
	if !ok {
		c.DataFormat(nil)
		return
	}
	if hubMode == "subscribe" {
		t := time.Now().Unix() + 86400
		yt, err := model.GetYoutubeChannelWithID(id)
		if err != nil {
			c.ServerError(nil)
			return
		}
		if yt == nil {
			c.NotFound("channel not found")
			return
		}
		err = yt.UpdateExpire(t)
		if err != nil {
			c.ServerError(nil)
		}
	}
	c.String(200, challenge)
}

// GetNotifyWebhook -
func GetNotifyWebhook(c *context.Context) {
	byteBody, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		c.DataFormat(nil)
		return
	}
	fmt.Println("Hook string ::: ", string(byteBody))

	id, ok := c.GetQuery("id")
	if !ok {
		c.DataFormat(nil)
		return
	}

	hook := &feed{}
	err = xml.Unmarshal(byteBody, &hook)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	log.Println("hook data", utils.ToMap(hook))

	if len(hook.Entry) == 0 {
		c.Success(nil)
		return
	}

	yt, err := model.GetYoutubeChannelWithID(id)
	log.Println("youtube and error", yt, err)
	if err != nil || yt == nil {
		c.ServerError(nil)
		return
	}

	if hook.Entry[0].ID == yt.LastVideo {
		c.Success(nil)
		return
	}

	err = yt.UpdateLastVideo(hook.Entry[0].ID)
	if err != nil {
		c.ServerError(nil)
		return
	}

	err = yt.GetGroups()
	if err != nil {
		log.Println("get groups error ::::", err)
		c.ServerError(nil)
		return
	}
	log.Println("yt groups ::::: ", yt.Groups)

	for _, v := range yt.Groups {
		log.Println("group data :::: ", v, v.Notify, v.Name, v.ID)
		bot, err := v.GetBot()
		if err != nil || bot == nil {
			log.Println("get group binding bot fail :: ", err)
			continue
		}
		if v.Notify == true {
			str := v.Tmpl
			log.Println("template :::: ", str)
			if len(str) == 0 {
				str = fmt.Sprintf("%s\n%s", hook.Entry[0].Title, hook.Entry[0].Link.Href)
			} else {
				str = strings.Replace(str, "{link}", hook.Entry[0].Link.Href, -1)
				str = strings.Replace(str, "{txt}", hook.Entry[0].Title, -1)
			}

			msg := lineapi.TextMessage{
				Text: str,
			}
			log.Println("msg ::::: ", msg)

			lineapi.PushMessage(bot.AccessToken, v.ID, msg)
		}
	}

	c.Success(nil)
}
