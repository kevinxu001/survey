package controllers

import ()

type IndexController struct {
	CommonController
}

func (this *IndexController) Get() {
	this.TplNames = "index.html"
}
