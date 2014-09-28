package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	// "github.com/kevinxu001/survey/lib"
	"fmt"
	"github.com/kevinxu001/survey/models"
	"strconv"
	// "time"
)

type ReportController struct {
	CommonController
}

func (this *ReportController) ReportMake() {
	this.TplNames = "report/make.html"
	this.Data["isReportAdmin"] = true
	this.Data["isReportMake"] = true
}

func (this *ReportController) ReportMakeByTaskId() {
	taskid, err := strconv.Atoi(this.Ctx.Input.Param(":taskid"))
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "请选择要生成报表数据的调研任务！"}
		this.ServeJson()
		return
	}

	//读取对应调研任务的所有一级分类，将来可实现读取调研任务的数据报表模版来生成数据
	o := orm.NewOrm()
	qsclass := o.QueryTable("survey_class")
	var firstClasses []models.SurveyClass
	_, err = qsclass.Filter("surveytask__id", taskid).Filter("pid", 0).All(&firstClasses)

	//根据一级分类来生成对应的组织机构，分类ID，总分等字段的数据统计表（表名为report_task调研任务id）
	tableName := "report_task" + this.Ctx.Input.Param(":taskid")
	sqlstring := fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", tableName)
	//每次都删除统计表
	_, err = o.Raw(sqlstring).Exec()
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "删除统计数据表失败！"}
		this.ServeJson()
		return
	}

	sqlstring = fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (`id` int(11) unsigned NOT NULL AUTO_INCREMENT,`organization_unit_id` int(11) NOT NULL,", tableName)
	for _, firstclass := range firstClasses {
		sqlstring += fmt.Sprintf("`class%s` double DEFAULT NULL,", strconv.Itoa(firstclass.Id))
	}
	sqlstring += fmt.Sprintf("`total` double DEFAULT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1;")

	//重新建表并统计，效率不高
	_, err = o.Raw(sqlstring).Exec()
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "生成统计数据表失败！"}
		this.ServeJson()
		return
	}

	//开始进行相应的统计，并数据入表
	//遍历所有的调研项目并生成查询的in列表

	//先遍历所有的组织机构及其ID
	// qsou := o.QueryTable("organization_unit")
	// var organizations []models.OrganizationUnit

	//组织机构两级，区和校
	//统计各个学校的报表数据，如果学校已有上报数据的情况下
	//先插入各个组织机构id数据
	var oulist orm.ParamsList
	_, err = o.Raw("SELECT distinct `organization_unit_id` FROM `survey_data` ORDER BY `organization_unit_id` ASC").ValuesFlat(&oulist)

	ouidlist := make([]string, 0, len(oulist))
	for _, orgid := range oulist {
		if orgstr, ok := orgid.(string); ok {
			ouidlist = append(ouidlist, orgstr)
		}
	}

	insertSql := fmt.Sprintf("INSERT INTO `%s` (`organization_unit_id`) VALUES(?);", tableName)
	ip, _ := o.Raw(insertSql).Prepare()
	for _, orgstr := range ouidlist {
		orgid, _ := strconv.Atoi(orgstr)
		ip.Exec(orgid)
	}
	ip.Close()

	//再根据firstclassid依次统计每个学校的分数
	type classsum struct {
		OrgId    int
		SumCount float32
	}
	for _, firstclass := range firstClasses {
		var classsums []classsum
		_, err = o.Raw("SELECT `organization_unit_id` AS org_id,SUM(`item_point`) AS sum_count FROM `survey_data` WHERE `first_class_id`=? GROUP BY `organization_unit_id` ORDER BY `organization_unit_id` ASC", firstclass.Id).QueryRows(&classsums)

		updateSql := fmt.Sprintf("UPDATE `%s` SET `class%d`= ? WHERE `organization_unit_id`= ?", tableName, firstclass.Id)
		up, _ := o.Raw(updateSql).Prepare()
		for _, sum := range classsums {
			up.Exec(sum.SumCount, sum.OrgId)
		}
		up.Close()
	}

	//再统计每个学校的总分
	var classsums []classsum
	_, err = o.Raw("SELECT `organization_unit_id` AS org_id,SUM(`item_point`) AS sum_count FROM `survey_data` GROUP BY `organization_unit_id` ORDER BY `organization_unit_id` ASC").QueryRows(&classsums)

	updateTotalSql := fmt.Sprintf("UPDATE `%s` SET `total`= ? WHERE `organization_unit_id`= ?", tableName)
	upt, _ := o.Raw(updateTotalSql).Prepare()
	for _, sum := range classsums {
		upt.Exec(sum.SumCount, sum.OrgId)
	}
	upt.Close()

	//统计各个区的报表数据

	//统计大市的报表数据，默认情况下有个虚拟的全市的概念，对应的组织机构id为0

	// for _, ou := range organizations {
	// 	beego.Info(ou)
	// }

	this.Data["json"] = true
	this.ServeJson()
}

func (this *ReportController) ReportCity() {
	this.TplNames = "report/city.html"
	this.Data["isReportAdmin"] = true
	this.Data["isReportCity"] = true

}

func (this *ReportController) ReportArea() {
	this.TplNames = "report/area.html"
	this.Data["isReportAdmin"] = true
	this.Data["isReportArea"] = true

}

func (this *ReportController) ReportSchool() {
	this.TplNames = "report/school.html"
	this.Data["isReportAdmin"] = true
	this.Data["isReportSchool"] = true

}
