package background

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"/golang/mtfosbot/model"

	twitchirc "git.trj.tw/golang/mtfosbot/module/twitch-irc"
)

func checkOpay() {
	channels, err := model.GetAllTwitchChannel()
	if err != nil {
		return
	}
	for _, v := range channels {
		if len(v.OpayID) > 0 && len(v.OpayID) == 32 {
			go getOpayData(v)
		}
	}
}

type opayResp struct {
	LstDonate []donateList `json:"lstDonate"`
	Settings  opaySetting  `json:"settings"`
}

type donateList struct {
	DonateID string `json:"donateid"`
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	MSG      string `json:"msg"`
}

type opaySetting struct {
	BGColor     string `json:"BgColor"`
	FontAnimate string `json:"FontAnimate"`
	MSGTemplate string `json:"MsgTemplate"`
}

func getOpayData(ch *model.TwitchChannel) {
	u := fmt.Sprintf("https://payment.opay.tw/Broadcaster/CheckDonate/%s", ch.OpayID)

	req, err := http.NewRequest("POST", u, strings.NewReader("{}"))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/63.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Opay http response code ::: ", resp.StatusCode)
		return
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		fmt.Println("Opay resp not json")
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	oResp := opayResp{}
	err = json.Unmarshal(bodyBytes, &oResp)
	if err != nil {
		return
	}

	if len(oResp.LstDonate) == 0 {
		return
	}

	var ids []string
	for _, v := range oResp.LstDonate {
		ids = append(ids, v.DonateID)
	}

	donateList, err := model.GetDonateListWithIDs(ids)
	if err != nil {
		return
	}

	if len(donateList) > 0 {
		for i := 0; i < len(oResp.LstDonate); i++ {
			for _, v := range donateList {
				if v.DonateID == oResp.LstDonate[i].DonateID {
					oResp.LstDonate[i].DonateID = ""
				}
			}
		}
	}

	for _, v := range oResp.LstDonate {
		if len(v.DonateID) > 0 {
			donateData := &model.OpayDonateList{
				OpayID:   ch.OpayID,
				DonateID: v.DonateID,
				Price:    v.Amount,
				Text:     v.MSG,
				Name:     v.Name,
			}
			err = donateData.InsertData()
			if err == nil && ch.Join {
				msg := fmt.Sprintf("/me 感謝 %s 贊助了 %d 元, %s", v.Name, v.Amount, v.MSG)
				twitchirc.SendMessage(ch.Name, msg)
			}
		}
	}

}
