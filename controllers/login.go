package controllers

import ()

type LoginController struct {
	CommonController
}

func (this *LoginController) Get() {
	this.TplNames = "login/login.html"
}

func (this *LoginController) Post() {

}

func (this *LoginController) Logout() {

}
