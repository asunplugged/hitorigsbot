package twitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"/golang/mtfosbot/module/apis"
	"/golang/mtfosbot/module/config"
)

// UserInfo - twitch user info data
type UserInfo struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

var baseURL = "https://api.twitch.tv"

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
	conf := config.GetConf()
	m["Content-Type"] = "application/json"
	if len(token) > 0 {
		switch token[0].(type) {
		case string:
			m["Authorization"] = fmt.Sprintf("Bearer %s", (token[0].(string)))
			break
		default:
		}
	}
	m["Client-ID"] = conf.Twitch.ClientID
	return m
}

// GetUserDataByToken - get token own user data
func GetUserDataByToken(token string) (userInfo *UserInfo) {
	if len(token) == 0 {
		return
	}
	url, ok := getURL("/helix/users")
	if !ok {
		return
	}

	reqObj := apis.RequestObj{}
	reqObj.Headers = getHeaders(token)
	reqObj.Method = "GET"
	reqObj.URL = url
	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return
	}

	apiData := struct {
		Data []*UserInfo `json:"data"`
	}{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bodyBytes, &apiData)
	if err != nil {
		return
	}

	if len(apiData.Data) == 0 {
		return
	}

	return apiData.Data[0]
}

// GetUserDataByName -
func GetUserDataByName(login string) (userInfo *UserInfo) {
	if len(login) == 0 {
		return
	}
	qsValue := url.Values{}
	qsValue.Add("login", login)
	url, ok := getURL("/helix/users", qsValue.Encode())
	if !ok {
		return
	}

	reqObj := apis.RequestObj{}
	reqObj.Headers = getHeaders()
	reqObj.Method = "GET"
	reqObj.URL = url
	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return
	}

	apiData := struct {
		Data []*UserInfo `json:"data"`
	}{}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bodyBytes, &apiData)
	if err != nil {
		return
	}

	if len(apiData.Data) == 0 {
		return
	}

	return apiData.Data[0]
}

// GetUserDataByID -
func GetUserDataByID(id string) (userInfo *UserInfo) {
	if len(id) == 0 {
		return
	}
	qsValue := url.Values{}
	qsValue.Add("id", id)
	url, ok := getURL("/helix/users", qsValue.Encode())
	if !ok {
		return
	}

	reqObj := apis.RequestObj{}
	reqObj.Headers = getHeaders()
	reqObj.Method = "GET"
	reqObj.URL = url
	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return
	}

	apiData := struct {
		Data []*UserInfo `json:"data"`
	}{}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bodyBytes, &apiData)
	if err != nil {
		return
	}

	if len(apiData.Data) == 0 {
		return
	}

	return apiData.Data[0]
}

// StreamInfo -
type StreamInfo struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	GameID       string `json:"game_id"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	ViewerCount  int    `json:"viewer_count"`
	StartedAt    string `json:"started_at"`
	Language     string `json:"language"`
	ThumbnailURL string `json:"thumbnail_url"`
}

// GetUserStreamStatus -
func GetUserStreamStatus(ids []string) (info []*StreamInfo) {
	if len(ids) == 0 {
		return
	}

	apiData := struct {
		Data []*StreamInfo `json:"data"`
	}{}

	qsValue := url.Values{}
	for _, v := range ids {
		qsValue.Add("user_id", v)
	}
	url, ok := getURL("/helix/streams", qsValue.Encode())
	if !ok {
		return
	}

	reqObj := apis.RequestObj{}
	reqObj.Headers = getHeaders()
	reqObj.Method = "GET"
	reqObj.URL = url
	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bodyBytes, &apiData)
	if err != nil {
		return nil
	}

	return apiData.Data
}

// TwitchTokenData -
type TwitchTokenData struct {
	AccessToken  string   `json:"access_token" cc:"access_token"`
	RefreshToken string   `json:"refresh_token" cc:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in" cc:"expires_in"`
	Scope        []string `json:"scope" cc:"scope"`
	TokenType    string   `json:"token_type" cc:"token_type"`
}

// GetTokenData -
func GetTokenData(code string) (token *TwitchTokenData, err error) {
	if len(code) == 0 {
		return nil, errors.New("code is empty")
	}
	conf := config.GetConf()
	twitchURL := "https://id.twitch.tv/oauth2/token"
	redirectTo := strings.TrimRight(conf.URL, "/") + "/twitch/oauth"

	qs := url.Values{}
	qs.Add("client_id", conf.Twitch.ClientID)
	qs.Add("client_secret", conf.Twitch.ClientSecret)
	qs.Add("code", code)
	qs.Add("grant_type", "authorization_code")
	qs.Add("redirect_uri", redirectTo)

	// u, err := url.Parse(twitchURL)
	// if err != nil {
	// 	return nil, err
	// }
	// u, err = u.Parse(qs.Encode())
	// if err != nil {
	// 	return nil, err
	// }

	reqObj := apis.RequestObj{
		URL:    twitchURL + "?" + qs.Encode(),
		Method: "POST",
	}
	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 || !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil, errors.New("api response error")
	}

	err = json.Unmarshal(bodyBytes, &token)

	return
}
