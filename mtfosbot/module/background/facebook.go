package background

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"/golang/mtfosbot/model"
	"/golang/mtfosbot/module/apis/line"
)

var idRegex = []*regexp.Regexp{
	regexp.MustCompile(`[\?|&]id\=(\d+)`),
	regexp.MustCompile(`\/posts\/(\d+)`),
	regexp.MustCompile(`\/photos\/.+?\/(\d+)`),
	regexp.MustCompile(`\/videos\/(\d+)`),
}

// PageData - facebook fan page data
type PageData struct {
	ID   string
	Text string
	Time int32
	Link string
}

type byTime []*PageData

func (pd byTime) Len() int           { return len(pd) }
func (pd byTime) Swap(i, j int)      { pd[i], pd[j] = pd[j], pd[i] }
func (pd byTime) Less(i, j int) bool { return pd[i].Time < pd[j].Time }

func readFacebookPage() {
	pages, err := model.GetAllFacebookPage()
	if err != nil {
		fmt.Println("get page data err ", err)
		return
	}

	for _, v := range pages {
		go getPageHTML(v)
	}
}

func getPageHTML(page *model.FacebookPage) {
	err := page.GetGroups()
	if err != nil {
		fmt.Println("get page group err ", err)
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://facebook.com/%s", page.ID))
	if err != nil {
		fmt.Println("get page html err ", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("parse doc err ", err)
		return
	}
	var pageData []*PageData
	sel := doc.Find("div.userContentWrapper")
	sel.Each(func(idx int, s *goquery.Selection) {
		timeEl := s.Find("abbr")
		time, timeExists := timeEl.Attr("data-utime")
		if !timeExists {
			fmt.Println("time not found")
			return
		}
		link, linkExists := timeEl.Parent().Attr("href")
		if !linkExists {
			fmt.Println("link not found")
			return
		}
		postContent := s.Find("div.userContent")
		text := postContent.Text()
		postID, idExists := postContent.First().Attr("id")

		if !idExists {
			idFlag := false
			for _, v := range idRegex {
				if v.MatchString(link) {
					idFlag = true
					m := v.FindStringSubmatch(link)
					postID = m[1]
				}
			}
			if !idFlag {
				fmt.Println("id not found")
				return
			}
		}
		fmt.Printf("Time: %s / Text: %s / ID: %s \n", time, text, postID)

		timeInt, err := strconv.ParseInt(time, 10, 32)

		if err != nil {
			fmt.Println("time parse err ", err)
			return
		}

		re := regexp.MustCompile(`^\/`)
		pageLink := fmt.Sprintf("https://www.facebook.com/%s", re.ReplaceAllString(link, ""))

		data := &PageData{
			ID:   postID,
			Text: text,
			Time: int32(timeInt),
			Link: pageLink,
		}

		pageData = append(pageData, data)
	})

	if len(pageData) == 0 {
		return
	}

	sort.Sort(sort.Reverse(byTime(pageData)))

	lastData := pageData[0]
	t := int32(time.Now().Unix())

	if (t-600) < lastData.Time && lastData.ID != page.LastPost {
		err = page.UpdatePost(lastData.ID)
		if err != nil {
			return
		}

		for _, v := range page.Groups {
			if v.Notify {
				bot, err := v.GetBot()
				if err != nil || bot == nil {
					log.Println("get group binding bot fail :: ", err)
					continue
				}
				tmpl := v.Tmpl
				if len(tmpl) > 0 {
					tmpl = strings.Replace(tmpl, "{link}", lastData.Link, -1)
					tmpl = strings.Replace(tmpl, "{txt}", lastData.Text, -1)
				} else {
					tmpl = fmt.Sprintf("%s\n%s", lastData.Text, lastData.Link)
				}
				msg := line.TextMessage{
					Text: tmpl,
				}
				line.PushMessage(bot.AccessToken, v.ID, msg)
			}
		}
	}
}
