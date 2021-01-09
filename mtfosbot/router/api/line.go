package api

import (
	"fmt"
	"log"
	"strconv"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/line"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/utils"
)

// PushMessage -
func PushLineMessage(c *context.Context) {
	bodyData := struct {
		Group   string `json:"group,required"`
		Message string `json:"message,required"`
	}{}

	err := c.BindData(&bodyData)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	textObj := line.TextMessage{}
	textObj.Text = bodyData.Message

	group, err := model.GetLineGroup(bodyData.Group)
	if err != nil {
		c.ServerError(nil)
		log.Println("get group :: ", err)
		return
	}
	bot, err := group.GetBot()
	if err != nil || bot == nil {
		log.Println("get group binding bot fail :: ", err)
		c.ServerError(nil)
		return
	}

	line.PushMessage(bot.AccessToken, bodyData.Group, textObj)

	c.Success(nil)
}

// GetLineList -
func GetLineList(c *context.Context) {
	ls, err := model.GetLineGroupList()
	if err != nil {
		c.ServerError(nil)
		return
	}
	list := make([]map[string]interface{}, 0)

	if ls != nil {
		for _, v := range ls {
			if v == nil {
				continue
			}
			list = append(list, utils.ToMap(v))
		}
	}

	c.Success(map[string]interface{}{
		"list": list,
	})
}

// GetLineMessageLog -
func GetLineMessageLog(c *context.Context) {
	numP := 1
	if p, ok := c.GetQuery("p"); ok {
		if i, err := strconv.Atoi(p); err == nil {
			numP = i
		}
	}
	numMax := 20
	if max, ok := c.GetQuery("max"); ok {
		if m, err := strconv.Atoi(max); err == nil {
			numMax = m
		}
	}

	g := c.DefaultQuery("group", "")
	u := c.DefaultQuery("user", "")
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	where := make(map[string]string)
	if len(g) > 0 {
		where["group"] = g
	}
	if len(u) > 0 {
		where["user"] = u
	}

	count, err := model.GetLineMessageLogCount(where)
	if err != nil {
		fmt.Println("count error :::: ", err)
		c.ServerError(nil)
		return
	}

	page := utils.CalcPage(count, numP, numMax)
	if numP == -1 {
		numP = page.Total
	}
	page = utils.CalcPage(count, numP, numMax)

	logs, err := model.GetLineMessageLogList(g, u, page.Offset, page.Limit, order)
	if err != nil {
		c.ServerError(nil)
		return
	}

	resMap := make([]map[string]interface{}, len(logs))

	for k, v := range logs {
		m := utils.ToMap(v.LineMessageLog)
		m["group_name"] = v.GroupName
		m["user_name"] = v.UserName
		resMap[k] = m
	}

	c.Success(map[string]interface{}{
		"list": resMap,
		"page": map[string]interface{}{
			"cur":   page.Page,
			"total": page.Total,
		},
	})
}

// GetCommandList -
func GetCommandList(c *context.Context) {
	numP := 1
	if p, ok := c.GetQuery("p"); ok {
		if i, err := strconv.Atoi(p); err == nil {
			numP = i
		}
	}

	numMax := 20
	if max, ok := c.GetQuery("max"); ok {
		if m, err := strconv.Atoi(max); err == nil {
			numMax = m
		}
	}

	g := c.DefaultQuery("group", "")
	cmd := c.DefaultQuery("cmd", "")

	whereMap := make(map[string]string)

	if len(g) > 0 {
		whereMap["group"] = g
	}
	if len(cmd) > 0 {
		whereMap["cmd"] = cmd
	}

	count, err := model.GetCommandCount(whereMap)
	if err != nil {
		log.Println("error :::: ", err)
		c.ServerError(nil)
		return
	}

	page := utils.CalcPage(count, numP, numMax)

	cmds, err := model.GetCommands(whereMap, page.Offset, page.Limit, nil)
	if err != nil {
		log.Println("list error :::: ", err)
		c.ServerError(nil)
		return
	}

	cmdList := make([]map[string]interface{}, len(cmds))
	for k, v := range cmds {
		tmp := utils.ToMap(v.Commands)
		tmp["group_name"] = v.GroupName
		cmdList[k] = tmp
		tmp = nil
	}

	c.Success(map[string]interface{}{
		"list": cmdList,
		"page": map[string]interface{}{
			"cur":   page.Page,
			"total": page.Total,
		},
	})
}

// AddLineGroupCommand -
func AddLineGroupCommand(c *context.Context) {
	bodyStruct := struct {
		Cmd     string `json:"cmd"`
		Message string `json:"message"`
		Group   string `json:"group"`
	}{}

	err := c.BindData(&bodyStruct)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	if len(bodyStruct.Cmd) == 0 || len(bodyStruct.Message) == 0 {
		c.DataFormat(nil)
		return
	}

	cmd, err := model.AddCommand(bodyStruct.Cmd, bodyStruct.Message, bodyStruct.Group)
	if err != nil {
		c.ServerError(nil)
		return
	}

	c.Success(map[string]interface{}{
		"cmd": utils.ToMap(cmd),
	})
}

// DeleteLineGroupCommand -
func DeleteLineGroupCommand(c *context.Context) {
	cmd, ok := c.Params.Get("cmd")
	if !ok {
		c.DataFormat(nil)
		return
	}
	g := c.DefaultQuery("group", "")

	exist, err := model.CheckCommand(cmd, g)
	if err != nil {
		c.ServerError(nil)
		return
	}
	if !exist {
		c.NotFound("command not found")
		return
	}
	err = model.DeleteCommand(cmd, g)
	if err != nil {
		c.ServerError(nil)
		return
	}

	c.Success(nil)
}

// EditLineGroupCommand -
func EditLineGroupCommand(c *context.Context) {
	bodyStruct := struct {
		Message string `json:"message"`
	}{}
	cmd, ok := c.Params.Get("cmd")
	if !ok {
		c.DataFormat(nil)
		return
	}
	g := c.DefaultQuery("group", "")

	err := c.BindData(&bodyStruct)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	exist, err := model.CheckCommand(cmd, g)
	if err != nil {
		c.ServerError(nil)
		return
	}
	if !exist {
		c.NotFound("command not found")
		return
	}

	err = model.UpdateCommand(cmd, g, bodyStruct.Message)
	if err != nil {
		c.ServerError(nil)
		return
	}

	c.Success(nil)
}
