package msgcmd

import (
	"fmt"
	"regexp"
	"strings"

	"git.trj.tw/golang/mtfosbot/model"

	"git.trj.tw/golang/mtfosbot/module/apis/line"
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
)

func parseCMD(in string) (c [][]string) {
	re, err := regexp.Compile("{{(.+?)}}")
	if err != nil {
		return
	}

	c = re.FindAllStringSubmatch(in, -1)

	return
}

// ParseLineMsg -
func ParseLineMsg(botid, txt, replyToken string, source *lineobj.SourceObject) {
	if source.Type != "group" {
		return
	}
	bot, err := model.GetBotInfo(botid)
	if err != nil || bot == nil {
		fmt.Println("get bot err or no bot :: ", err)
		return
	}
	strs := strings.Split(strings.Trim(txt, " "), " ")

	// skip empty string
	if len(strs[0]) == 0 {
		return
	}

	firstNum := []rune(strs[0])[0]

	if firstNum == 33 || firstNum == 65281 {
		// nor cmd
		cmd := ""
		if firstNum == 65281 {
			cmd = strs[0][3:]
		} else {
			cmd = strs[0][1:]
		}
		if len(cmd) == 0 {
			return
		}
		c, err := model.GetGroupCommand(strings.ToLower(cmd), source.GroupID)
		if err != nil || c == nil {
			return
		}

		resStrs := runCMD(strings.Join(strs[1:], " "), c.Message, source)
		msgs := make([]interface{}, 0)
		for _, v := range resStrs {
			m := parseResult(v)
			msgs = append(msgs, m)
		}
		line.ReplyMessage(bot.AccessToken, replyToken, msgs...)

	} else {
		// key cmd
		c, err := model.GetGroupKeyCommand(strings.ToLower(strs[0]), source.GroupID)
		if err != nil || c == nil {
			return
		}

		resStrs := runCMD(strings.Join(strs[1:], " "), c.Message, source)
		msgs := make([]interface{}, 0)
		for _, v := range resStrs {
			m := parseResult(v)
			msgs = append(msgs, m)
		}

		line.ReplyMessage(bot.AccessToken, replyToken, msgs...)

	}
}

func parseResult(str string) interface{} {
	var m interface{}

	if strings.HasPrefix(str, "$image$") {
		str = strings.Replace(str, "$image$", "", 1)
		strs := strings.Split(str, ";")
		m = line.ImageMessage{
			OriginalContentURL: strs[0],
			PreviewImageURL:    strs[1],
		}
	} else if strings.HasPrefix(str, "$video$") {
		str = strings.Replace(str, "$video$", "", 1)
		strs := strings.Split(str, ";")
		m = line.VideoMessage{
			OriginalContentURL: strs[0],
			PreviewImageURL:    strs[1],
		}
	} else {
		m = line.TextMessage{
			Text: str,
		}
	}

	return m
}

func runCMD(txt, c string, s *lineobj.SourceObject) (res []string) {
	// TODO: 把多重回覆指令 回應內容獨立資料表
	cmds := strings.Split(c, "$#$")
	if len(cmds) == 0 {
		return
	}
	for _, c := range cmds {
		cmdAct := parseCMD(c)
		if len(cmdAct) == 0 {
			res = append(res, c)
			continue
		}
		tmpRes := c
		for _, v := range cmdAct {
			if len(v) > 1 {
				// run cmd
				m := strings.Split(v[1], "=")
				sub := ""
				if len(m) > 1 {
					sub = strings.Join(m[1:], " ")
				}
				cmdRes := selectAct(m[0], sub, txt, s)
				tmpRes = strings.Replace(tmpRes, v[0], cmdRes, 1)
			}
		}
		res = append(res, tmpRes)
	}
	return
}
