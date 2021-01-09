package msgcmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"git.trj.tw/golang/mtfosbot/module/apis/twitch"
	"git.trj.tw/golang/mtfosbot/module/config"

	"git.trj.tw/golang/mtfosbot/model"
	googleapi "git.trj.tw/golang/mtfosbot/module/apis/google"
	lineobj "git.trj.tw/golang/mtfosbot/module/line-message/line-object"
)

func selectAct(cmd, sub, txt string, s *lineobj.SourceObject) (res string) {
	switch cmd {
	case "addgroup":
		return addLineGroup(sub, txt, s)
	case "addpage":
		return addFacebookPage(sub, txt, s)
	case "addtwitch":
		return addTwitchChannel(sub, txt, s)
	case "delpage":
		return delFacebookPage(sub, txt, s)
	case "deltwitch":
		return delTwitchChannel(sub, txt, s)
	case "image":
		return fmt.Sprintf("$image$%s", sub)
	case "video":
		return fmt.Sprintf("$video$%s", sub)
	case "addyoutube":
		return addYoutubeChannel(sub, txt, s)
	case "delyoutube":
		return delYoutubeChannel(sub, txt, s)
	case "lottery":
		return lottery(sub, txt, s)
	case "hello":
		return "World!!"
	}
	return
}

func addLineGroup(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = groupName notify
	exists, err := model.CheckGroup(s.GroupID)
	if err != nil {
		return "run check group error"
	}
	if exists {
		return "group exists"
	}
	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 2 {
		return "command args not match"
	}
	i, err := strconv.ParseInt(args[1], 10, 8)
	if err != nil || i < 0 || i > 1 {
		return "notify plases input 1 or 0"
	}

	_, err = model.AddLineGroup(args[0], s.UserID, i == 1)
	if err != nil {
		return "add group fail"
	}

	return "Success"
}

