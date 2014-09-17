package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
)

//初始化cache
func InitCache() {

}

//初始化session
func InitSession() {

}

//初始化logs
func InitLogs() {

}

//初始化i18n
func InitI18n() {

}

func FilterUser(ctx *context.Context) {
	adminUrl := []string{"organizations", "users", "reports"}

	for _, url := range adminUrl {
		if strings.Contains(ctx.Request.RequestURI, url) {
			if adminUser, ok := ctx.Input.Session("adminUser").(string); !ok || adminUser != beego.AppConfig.String("conf::admin_user") {
				ctx.Redirect(302, "/login")
			}
			return
		}
	}

}

//初始化RBAC权限控制
func InitAccessControl() {
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
}
