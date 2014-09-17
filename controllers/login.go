package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/lib"
	"github.com/kevinxu001/survey/models"
	"time"
)

type LoginController struct {
	CommonController
}

func (this *LoginController) Get() {
	this.TplNames = "login/login.html"

	vadmin := this.GetSession("adminUser")
	if admin_user, ok := vadmin.(string); ok && admin_user == beego.AppConfig.String("conf::admin_user") {
		this.Redirect("/", 302)
		return
	}

	v := this.GetSession("currentUser")
	if _, ok := v.(models.User); ok {
		this.Redirect("/", 302)
	}
}

func (this *LoginController) Login() {
	username := this.GetString("username")
	//密码未加密
	password := this.GetString("password")
	encryptPassword := lib.StrToMD5(password)
	//rememberme, _ := this.GetBool("rememberme")

	if username == "" || password == "" {
		this.Data["json"] = &Rsp{Success: false, Msg: "用户名或密码为空，登录失败，请重新输入用户信息！"}
		this.ServeJson()
		return
	}

	//判断超级管理员登录
	if username == beego.AppConfig.String("conf::admin_user") {
		if encryptPassword == beego.AppConfig.String("conf::admin_pass") {
			this.SetSession("adminUser", username)
			this.SetSession("adminRealname", beego.AppConfig.String("conf::admin_realname"))

			this.Data["json"] = &Rsp{Success: true, Msg: "登录成功！"}
			this.ServeJson()
			return
		}
	}

	//判断普通登录
	o := orm.NewOrm()
	qs := o.QueryTable("user")

	var loginUser models.User

	err := qs.Filter("username", username).Filter("password", encryptPassword).One(&loginUser)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "用户登录失败，请重新输入用户信息！"}
		this.ServeJson()
		return
	}

	if loginUser.UserName != username || loginUser.PassWord != encryptPassword {
		this.Data["json"] = &Rsp{Success: false, Msg: "用户登录失败，请重新输入用户信息！"}
		this.ServeJson()
		return
	}

	loginUser.LastLogin = time.Now()
	_, err = o.Update(&loginUser)
	if err != nil {
		beego.Error(err)
	}

	//用户验证成功，记录用户session
	this.SetSession("currentUser", loginUser)

	this.Data["json"] = &Rsp{Success: true, Msg: "登录成功！"}
	this.ServeJson()
}

func (this *LoginController) Logout() {
	this.TplNames = "login/login.html"

	vadmin := this.GetSession("adminUser")
	if admin_user, ok := vadmin.(string); ok && admin_user == beego.AppConfig.String("conf::admin_user") {
		this.DelSession("adminUser")
		this.DelSession("adminRealname")
		return
	}

	v := this.GetSession("currentUser")
	if _, ok := v.(models.User); ok {
		this.DelSession("currentUser")
	}
}
