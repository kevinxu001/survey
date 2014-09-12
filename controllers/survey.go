package controllers

import ()

type SurveyController struct {
	CommonController
}

func (this *SurveyController) Get() {
	this.TplNames = "survey/list.html"
	this.Data["isSurveys"] = true
}

func (this *SurveyController) GetSurveyById() {
	this.TplNames = "survey/listsurvey.html"
	this.Data["isSurveys"] = true

	// id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	// beego.Info(id)

	// o:=orm.NewOrm()
	// qs:=o.QueryTable("surveytask")

}
