package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/kevinxu001/survey/models"
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

func (this *UserController) AdminList() {
	this.TplNames = "user/adminlist.html"
}
