package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	_ = iota
	EventStart
	EventEnd
)

type Event struct {
	Id           int
	EventContent string    `orm:"type(text)" form:"EventContent" valid:"Required;MaxSize(300)"`
	EventTime    time.Time `orm:"type(datetime)" form:"EventTime"`
	//事件类型 0:开始事件 1:结束事件
	EventType uint8 `orm:"default(1)" form:"EventType" valid:"Range(1，2)"`
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
	orm.RegisterModelWithPrefix(beego.AppConfig.String("tableprefix"), new(Event))
}
