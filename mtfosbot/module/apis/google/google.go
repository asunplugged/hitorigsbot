package google

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"/golang/mtfosbot/module/apis"

	jsoniter "github.com/json-iterator/go"

	"/golang/mtfosbot/module/config"
)

var baseURL = "https://www.googleapis.com"
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func getURL(p string, querystring ...interface{}) (string, bool) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", false
	}
	ref, err := u.Parse(p)
	if err != nil {
		return "", false
	}
	if len(querystring) > 0 {
		switch querystring[0].(type) {
		case string:
			ref, err = ref.Parse(fmt.Sprintf("?%s", (querystring[0].(string))))
			if err != nil {
				return "", false
			}
			break
		default:
		}
	}

	str := ref.String()
	return str, true
}

func getHeaders(token ...interface{}) map[string]string {
	m := make(map[string]string)
	m["Content-Type"] = "application/json"
	return m
}

// ChannelItem -
type ChannelItem struct {
	ID      string         `json:"id"`
	Snippet ChannelSnippet `json:"snippet"`
}

// ChannelSnippet -
type ChannelSnippet struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CustomURL   string `json:"customUrl"`
}

// QueryYoutubeName -
func QueryYoutubeName(id string) (n string, err error) {
	conf := config.GetConf()
	if len(id) == 0 {
		return "", errors.New("id is empty")
	}
	qs := url.Values{}
	qs.Add("id", id)
	qs.Add("key", conf.Google.APIKey)
	qs.Add("part", "snippet")

	apiURL, ok := getURL("/youtube/v3/channels", qs.Encode())
	if !ok {
		return "", errors.New("url parser fail")
	}
	reqObj := apis.RequestObj{
		Method:  "GET",
		URL:     apiURL,
		Headers: getHeaders(),
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 || !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return "", errors.New("api response fail")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	apiRes := struct {
		Items []ChannelItem `json:"items"`
	}{}

	err = json.Unmarshal(bodyBytes, &apiRes)
	if err != nil {
		return "", err
	}

	if len(apiRes.Items) == 0 {
		return "", errors.New("channel data not found")
	}

	for _, v := range apiRes.Items {
		if v.ID == id {
			return v.Snippet.Title, nil
		}
	}

	return "", errors.New("channel data not found")
}

// SubscribeYoutube -
func SubscribeYoutube(id string) {
	if len(id) == 0 {
		return
	}
	conf := config.GetConf()
	apiURL := "https://pubsubhubbub.appspot.com/subscribe"
	cbURL, err := url.Parse(conf.URL)
	if err != nil {
		return
	}
	cbURL, err = cbURL.Parse(fmt.Sprintf("/google/youtube/webhook?id=%s", id))
	if err != nil {
		return
	}

	qs := url.Values{}
	qs.Add("hub.mode", "subscribe")
	qs.Add("hub.verify", "async")
	qs.Add("hub.topic", fmt.Sprintf("https://www.youtube.com/xml/feeds/videos.xml?channel_id=%s", id))
	qs.Add("hub.callback", cbURL.String())
	qs.Add("hub.lease_seconds", "86400")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(qs.Encode()))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}
