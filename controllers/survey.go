package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/models"
	"strconv"
	"strings"
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

	this.Data["SurveyTask"] = &surveytask

	//读取登陆用户信息，并读取该用户已提交的填报数据并生成对应的map[itemid]point
	//判断是否是普通用户登录
	v := this.GetSession("currentUser")
	user, ok := v.(models.User)
	if !ok {
		this.Data["hasErrorMsg"] = true
		this.Data["errorMsg"] = "只有注册用户才能参与填报，请退出并重新登录！"
		return
	}
	//读取用户的组织机构
	o.Read(user.OrgUnit)
	this.Data["OrganizationUnitName"] = user.OrgUnit.UnitName

	//一级分类的名称及总填报项目数和本用户已填报项目数
	type ClassCount struct {
		FirstClassDesc string
		FilledSum      float32
	}
	classFilledSum := make(map[int]ClassCount)

	qsclass := o.QueryTable("survey_class")
	var firstclasses []models.SurveyClass
	qsclass.Filter("surveytask__id", surveytask.Id).Filter("Pid", 0).OrderBy("SortRank").All(&firstclasses)

	for _, sClass := range firstclasses {
		classFilledSum[sClass.Id] = ClassCount{FirstClassDesc: sClass.ClassDesc}
	}

	type ItemSum struct {
		FirstClassId int
		ItemSum      float32
	}
	var itemsums []ItemSum
	_, err = o.Raw("SELECT first_class_id,sum(item_point) as item_sum FROM `survey_data` WHERE survey_task_id=? and organization_unit_id=? group by first_class_id order by first_class_id asc", surveytask.Id, user.OrgUnit.Id).QueryRows(&itemsums)

	for _, itemsum := range itemsums {
		t := classFilledSum[itemsum.FirstClassId]
		t.FilledSum = itemsum.ItemSum
		classFilledSum[itemsum.FirstClassId] = t
	}
	//增加总分和当前总得分
	var totalsum orm.ParamsList
	_, err = o.Raw("SELECT sum(point_max) as total FROM `survey_item` WHERE survey_task_id=?", surveytask.Id).ValuesFlat(&totalsum)
	var total string
	total, ok = totalsum[0].(string)

	var allsum orm.ParamsList
	_, err = o.Raw("SELECT sum(item_point) as item_sum FROM `survey_data` WHERE survey_task_id=? and organization_unit_id=?", surveytask.Id, user.OrgUnit.Id).ValuesFlat(&allsum)
	var sumf float64
	if sum, ok := allsum[0].(string); ok {
		sumf, _ = strconv.ParseFloat(sum, 32)
	}

	classFilledSum[0] = ClassCount{FirstClassDesc: "总分(" + total + "分)", FilledSum: float32(sumf)}

	this.Data["ClassFilledSum"] = &classFilledSum

	//读取已上传文件信息
	qsfile := o.QueryTable("file")
	var uploadfiles []models.File
	_, err = qsfile.Filter("organizationunit__id", user.OrgUnit.Id).Filter("surveytask__id", surveytask.Id).All(&uploadfiles)
	if err != nil {
		beego.Error(err)
	}

	this.Data["UploadFiles"] = &uploadfiles
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
		return
	}
	if tnow.After(surveytask.TaskEnded) {
		this.Data["hasErrorMsg"] = true
		this.Data["errorMsg"] = "调研已经结束！"
		return
	}

	this.Data["SurveyTask"] = &surveytask

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

	//读取登陆用户信息，并读取该用户已提交的填报数据并生成对应的map[itemid]point
	//判断是否是普通用户登录
	v := this.GetSession("currentUser")
	user, ok := v.(models.User)
	if !ok {
		this.Data["hasErrorMsg"] = true
		this.Data["errorMsg"] = "只有注册用户才能参与填报，请退出并重新登录！"
		return
	}
	o.Read(user.OrgUnit)
	//读取用户已填报的数据，并建立对应的map数据返回页面，如果没有填报的数据，则对应的point为该项最大值
	var itemValues = make(map[int]float32, len(sitems))
	for _, v := range sitems {
		itemValues[v.Id] = v.PointMax
	}

	qsdata := o.QueryTable("survey_data")
	var filledData []models.SurveyData
	_, err = qsdata.Filter("organizationunit__id", user.OrgUnit.Id).Filter("surveytask__id", taskid).All(&filledData)
	if err != nil {
		beego.Error(err)
	}

	for _, filldata := range filledData {
		itemValues[filldata.SurveyItem.Id] = filldata.ItemPoint
	}

	this.Data["ItemValues"] = &itemValues

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

