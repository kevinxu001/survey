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
	ClassDesc  string      `orm:"null;type(text)" form:"ClassDesc" valid:"MaxSize(1000)"`
	SortRank   int         `orm:"default(1)" form:"SortRank" valid:"Range(1,100)"`
	SurveyTask *SurveyTask `orm:"rel(fk)"`
}

type SurveyItem struct {
	Id          int
	ItemName    string       `orm:"null;size(60)" form:"ItemName" valid:"Required;MaxSize(60)"`
	ItemDesc    string       `orm:"null;type(text)" form:"ItemDesc" valid:"MaxSize(1000)"`
	PointMin    float32      `orm:"default(0.0)" form:"PointMin" valid:"Numeric"`
	PointMax    float32      `orm:"default(5.0)" form:"PointMax" valid:"Numeric"`
	PointStep   float32      `orm:"default(1.0)" form:"PointStep" valid:"Numeric"`
	SortRank    int          `orm:"default(1)" form:"SortRank" valid:"Range(1,100)"`
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
