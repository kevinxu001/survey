package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/models"
	"strconv"
)

type OrganizationController struct {
	CommonController
}

func (this *OrganizationController) Get() {
	this.TplNames = "organization/list.html"
	this.Data["isSystemAdmin"] = true
	this.Data["isOrganizations"] = true
}

func (this *OrganizationController) GetOrganizations() {
	//根据select2传入的q参数查找类似unitname的查找并返回数据
	q := this.GetString("q")
	pageLimit, _ := this.GetInt("page_limit")
	page, _ := this.GetInt("page")
	//根据pid获取所有pid下的下级组织机构数据
	pid, _ := this.GetInt("pid")

	o := orm.NewOrm()
	qs := o.QueryTable("organization_unit")

	var organizations []*models.OrganizationUnit

	if q != "" {
		qs.Filter("unitname__icontains", q).Limit(pageLimit, (page-1)*pageLimit).All(&organizations)
	} else if pid >= 0 {
		qs.Filter("pid", pid).All(&organizations)
	}

	this.Data["json"] = organizations
	this.ServeJson()
}

func (this *OrganizationController) GetOrganizationsById() {
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "无法获取组织机构相关信息！"}
		this.ServeJson()
		return
	}

	o := orm.NewOrm()
	qs := o.QueryTable("organization_unit")

	var organizations []*models.OrganizationUnit
	qs.Filter("id", id).All(&organizations)

	this.Data["json"] = organizations
	this.ServeJson()
}

func (this *OrganizationController) Add() {
	pid, err := this.GetInt("Pid")
	unitname := this.GetString("UnitName")
	sortrank, err := this.GetInt("SortRank")
	remark := this.GetString("Remark")

	o := orm.NewOrm()
	qs := o.QueryTable("organization_unit")

	num, _ := qs.Filter("unitname", unitname).Count()
	if num > 0 {
		this.Data["json"] = &Rsp{Success: false, Msg: "已有相同名称的组织机构！"}
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
		this.Data["json"] = &Rsp{Success: false, Msg: "无法新增组织机构，插入数据出错！"}
		this.ServeJson()
		return
	}

	_, err = qs.Filter("id", pid).Update(orm.Params{"status": 0})

	this.Data["json"] = &Rsp{Success: true, Msg: "成功新增组织机构！"}
	this.ServeJson()
}

func (this *OrganizationController) Modify() {
	id, err := this.GetInt("id")
	beego.Info(id)
	pid, err := this.GetInt("Pid")
	unitname := this.GetString("UnitName")
	sortrank, err := this.GetInt("SortRank")
	remark := this.GetString("Remark")

	o := orm.NewOrm()
	qs := o.QueryTable("organization_unit")

	num, _ := qs.Filter("id", id).Count()
	if num != 1 {
		this.Data["json"] = &Rsp{Success: false, Msg: "找不到此组织机构，无法修改！"}
		this.ServeJson()
		return
	}

	var organization models.OrganizationUnit
	err = qs.Filter("id", id).One(&organization)
	if err != nil {
		this.Data["json"] = &Rsp{Success: false, Msg: "读取组织机构信息出错！"}
		this.ServeJson()
		return
	}
	organization.Pid = int(pid)
	organization.UnitName = unitname
	organization.SortRank = uint8(sortrank)
	organization.Remark = remark
	organization.Status = 1

	_, err = o.Update(&organization)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "无法修改组织机构，插入数据有错！"}
		this.ServeJson()
		return
	}

	//更新有子节点的status状态为0,无子结点的状态为1
	_, err = o.Raw("update `organization_unit` a inner join (select id from `organization_unit` where id in (select distinct pid from `organization_unit`)) b on a.id=b.id set status=0").Exec()
	if err != nil {
		beego.Error(err)
	}
	_, err = o.Raw("update `organization_unit` a inner join (select id from `organization_unit` where id not in (select distinct pid from `organization_unit`)) b on a.id=b.id set status=1").Exec()
	if err != nil {
		beego.Error(err)
	}

	this.Data["json"] = &Rsp{Success: true, Msg: "成功修改组织机构！"}
	this.ServeJson()
}

func (this *OrganizationController) Delete() {
	ids := this.GetString("ids")

	o := orm.NewOrm()
	//更新所有删除id的对应pid的status
	// var pids orm.ParamsList
	// num, err := o.Raw("select distinct pid from `organization_unit`").ValuesFlat(&pids)
	// if err == nil && num > 0 {
	// 	ids := make([]string, len(pids), len(pids))
	// 	for i, id := range pids {
	// 		ids[i] = id.(string)
	// 	}
	// 	_, err = o.Raw("update `organization_unit` set status=0 where id in (" + strings.Join(ids, ",") + ")").Exec()
	// }

	_, err := o.Raw("delete from `organization_unit` where id in (" + ids + ")").Exec()
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "删除组织机构出错！"}
		this.ServeJson()
		return
	}

	//更新有子节点的status状态为0,无子结点的状态为1
	_, err = o.Raw("update `organization_unit` a inner join (select id from `organization_unit` where id in (select distinct pid from `organization_unit`)) b on a.id=b.id set status=0").Exec()
	if err != nil {
		beego.Error(err)
	}
	_, err = o.Raw("update `organization_unit` a inner join (select id from `organization_unit` where id not in (select distinct pid from `organization_unit`)) b on a.id=b.id set status=1").Exec()
	if err != nil {
		beego.Error(err)
	}

	this.Data["json"] = &Rsp{Success: true, Msg: "成功删除组织机构！"}
	this.ServeJson()
}
