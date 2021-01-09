package api

import (
	"time"

	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/twitch"
	"git.trj.tw/golang/mtfosbot/module/context"
	"git.trj.tw/golang/mtfosbot/module/twitch-irc"
	"git.trj.tw/golang/mtfosbot/module/utils"
	"github.com/gin-gonic/contrib/sessions"
)

func sessionTypeTwitch(id string) (ch []*model.TwitchChannel, err error) {
	chdata, err := model.GetTwitchChannelWithID(id)
	if err != nil {
		return
	}
	if chdata == nil {
		return nil, nil
	}
	list := make([]*model.TwitchChannel, 1)
	list = append(list, chdata)
	return list, nil
}

func sessionTypeSystem() (ch []*model.TwitchChannel, err error) {
	ch, err = model.GetAllTwitchChannel()
	return
}

// GetChannels - middleware
func GetChannels(c *context.Context) {
	sess := sessions.Default(c.Context)
	logingType := sess.Get("loginType")
	if logingType == nil {
		c.LoginFirst(nil)
		return
	}

	ltype, ok := logingType.(string)
	if !ok {
		c.LoginFirst(nil)
		return
	}
	if ltype == "twitch" {
		u := sess.Get("user")
		if u == nil {
			c.LoginFirst(nil)
			return
		}
		user, ok := u.(twitch.UserInfo)
		if !ok {
			c.LoginFirst(nil)
			return
		}
		chs, err := sessionTypeTwitch(user.ID)
		if err != nil || chs == nil {
			c.ServerError(nil)
			return
		}
		c.Set("channels", chs)
	} else if ltype == "system" {
		u := sess.Get("user")
		if u == nil {
			c.LoginFirst(nil)
			return
		}
		_, ok := u.(model.Account)
		if !ok {
			c.LoginFirst(nil)
			return
		}
		chs, err := sessionTypeSystem()
		if err != nil || chs == nil {
			c.LoginFirst(nil)
			return
		}
		c.Set("channels", chs)
	} else {
		c.LoginFirst(nil)
		return
	}

	c.Next()
}

func hasChannel(id string, c *context.Context) *model.TwitchChannel {
	if len(id) == 0 {
		return nil
	}

	channels, exists := c.Get("channels")
	if !exists {
		return nil
	}
	chs, ok := channels.([]*model.TwitchChannel)
	if !ok {
		return nil
	}

	for _, v := range chs {
		if v.ID == id {
			return v
		}
	}

	return nil
}

// GetChannelList -
func GetChannelList(c *context.Context) {
	channels, exists := c.Get("channels")
	if !exists {
		c.ServerError(nil)
		return
	}
	list, ok := channels.([]*model.TwitchChannel)
	if !ok {
		c.ServerError(nil)
		return
	}

	mapList := make([]map[string]interface{}, 0)
	for _, v := range list {
		if v == nil {
			continue
		}
		mapList = append(mapList, utils.ToMap(v))
	}

	c.Success(map[string]interface{}{
		"list": mapList,
	})
}

// GetChannelData -
func GetChannelData(c *context.Context) {
	chid := c.Param("chid")
	chdata := hasChannel(chid, c)
	if chdata == nil {
		c.NotFound(nil)
		return
	}

	c.Success(map[string]interface{}{
		"channel": utils.ToMap(chdata),
	})
}

// BotJoinChannel -
func BotJoinChannel(c *context.Context) {
	chid := c.Param("chid")
	chdata := hasChannel(chid, c)
	if chdata == nil {
		c.NotFound(nil)
		return
	}

	bodyArg := struct {
		Join int `json:"join"`
	}{}
	err := c.BindData(&bodyArg)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	if bodyArg.Join != 0 && bodyArg.Join != 1 {
		c.DataFormat(nil)
		return
	}

	err = chdata.UpdateJoin(bodyArg.Join == 1)
	if err != nil {
		c.ServerError(nil)
		return
	}

	if bodyArg.Join == 1 {
		twitchirc.JoinChannel(chdata.Name)
	} else {
		twitchirc.LeaveChannel(chdata.Name)
	}

	c.Success(nil)
}

// OpayIDChange -
func OpayIDChange(c *context.Context) {
	chid := c.Param("chid")
	chdata := hasChannel(chid, c)
	if chdata == nil {
		c.NotFound(nil)
		return
	}

	bodyArg := struct {
		Opay string `json:"opay" binding:"required"`
	}{}
	err := c.BindData(&bodyArg)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	err = chdata.UpdateOpayID(bodyArg.Opay)
	if err != nil {
		c.ServerError(nil)
		return
	}

	c.Success(nil)
}

// GetDonateSetting -
func GetDonateSetting(c *context.Context) {
	chid := c.Param("chid")
	chdata := hasChannel(chid, c)
	if chdata == nil {
		c.NotFound(nil)
		return
	}

	ds, err := model.GetDonateSettingByChannel(chdata.ID)
	if err != nil {
		c.ServerError(nil)
		return
	}

	var mapData map[string]interface{}
	if ds != nil {
		mapData = utils.ToMap(ds)
	} else {
		mapData = map[string]interface{}{}
	}

	c.Success(map[string]interface{}{
		"setting": mapData,
	})
}

// UpdateDonateSetting -
func UpdateDonateSetting(c *context.Context) {
	chid := c.Param("chid")
	chdata := hasChannel(chid, c)
	if chdata == nil {
		c.NotFound(nil)
		return
	}

	bodyArg := struct {
		End         int64  `json:"end" binding:"exists"`
		Title       string `json:"title" binding:"required"`
		Amount      int    `json:"amount" binding:"exists"`
		StartAmount int    `json:"start_amount" binding:"exists"`
	}{}
	err := c.BindData(&bodyArg)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	if bodyArg.End > 10000000000-1 {
		bodyArg.End = bodyArg.End / 1000
	}

	t := time.Unix(bodyArg.End, 0)

	ds := &model.DonateSetting{
		Title:        bodyArg.Title,
		EndDate:      t,
		StartDate:    time.Now(),
		StartAmount:  bodyArg.StartAmount,
		TargetAmount: bodyArg.Amount,
		Twitch:       chdata.ID,
	}
	err = ds.InsertOrUpdate()
	if err != nil {
		c.ServerError(nil)
		return
	}

	c.Success(nil)
}

// GetDonateBarStatus -
func GetDonateBarStatus(c *context.Context) {
	chid := c.Param("chid")
	chdata, err := model.GetTwitchChannelWithID(chid)
	if err != nil {
		c.ServerError(nil)
		return
	}
	if chdata == nil {
		c.NotFound(nil)
		return
	}

	ds, err := model.GetDonateSettingByChannel(chdata.ID)
	if err != nil {
		c.ServerError(nil)
		return
	}

	sum := 0
	mapData := map[string]interface{}{}
	if ds != nil {
		sum, err = model.SumChannelDonatePriceSinceTime(chdata.ID, ds.StartDate)
		if err != nil {
			c.ServerError(nil)
			return
		}
		sum += ds.StartAmount
		mapData = utils.ToMap(ds)
		mapData["total"] = sum
	}
	c.Success(map[string]interface{}{
		"setting": mapData,
	})
}