func checkGroupOwner(s *lineobj.SourceObject) (ok bool, err error) {
	exists, err := model.CheckGroup(s.GroupID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	ok, err = model.CheckGroupOwner(s.UserID, s.GroupID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func addFacebookPage(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = pageid tmpl
	ok, err := checkGroupOwner(s)
	if err != nil {
		return "check group fail"
	}
	if !ok {
		return "not owner"
	}

	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 2 {
		return "command args not match"
	}

	page, err := model.GetFacebookPage(args[0])
	if err != nil {
		return "check facebook page fail"
	}
	if page == nil {
		page = &model.FacebookPage{
			ID: args[0],
		}
		err = page.AddPage()
		if err != nil {
			return "add facebook page fail"
		}
	}

	rt := &model.LineFacebookRT{
		Line:     s.GroupID,
		Facebook: args[0],
		Tmpl:     strings.Join(args[1:], " "),
	}

	err = rt.AddRT()
	if err != nil {
		return "add facebook page fail"
	}

	return "Success"
}

func delFacebookPage(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = pageid
	ok, err := checkGroupOwner(s)
	if err != nil {
		return "check group fail"
	}
	if !ok {
		return "not owner"
	}

	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 1 {
		return "commage arg not match"
	}

	rt := &model.LineFacebookRT{
		Line:     s.GroupID,
		Facebook: args[0],
	}
	err = rt.DelRT()
	if err != nil {
		return "remove facebook page fail"
	}

	return "Success"
}

func checkTwitchType(t string) bool {
	switch t {
	case "live":
	default:
		return false
	}
	return true
}

func addTwitchChannel(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = twitchLogin type tmpl
	ok, err := checkGroupOwner(s)
	if err != nil {
		return "check group fail"
	}
	if !ok {
		return "not owner"
	}

	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 3 {
		return "command args not match"
	}

	if !checkTwitchType(args[1]) {
		return "type not allow"
	}

	info := twitch.GetUserDataByName(args[0])
	if info == nil {
		return "get twitch user id fail"
	}

	ch, err := model.GetTwitchChannelWithName(args[0])
	if err != nil {
		return "check channel fail"
	}
	if ch == nil {
		ch = &model.TwitchChannel{
			ID:   info.ID,
			Name: info.Login,
		}
		err = ch.Add()
		if err != nil {
			return "add twitch channel fail"
		}
	}

	rt := &model.LineTwitchRT{
		Line:   s.GroupID,
		Twitch: info.ID,
		Type:   args[1],
		Tmpl:   strings.Join(args[2:], " "),
	}
	err = rt.AddRT()
	if err != nil {
		return "add rt data fail"
	}

	return "Success"
}

func delTwitchChannel(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = twitchLogin type
	ok, err := checkGroupOwner(s)
	if err != nil {
		return "check group fail"
	}
	if !ok {
		return "not owner"
	}

	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 2 {
		return "command arg not match"
	}

	if !checkTwitchType(args[1]) {
		return "type not allow"
	}

	ch := &model.TwitchChannel{
		Name: args[0],
	}
	err = ch.GetWithName()
	if err != nil {
		return "get channel data fail"
	}
	if ch == nil {
		return "Success"
	}

	rt := &model.LineTwitchRT{
		Line:   s.GroupID,
		Twitch: ch.ID,
		Type:   args[1],
	}

	err = rt.DelRT()
	if err != nil {
		return "delete rt fail"
	}

	return "Success"
}

func addYoutubeChannel(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = youtubeID tmpl
	ok, err := checkGroupOwner(s)
	if err != nil {
		return "check group fail"
	}
	if !ok {
		return "not owner"
	}

	args := strings.Split(strings.Trim(txt, " "), " ")
	if len(args) < 2 {
		return "command arg not match"
	}
	ytName, err := googleapi.QueryYoutubeName(args[0])
	if err != nil || len(ytName) == 0 {
		fmt.Println("get youtube channel name fail :: ", err, ytName)
		return "get youtube channel name fail"
	}

	ytData, err := model.GetYoutubeChannelWithID(args[0])
	if err != nil {
		fmt.Println("get channel info fail :: ", err)
		return "check youtube fail"
	}
	if ytData == nil {
		ytData = &model.YoutubeChannel{
			ID:   args[0],
			Name: ytName,
		}
		err = ytData.Add()
		if err != nil {
			return "add youtube channel fail"
		}
	}

	rt := &model.LineYoutubeRT{
		Line:    s.GroupID,
		Youtube: args[0],
		Tmpl:    strings.Join(args[1:], " "),
	}
	err = rt.AddRT()
	if err != nil {
		return "add youtube channel rt fail"
	}

	return "Success"
}

func delYoutubeChannel(sub, txt string, s *lineobj.SourceObject) (res string) {
	// args = youtubeID
	ok, err := checkGroupOwner(s)
	if err != nil {
		return "check group fail"
	}
	if !ok {
		return "not owner"
	}
	txt = strings.Trim(txt, " ")
	ytData, err := model.GetYoutubeChannelWithID(txt)
	if err != nil {
		return "check channel fail"
	}
	if ytData == nil {
		return "channel not exists"
	}
	rt := &model.LineYoutubeRT{
		Line:    s.GroupID,
		Youtube: ytData.ID,
	}
	err = rt.DelRT()
	if err != nil {
		return "delete channel fail"
	}
	return "Success"
}

func lottery(sub, txt string, s *lineobj.SourceObject) (res string) {
	if len(sub) == 0 {
		return ""
	}
	arr, err := model.GetRandomLotteryByTypeAndLimit(sub, 100)
	if err != nil || arr == nil {
		return ""
	}

	// shuffle array
	dest := make([]*model.Lottery, len(arr))
	perm := rand.Perm(len(arr))
	for idx, val := range perm {
		dest[val] = arr[idx]
	}
	defer func() {
		for i := range perm {
			dest[i] = nil
			arr[i] = nil
		}
		dest = nil
		arr = nil
		perm = nil
	}()

	data := dest[0]

	conf := config.GetConf()
	u := conf.URL
	if last := len(u); last > 0 && u[last-1] == '/' {
		u = u[:last]
	}
	oriURL := "/image/origin"
	thumbURL := "/image/thumbnail"
	if len(data.Message) == 0 {
		return ""
	}
	o := u + oriURL + "/" + data.Message + "?d=" + sub
	t := u + thumbURL + "/" + data.Message + "?d=" + sub
	return fmt.Sprintf("$image$%s;%s", o, t)
}
