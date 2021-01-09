package api

import (
	"git.trj.tw/golang/mtfosbot/model"
	"git.trj.tw/golang/mtfosbot/module/apis/twitch"
	"git.trj.tw/golang/mtfosbot/module/context"
	"github.com/gin-gonic/contrib/sessions"
	"golang.org/x/crypto/bcrypt"
)

// CheckSession -
func CheckSession(c *context.Context) {
	session := sessions.Default(c.Context)
	userData := session.Get("user")
	loginType := session.Get("loginType")
	if userData == nil || loginType == nil {
		c.LoginFirst(nil)
		return
	}
	var name string
	var ltype string
	var ok bool
	switch userData.(type) {
	case model.Account:
		name = userData.(model.Account).Account
	case twitch.UserInfo:
		name = userData.(twitch.UserInfo).Login
	default:
		c.LoginFirst(nil)
		return
	}
	if ltype, ok = loginType.(string); !ok {
		c.LoginFirst(nil)
		return
	}

	loginUser := map[string]string{
		"name": name,
		"type": ltype,
	}
	session.Set("loginUser", loginUser)
	session.Save()
	c.Next()
}

// UserLogin - system user login
func UserLogin(c *context.Context) {
	bodyArg := struct {
		Account  string `form:"account" json:"account" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}{}
	err := c.BindData(&bodyArg)
	if err != nil {
		c.DataFormat(nil)
		return
	}

	acc, err := model.GetAccount(bodyArg.Account)
	if err != nil {
		c.ServerError(`account or password error`)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(bodyArg.Password))
	if err != nil {
		c.DataFormat(`account or password error`)
		return
	}

	session := sessions.Default(c.Context)

	acc.Password = ""
	session.Set("user", acc)
	session.Set("loginType", "system")
	session.Save()

	c.Success(nil)
}

// UserLogout -
func UserLogout(c *context.Context) {
	session := sessions.Default(c.Context)
	session.Clear()
	session.Save()

	c.Success(nil)
}

// GetSessionData -
func GetSessionData(c *context.Context) {
	session := sessions.Default(c.Context)
	loginUser := session.Get("loginUser")
	if loginUser == nil {
		c.LoginFirst(nil)
		return
	}

	user := map[string]interface{}{
		"user": loginUser,
	}
	c.Success(user)
}
