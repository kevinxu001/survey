package controllers

import (
	"github.com/astaxie/beego"
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
}
