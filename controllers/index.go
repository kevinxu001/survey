package controllers

import ()

type IndexController struct {
	CommonController
}

	Day    string
func (this *IndexController) Get() {
	this.TplNames = "index.html"
}

func (this *IndexController) Dashboard() {
	this.TplNames = "index.html"
	this.Data["isDashboard"] = true
}

// func (this *IndexController) Blank() {
// 	this.TplNames = "blank.html"
// }
