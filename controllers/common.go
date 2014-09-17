package controllers

import (
	"github.com/astaxie/beego"
	"github.com/kevinxu001/survey/models"
)

type CommonController struct {
	beego.Controller
}

type Rsp struct {
	Success bool
	Msg     string
}

func (this *CommonController) Prepare() {
	this.Data["SiteName"] = beego.AppConfig.String("conf::site_name")

	path := this.Ctx.Request.RequestURI
	if path == "/login" || path == "/api/login" {
		return
	} else {
		//判断是不是管理员登录
		vadmin := this.GetSession("adminUser")
		if admin_user, ok := vadmin.(string); ok {
			if admin_user != beego.AppConfig.String("conf::admin_user") {
				this.Redirect("/login", 302)
				return
			} else {
				admin_realname := beego.AppConfig.String("conf::admin_realname")
				this.Data["isAdminUser"] = true
				this.Data["Realname"] = admin_realname
				return
			}
		}

		//判断是否是普通用户登录
		v := this.GetSession("currentUser")
		if user, ok := v.(models.User); !ok {
			this.Redirect("/login", 302)
			return
		} else {
			this.Data["Realname"] = user.RealName
			return
		}
	}
}
