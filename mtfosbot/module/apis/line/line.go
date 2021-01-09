package line

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"/golang/mtfosbot/module/apis"
)

// TextMessage - line text message object
type TextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ImageMessage - line image message object
type ImageMessage struct {
	Type               string `json:"type"`
	OriginalContentURL string `json:"originalContentUrl"`
	PreviewImageURL    string `json:"previewImageUrl"`
}

// VideoMessage - line video message object
type VideoMessage struct {
	Type               string `json:"type"`
	OriginalContentURL string `json:"OriginalContentUrl"`
	PreviewImageURL    string `json:"previewImageUrl"`
}

// LineUserInfo -
type LineUserInfo struct {
	DisplayName string `json:"displayName"`
	UserID      string `json:"userId"`
}

type pushBody struct {
	To       string        `json:"to"`
	Messages []interface{} `json:"messages"`
}
type replyBody struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

var baseURL = "https://api.line.me/"

func getURL(p string) (string, bool) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", false
	}
	ref, err := u.Parse(p)
	if err != nil {
		return "", false
	}
	str := ref.String()
	return str, true
}

func getHeaders(token string) map[string]string {
	m := make(map[string]string)
	m["Content-Type"] = "application/json"
	m["Authorization"] = fmt.Sprintf("Bearer %s", token)
	return m
}

func checkMessageObject(m interface{}) interface{} {
	if m == nil {
		return nil
	}

	var obj interface{}
	switch m.(type) {
	case ImageMessage:
		tmp := (m.(ImageMessage))
		tmp.Type = "image"
		obj = tmp
		break
	case TextMessage:
		tmp := (m.(TextMessage))
		tmp.Type = "text"
		obj = tmp
		break
	case VideoMessage:
		tmp := (m.(VideoMessage))
		tmp.Type = "video"
		obj = tmp
		break
	default:
		return nil
	}

	return obj
}

// PushMessage -
func PushMessage(accessToken, target string, message ...interface{}) {
	log.Println("push target :::: ", target)
	if len(target) == 0 || len(message) == 0 {
		return
	}
	urlPath := "/v2/bot/message/push"

	body := &pushBody{
		To: target,
	}

	checked := make([]interface{}, 0)
	for _, v := range message {
		tmp := checkMessageObject(v)
		if tmp == nil {
			continue
		}
		checked = append(checked, tmp)
	}

	body.Messages = append(body.Messages, checked...)
	if len(body.Messages) > 5 {
		body.Messages = body.Messages[:5]
	}
	dataByte, err := json.Marshal(body)
	if err != nil {
		log.Println("to json error ::::", err)
		return
	}

	byteReader := bytes.NewReader(dataByte)

	apiURL, ok := getURL(urlPath)
	if !ok {
		log.Println("get url fail ::::::")
		return
	}

	reqObj := apis.RequestObj{
		Method:  "POST",
		URL:     apiURL,
		Headers: getHeaders(accessToken),
		Body:    byteReader,
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		log.Println("get req fail :::::: ", err)
		return
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("send to line error :::: ", err)
		return
	}
}

// ReplyMessage -
func ReplyMessage(accessToken, replyToken string, message ...interface{}) {
	if len(replyToken) == 0 || len(message) == 0 {
		return
	}
	urlPath := "/v2/bot/message/reply"

	body := &replyBody{
		ReplyToken: replyToken,
	}

	checked := make([]interface{}, 0)
	for _, v := range message {
		tmp := checkMessageObject(v)
		if tmp == nil {
			continue
		}
		checked = append(checked, tmp)
	}

	body.Messages = append(body.Messages, checked...)
	if len(body.Messages) > 5 {
		body.Messages = body.Messages[:5]
	}
	dataByte, err := json.Marshal(body)
	if err != nil {
		return
	}

	byteReader := bytes.NewReader(dataByte)

	apiURL, ok := getURL(urlPath)
	if !ok {
		return
	}

	reqObj := apis.RequestObj{
		Method:  "POST",
		URL:     apiURL,
		Headers: getHeaders(accessToken),
		Body:    byteReader,
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
}

// GetUserInfo -
func GetUserInfo(accessToken, u, g string) (user *LineUserInfo, err error) {
	urlPath := fmt.Sprintf("/v2/bot/group/%s/member/%s", g, u)
	header := getHeaders(accessToken)
	apiURL, ok := getURL(urlPath)
	if !ok {
		return nil, errors.New("url parser fail")
	}

	reqObj := apis.RequestObj{
		Method:  "GET",
		URL:     apiURL,
		Headers: header,
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

	if resp.StatusCode != 200 {
		return nil, errors.New("api response not 200")
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		return nil, errors.New("response body not json")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		return nil, err
	}

	return
}

// GetContentHead -
func GetContentHead(accessToken, id string) (mime string, err error) {
	urlPath := fmt.Sprintf("/v2/bot/message/%s/content", id)
	header := getHeaders(accessToken)
	u, ok := getURL(urlPath)
	if !ok {
		return "", errors.New("get url fail")
	}

	reqObj := apis.RequestObj{
		Method:  "HEAD",
		URL:     u,
		Headers: header,
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

	mime = resp.Header.Get("Content-Type")

	return
}

// DownloadContent -
func DownloadContent(accessToken, id string, w io.Writer) (err error) {
	urlPath := fmt.Sprintf("/v2/bot/message/%s/content", id)
	header := getHeaders(accessToken)
	u, ok := getURL(urlPath)
	if !ok {
		return errors.New("get url fail")
	}

	reqObj := apis.RequestObj{
		Method:  "GET",
		URL:     u,
		Headers: header,
	}

	req, err := apis.GetRequest(reqObj)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)

	return
}
