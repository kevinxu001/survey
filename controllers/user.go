package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/lib"
	"github.com/kevinxu001/survey/models"
	"strconv"
)

type UserController struct {
	CommonController
}

type UserDataRsp struct {
	Draw            int64           `json:"draw"`
	RecordsTotal    int64           `json:"recordsTotal"`
	RecordsFiltered int64           `json:"recordsFiltered"`
	Data            *[]*models.User `json:"data"`
}

func (this *UserController) Get() {
	this.TplNames = "user/list.html"
	this.Data["isSystemAdmin"] = true
	this.Data["isUsers"] = true
}

func (this *UserController) GetUsers() {
	//根据select2传入的q参数查行类似unitname的查找并返回数据
	draw, err := this.GetInt("draw")
	start, _ := this.GetInt("start")
	length, _ := this.GetInt("length")
	searchValue := this.GetString("search[value]")

	//根据oid获取所有oid下的用户数据
	oid, _ := this.GetInt("oid")

	o := orm.NewOrm()
	qs := o.QueryTable("user")

	var (
		recordsTotal    int64
		recordsFiltered int64
	)

	var users []*models.User

	var searchCond, condAll *orm.Condition
	cond := orm.NewCondition()
	if searchValue != "" {
		searchCond = cond.And("username__icontains", searchValue).Or("realname__icontains", searchValue).Or("idcard__icontains", searchValue)
	}

	if oid > 0 {
		condAll = cond.AndCond(searchCond).AndCond(cond.And("orgunit", oid))
		recordsTotal, err = qs.SetCond(condAll).Count()
		recordsFiltered = recordsTotal

		qs.SetCond(condAll).Limit(length, start).All(&users)
	} else if oid == 0 {
		recordsTotal, err = qs.SetCond(searchCond).Count()
		recordsFiltered = recordsTotal

		_, err = qs.SetCond(searchCond).Limit(length, start).All(&users)
	}
	if err != nil {
		beego.Error(err)
	}

	this.Data["json"] = &UserDataRsp{Draw: draw, RecordsTotal: recordsTotal, RecordsFiltered: recordsFiltered, Data: &users}
	this.ServeJson()
}

func (this *UserController) Add() {
	oid, err := this.GetInt("oid")
	username := this.GetString("username")
	realname := this.GetString("realname")
	password := this.GetString("password")
	confirmpassword := this.GetString("confirmpassword")
	if password != confirmpassword {
		this.Data["json"] = &Rsp{Success: false, Msg: "密码与确认密码不相同，插入数据出错！"}
		this.ServeJson()
		return
	}
	mobile := this.GetString("mobile")
	phone := this.GetString("phone")
	idcard := this.GetString("idcard")

	o := orm.NewOrm()
	qs := o.QueryTable("user")

	num, _ := qs.Filter("username", username).Count()
	if num > 0 {
		this.Data["json"] = &Rsp{Success: false, Msg: "已有相同用户名的用户，插入数据出错！"}
		this.ServeJson()
		return
	}

	var organization models.OrganizationUnit
	qsou := o.QueryTable("organization_unit")
	err = qsou.Filter("id", oid).One(&organization)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "无法找到对应的组织机构，插入数据出错！"}
		this.ServeJson()
		return
	}

	var user models.User
	user.OrgUnit = &organization
	user.UserName = username
	user.PassWord = lib.StrToMD5(password)
	user.RealName = realname
	user.Mobile = mobile
	user.Phone = phone
	user.IdCard = idcard

	_, err = o.Insert(&user)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "无法新增用户，插入数据出错！"}
		this.ServeJson()
		return
	}

	this.Data["json"] = &Rsp{Success: true, Msg: "成功新增用户！"}
	this.ServeJson()
}

func (this *UserController) GetUserById() {
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "无法获取用户信息！"}
		this.ServeJson()
		return
	}

	o := orm.NewOrm()
	qs := o.QueryTable("user")

	var user models.User
	qs.Filter("id", id).One(&user)

	this.Data["json"] = &user
	this.ServeJson()
}

func (this *UserController) Modify() {
	id, err := this.GetInt("id")
	o := orm.NewOrm()
	qs := o.QueryTable("user")

	var user models.User
	err = qs.Filter("id", id).One(&user)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "用户数据不存在，编辑用户失败！"}
		this.ServeJson()
		return
	}

	oid, _ := this.GetInt("oid")
	username := this.GetString("username")
	realname := this.GetString("realname")
	password := this.GetString("password")
	confirmpassword := this.GetString("confirmpassword")
	if password != confirmpassword {
		this.Data["json"] = &Rsp{Success: false, Msg: "密码与确认密码不相同，编辑用户出错！"}
		this.ServeJson()
		return
	}
	mobile := this.GetString("mobile")
	phone := this.GetString("phone")
	idcard := this.GetString("idcard")

	if user.UserName != username {
		num, _ := qs.Filter("username", username).Count()
		if num > 0 {
			this.Data["json"] = &Rsp{Success: false, Msg: "已有相同用户名的用户，编辑用户出错！"}
			this.ServeJson()
			return
		}
	}

	var organization models.OrganizationUnit
	qsou := o.QueryTable("organization_unit")
	err = qsou.Filter("id", oid).One(&organization)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "无法找到对应的组织机构，编辑用户出错！"}
		this.ServeJson()
		return
	}

	user.OrgUnit = &organization
	user.UserName = username
	if password != "" {
		user.PassWord = lib.StrToMD5(password)
	}
	user.RealName = realname
	user.Mobile = mobile
	user.Phone = phone
	user.IdCard = idcard

	_, err = o.Update(&user)
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "编辑用户出错！"}
		this.ServeJson()
		return
	}

	this.Data["json"] = &Rsp{Success: true, Msg: "成功编辑用户！"}
	this.ServeJson()
}

func (this *UserController) Delete() {
	ids := this.GetString("ids")

	o := orm.NewOrm()

	_, err := o.Raw("delete from `user` where id in (" + ids + ")").Exec()
	if err != nil {
		beego.Error(err)
		this.Data["json"] = &Rsp{Success: false, Msg: "删除用户出错！"}
		this.ServeJson()
		return
	}

	this.Data["json"] = &Rsp{Success: true, Msg: "成功删除用户！"}
	this.ServeJson()
}
