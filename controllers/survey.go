package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/models"
	"strconv"
	"time"
)

type SurveyController struct {
	CommonController
}

func (this *SurveyController) Get() {
	this.TplNames = "survey/showsurvey.html"
	this.Data["isSurveys"] = true
}

func (this *SurveyController) GetSurveyById() {
	this.TplNames = "survey/showsurvey.html"
	this.Data["isSurveys"] = true

	//taskid, err := strconv.Atoi(this.Ctx.Input.Param(":id"))

}

func (this *SurveyController) FillinSurveyById() {
	this.TplNames = "survey/fillinsurvey.html"
	this.Data["isSurveys"] = true

	//读取调研的名称和时间范围，并返回未开始、调研结束等信息。
	taskid, err := strconv.Atoi(this.Ctx.Input.Param(":id"))

	o := orm.NewOrm()
	qstask := o.QueryTable("survey_task")

	var surveytask models.SurveyTask
	err = qstask.Filter("id", taskid).One(&surveytask)
	if err != nil {
		beego.Error(err)
		this.Redirect("/errorpage/500", 302)
		return
	}

	tnow := time.Now()
	// year, month, day := t.Date()
	// t = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	if tnow.Before(surveytask.TaskStarted) {
		this.Data["hasErrorMsg"] = true
		this.Data["errorMsg"] = "调研还未开始！"
	}
	if tnow.After(surveytask.TaskEnded) {
		this.Data["hasErrorMsg"] = true
		this.Data["errorMsg"] = "调研已经结束！"
	}

	this.Data["taskName"] = surveytask.TaskName

	//读取所有调研项目的classid和对应id的项目数量
	type ClassCount struct {
		ClassId   int
		ItemCount int
	}
	var classids []ClassCount
	_, err = o.Raw("select survey_class_id as class_id,count(1) as item_count from `survey_item` where survey_task_id=? group by survey_class_id order by survey_class_id asc", surveytask.Id).QueryRows(&classids)

	surveyItemMaps := make(map[int][]models.SurveyItem, len(classids))

	//读取所有调研项目内容
	qsitem := o.QueryTable("survey_item")
	var sitems []models.SurveyItem
	_, err = qsitem.Filter("surveytask__id", surveytask.Id).OrderBy("surveyclass__id", "sortrank").All(&sitems)

	var n int = 0
	for _, classid := range classids {
		surveyItemMaps[classid.ClassId] = sitems[n : n+classid.ItemCount]
		n += classid.ItemCount
	}

	this.Data["SurveyItemMaps"] = &surveyItemMaps

	//读取对应的调研任务的分类,本项目目前最多支持两级分类
	//先读取一级分类
	qsclass := o.QueryTable("survey_class")
	var firstClasses []models.SurveyClass
	qsclass.Filter("surveytask__id", surveytask.Id).Filter("pid", 0).OrderBy("id", "sortrank").All(&firstClasses)

	firstClassMaps := make(map[int]models.SurveyClass, len(firstClasses))
	for k, v := range firstClasses {
		firstClassMaps[k+1] = v
	}
	this.Data["FirstClassMaps"] = &firstClassMaps

	//读取二级分类
	//读取一级分类的id和下级类别的数量
	_, err = o.Raw("select pid as class_id,count(1) as item_count from `survey_class` where survey_task_id=? and pid>0 group by pid order by pid asc", surveytask.Id).QueryRows(&classids)

	secondClassMaps := make(map[int][]models.SurveyClass, len(classids))

	//读取所有二级分类内容
	var secondClasses []models.SurveyClass
	qsclass.Filter("surveytask__id", surveytask.Id).Filter("pid__gt", 0).OrderBy("pid", "sortrank").All(&secondClasses)

	n = 0
	for _, classid := range classids {
		secondClassMaps[classid.ClassId] = secondClasses[n : n+classid.ItemCount]
		n += classid.ItemCount
	}

	this.Data["SecondClassMaps"] = &secondClassMaps
}
