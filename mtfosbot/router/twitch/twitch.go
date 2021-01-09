package twitch

import (
	"net/url"
	"strings"

	"git.trj.tw/golang/mtfosbot/model"
	twitchapi "git.trj.tw/golang/mtfosbot/module/apis/twitch"
	"git.trj.tw/golang/mtfosbot/module/config"
	"git.trj.tw/golang/mtfosbot/module/context"
	"github.com/gin-gonic/contrib/sessions"
)

// OAuthLogin -
func OAuthLogin(c *context.Context) {
	session := sessions.Default(c.Context)
	twOauth := `https://id.twitch.tv/oauth2/authorize`
	conf := config.GetConf()
	redirectTo := strings.TrimRight(conf.URL, "/")
	redirectTo += "/twitch/oauth"
	qs := url.Values{}
	qs.Add("client_id", conf.Twitch.ClientID)
	qs.Add("redirect_uri", redirectTo)
	qs.Add("response_type", "code")
	qs.Add("scope", "user:read:email")

	toURL, ok := c.GetQuery("tourl")
	if ok && len(toURL) > 0 {
		session.Set("backUrl", toURL)
		session.Save()
	}

	c.Redirect(302, twOauth+"?"+qs.Encode())
}

// OAuthProc -
func OAuthProc(c *context.Context) {
	code, ok := c.GetQuery("code")
	if !ok || len(code) == 0 {
		c.DataFormat(nil)
		return
	}

	tokenData, err := twitchapi.GetTokenData(code)
	if err != nil {
		c.DataFormat("token get fail")
		return
	}

	session := sessions.Default(c.Context)

	userData := twitchapi.GetUserDataByToken(tokenData.AccessToken)
	if userData == nil {
		c.ServerError(nil)
		return
	}

	session.Set("token", tokenData)
	session.Set("user", userData)
	session.Set("loginType", "twitch")

	chData, err := model.GetTwitchChannelWithID(userData.ID)
	if err != nil {
		c.ServerError(nil)
		return
	}
	if chData == nil {
		chData = &model.TwitchChannel{
			ID:   userData.ID,
			Name: userData.Login,
		}
		err = chData.Add()
		if err != nil {
			c.ServerError(nil)
			return
		}
	} else {
		if userData.Login != chData.Name {
			chData.UpdateName(userData.Login)
		}
	}

	conf := config.GetConf()
	goURL := strings.TrimRight(conf.URL, "/") + "/web"
	tourl := session.Get("backUrl")
	if tourl != nil {
		if str, ok := tourl.(string); ok {
			goURL = str
			session.Delete("backUrl")
		}
	}
	session.Save()
	c.Redirect(301, goURL)
}
