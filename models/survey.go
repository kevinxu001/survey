package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type SurveyTask struct {
	Id          int
	TaskName    string    `orm:"size(60)" form:"TaskName" valid:"Required;MaxSize(60)"`
	TaskCreated time.Time `orm:"auto_now_add;type(datetime)" form:"TaskCreated"`
	TaskStarted time.Time `orm:"type(datetime)" form:"TaskStarted"`
	TaskEnded   time.Time `orm:"type(datetime)" form:"TaskEnded"`
	User        *User     `orm:"rel(fk)"`
}

type SurveyClass struct {
	Id         int
	Pid        int         `orm:"default(0)" form:"Pid" valid:"Numeric"`
	ClassName  string      `orm:"size(60)" form:"ClassName" valid:"Required;MaxSize(60)"`
	Desc       string      `orm:"null;size(200)" form:"Desc" valid:"MaxSize(200)"`
	SurveyTask *SurveyTask `orm:"rel(fk)"`
}

type SurveyItem struct {
	Id          int
	ItemName    string       `orm:"size(60)" form:"ItemName" valid:"Required;MaxSize(60)"`
	Desc        string       `orm:"null;size(200)" form:"Desc" valid:"MaxSize(200)"`
	PointMin    float32      `orm:"default(0.0)" form:"PointMin" valid:"Numeric"`
	PointMax    float32      `orm:"default(5.0)" form:"PointMax" valid:"Numeric"`
	PointStep   float32      `orm:"default(1.0)" form:"PointStep" valid:"Numeric"`
	SurveyTask  *SurveyTask  `orm:"rel(fk)"`
	SurveyClass *SurveyClass `orm:"rel(fk)"`
}

// func (o *OrganizationUnit) TableName() string {
// 	return "organization_unit"
// }

// func (o *OrganizationUnit) TableIndex() [][]string {
// 	return [][]string{
// 		[]string{"Id","UnitName"},
// 	}
// }

// func (o *OrganizationUnit) TableUnique() [][]string {
// 	return [][]string{
// 		[]string{"UnitName"}
// 	}
// }

func init() {
	orm.RegisterModelWithPrefix(beego.AppConfig.String("tableprefix"), new(SurveyTask), new(SurveyClass), new(SurveyItem))
}
