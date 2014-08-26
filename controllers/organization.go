package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/models"
)

type OrganizationController struct {
	CommonController
}

func (this *OrganizationController) Get() {
	this.TplNames = "organization/list.html"
	this.Data["isSystemAdmin"] = true
	this.Data["isOrganization"] = true

	pid, err := this.GetInt("pid")
	if err != nil {
		beego.Error(err)
		return
	}

	o := orm.NewOrm()
	qs := o.QueryTable("organization_unit")

	var organizations []*models.OrganizationUnit
	qs.Filter("pid", pid).All(&organizations)

	this.Data["json"] = organizations
	this.ServeJson()
}

func (this *OrganizationController) Add() {
	pid, err := this.GetInt("pid")
	unitname := this.GetString("UnitName")
	sortrank, _ := this.GetInt("SortRank")
	remark := this.GetString("Remark")

	o := orm.NewOrm()
	qs := o.QueryTable("organization_unit")

	num, _ := qs.Filter("unitname", unitname).Count()
	if num > 0 {
		this.Data["json"] = &Rsp{Success: false, Msg: "已有相同名称的学校！"}
		this.ServeJson()
		return
	}

	var organization models.OrganizationUnit
	organization.UnitName = unitname
	organization.Pid = int(pid)
	organization.SortRank = uint8(sortrank)
	organization.Remark = remark
	organization.Status = 1

	_, err = o.Insert(&organization)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "新建学校失败！"}
		this.ServeJson()
		return
	}

	this.Data["json"] = &Rsp{Success: true, Msg: "新建学校成功！"}
	this.ServeJson()
}
