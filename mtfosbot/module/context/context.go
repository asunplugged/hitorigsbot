package context

import (
	"git.trj.tw/golang/mtfosbot/module/apimsg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Context custom context struct
type Context struct {
	*gin.Context
}

// CustomMiddle func
type CustomMiddle func(*Context)

// PatchCtx - patch ctx to custom middle
func PatchCtx(handler CustomMiddle) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			Context: c,
		}
		handler(ctx)
	}
}

// BindData client body data
func (c *Context) BindData(i interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(i, b)
}

// CustomRes -
func (c *Context) CustomRes(status int, msg interface{}) {
	c.AbortWithStatusJSON(status, msg)
}

// LoginFirst -
func (c *Context) LoginFirst(msg interface{}) {
	obj := apimsg.GetRes("LoginFirst", msg)
	c.AbortWithStatusJSON(obj.Status, obj.Obj)
}

// NotFound -
func (c *Context) NotFound(msg interface{}) {
	obj := apimsg.GetRes("NotFound", msg)
	c.AbortWithStatusJSON(obj.Status, obj.Obj)
}

// DataFormat -
func (c *Context) DataFormat(msg interface{}) {
	obj := apimsg.GetRes("DataFormat", msg)
	c.AbortWithStatusJSON(obj.Status, obj.Obj)
}

// Forbidden -
func (c *Context) Forbidden(msg interface{}) {
	obj := apimsg.GetRes("Forbidden", msg)
	c.AbortWithStatusJSON(obj.Status, obj.Obj)
}

// Success -
func (c *Context) Success(msg interface{}) {
	obj := apimsg.GetRes("Success", msg)
	c.AbortWithStatusJSON(obj.Status, obj.Obj)
}

// ServerError response
func (c *Context) ServerError(msg interface{}) {
	obj := apimsg.GetRes("InternalError", msg)
	c.AbortWithStatusJSON(obj.Status, obj.Obj)
}
