package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type File struct {
	Id               int
	FileName         string            `orm:"size(60)" form:"FileName" valid:"Required;MaxSize(60)"`
	FilePath         string            `orm:"size(60)" form:"FilePath" valid:"Required;MaxSize(60)"`
	FileExt          string            `orm:"size(5)" form:"FileExt" valid:"Required;MaxSize(5)"`
	FileSize         int               `orm:"default(0)" form:"FileSize" valid:"Numeric"`
	Uploaded         time.Time         `orm:"auto_now_add;type(datetime)" form:"Uploaded"`
	SurveyTask       *SurveyTask       `orm:"rel(fk)"`
	OrganizationUnit *OrganizationUnit `orm:"rel(fk)"`
	User             *User             `orm:"rel(fk)"`
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
	orm.RegisterModelWithPrefix(beego.AppConfig.String("tableprefix"), new(File))
}
