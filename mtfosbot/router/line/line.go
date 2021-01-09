package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/context"
	linemsg "git.trj.tw/golang/mtfosbot/module/line-message"
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
)

// GetRawBody - line webhook body get
func GetRawBody(c *context.Context) {
	byteBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.DataFormat("body read fail")
		return
	}
	c.Set("rawbody", byteBody)
	c.Next()
}

// VerifyLine - middleware
func VerifyLine(c *context.Context) {
	rawbody, ok := c.Get("rawbody")
	if !ok {
		c.DataFormat("body read fail")
		return
	}
	var raw []byte
	if raw, ok = rawbody.([]byte); !ok {
		c.DataFormat("body type error")
		return
	}
	sign := c.GetHeader("X-Line-Signature")
	if len(sign) == 0 {
		c.Next()
		return
	}

	botid, ok := c.GetQuery("id")
	if !ok || len(botid) == 0 {
		c.CustomRes(403, map[string]string{
			"message": "no bot data",
		})
	}

	bot, err := model.GetBotInfo(botid)
	if err != nil {
		c.ServerError(nil)
		return
	}

	hash := hmac.New(sha256.New, []byte(bot.Secret))
	_, err = hash.Write(raw)
	if err != nil {
		c.ServerError(nil)
		return
	}
	hashSign := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if hashSign != sign {
		c.CustomRes(403, map[string]string{
			"message": "sign verify fail",
		})
		return
	}
	c.Next()
}

// GetLineMessage -
func GetLineMessage(c *context.Context) {
	rawbody, ok := c.Get("rawbody")
	if !ok {
		c.DataFormat("body read fail")
	}
	var raw []byte
	if raw, ok = rawbody.([]byte); !ok {
		c.DataFormat("body type error")
	}
	botid, ok := c.GetQuery("id")
	if !ok || len(botid) == 0 {
		c.CustomRes(403, map[string]string{
			"message": "no bot data",
		})
	}

	events := struct {
		Events []*lineobj.EventObject `json:"events"`
	}{}

	err := json.Unmarshal(raw, &events)
	if err != nil {
		c.ServerError(nil)
		return
	}

	if len(events.Events) > 0 {
		for _, v := range events.Events {
			log.Println("get line message :: ", v)
			go linemsg.MessageEvent(botid, v)
		}
	}

	c.Success(nil)
}