func (this *SurveyController) GetSurveys() {
	//根据select2传入的q参数查找对应taskname的调研任务并返回数据
	q := this.GetString("q")
	pageLimit, _ := this.GetInt("page_limit")
	page, _ := this.GetInt("page")

	o := orm.NewOrm()
	qs := o.QueryTable("survey_task")

	var surveyTasks []models.SurveyTask

	if q != "" {
		qs.Filter("taskname__icontains", q).Limit(pageLimit, (page-1)*pageLimit).All(&surveyTasks)
	} else {
		qs.All(&surveyTasks)
	}

	this.Data["json"] = surveyTasks
	this.ServeJson()
}

func (this *SurveyController) PostSurveyById() {
	//读取调研的名称和时间范围，并返回未开始、调研结束等信息。
	taskid, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	firstclassid, _ := this.GetInt("firstclassid")

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
		this.Data["json"] = &Rsp{Success: false, Msg: "调研还未开始！"}
		this.ServeJson()
		return
	}
	if tnow.After(surveytask.TaskEnded) {
		this.Data["json"] = &Rsp{Success: false, Msg: "调研已经结束！"}
		this.ServeJson()
		return
	}

	//读取post的记录数据，抽取出对应的item-id，比对数据库中此用户已入库的数据，有则不做操作，无则入库
	var fillValues = make(map[int]float32, len(this.Input()))
	var itemid int
	var point float64
	for k, v := range this.Input() {
		if strings.HasPrefix(k, "item") {
			itemid, _ = strconv.Atoi(strings.TrimPrefix(k, "item"))
			point, _ = strconv.ParseFloat(v[0], 32)
			fillValues[itemid] = float32(point)
		}
	}
	//读取数据库中已填报的记录数据
	//判断是否是普通用户登录
	v := this.GetSession("currentUser")
	user, ok := v.(models.User)
	if !ok {
		this.Data["hasErrorMsg"] = true
		this.Data["errorMsg"] = "只有注册用户才能参与填报，请退出并重新登录！"
		return
	}
	//读取用户的组织机构信息
	o.Read(user.OrgUnit)

	qsdata := o.QueryTable("survey_data")
	var filledData []models.SurveyData
	_, err = qsdata.Filter("organizationunit__id", user.OrgUnit.Id).Filter("surveytask__id", taskid).RelatedSel().All(&filledData, "id")

	if err != nil {
		beego.Error(err)
	}
	//如果已有填报数据，则更新此数据内容
	for _, filldata := range filledData {
		if _, ok := fillValues[filldata.SurveyItem.Id]; ok {
			filldata.ItemPoint = fillValues[filldata.SurveyItem.Id]
			_, err = o.Update(&filldata)
			if err != nil {
				beego.Error(err)
			}
			delete(fillValues, filldata.SurveyItem.Id)
		}
	}
	//如果还没有，则插入数据
	surveydatas := make([]*models.SurveyData, 0, len(fillValues))
	sTask := new(models.SurveyTask)
	sTask.Id = taskid
	for itemid, filldata := range fillValues {
		sClass := new(models.SurveyClass)
		sClass.Id = int(firstclassid)

		sItem := new(models.SurveyItem)
		sItem.Id = itemid

		sData := new(models.SurveyData)

		sData.ItemPoint = filldata
		sData.SurveyTask = sTask
		sData.FirstClass = sClass
		sData.SurveyItem = sItem
		sData.OrganizationUnit = user.OrgUnit
		sData.User = &user
		surveydatas = append(surveydatas, sData)
	}
	// for k, v := range surveydatas {
	// 	beego.Info(k, ":", v)
	// }
	_, err = o.InsertMulti(100, surveydatas)
	if err != nil {
		beego.Info(err)
	}

	this.Data["json"] = &Rsp{Success: true, Msg: "数据已提交。"}
	this.ServeJson()
}
