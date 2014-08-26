package controllers

import ()

type UserController struct {
	CommonController
}

func (this *UserController) Get() {
	this.TplNames = "user/list.html"
}

func (this *UserController) AdminList() {
	this.TplNames = "user/adminlist.html"
}
